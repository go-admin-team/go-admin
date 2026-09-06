package storage

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
	"github.com/go-admin-team/go-admin-core/v2/storage/queue"
)

func swapMsg() corestorage.Messager {
	m := new(queue.Message)
	m.SetStream("t")
	m.SetValues(map[string]interface{}{"a": "b"})
	return m
}

// A reload must never leave producers holding a queue that has stopped
// accepting.
//
// Shutdown waits for its consumers to deliver what the queue still holds. Taking
// the old adapter down before installing the new one meant the runtime pointed
// at a closed queue for that entire wait: every Append in the window came back
// ErrQueueClosed, and both call sites in common/middleware log it at error
// level. Installing first leaves no window - a producer gets the new queue or
// the old one, and both accept.
//
// The difference is only visible during that wait, which is why the test holds
// a consumer rather than checking the state after Setup has returned: by then
// the two orders look identical.
func TestAReloadNeverPointsProducersAtAClosedQueue(t *testing.T) {
	prevQ, prevC := config.QueueConfig, config.CacheConfig
	prevRuntime := sdk.Runtime
	prevInstalled, prevGen := installed, installedGen
	t.Cleanup(func() {
		config.QueueConfig, config.CacheConfig = prevQ, prevC
		sdk.Runtime = prevRuntime
		queueMu.Lock()
		installed, installedGen = prevInstalled, prevGen
		queueMu.Unlock()
	})
	sdk.Runtime = runtime.NewConfig()
	config.CacheConfig = &config.Cache{Memory: struct{}{}}
	config.QueueConfig = &config.Queue{Memory: &config.QueueMemory{PoolSize: 64}}

	Setup()

	// A consumer that will not finish until this test lets it, so the reload's
	// Shutdown has something to wait for.
	release := make(chan struct{})
	first := sdk.Runtime.GetQueuePrefix("")
	first.Register("t", func(corestorage.Messager) error {
		<-release
		return nil
	})
	go first.Run()
	for i := 0; i < 4; i++ {
		if err := first.Append(swapMsg()); err != nil {
			t.Fatalf("seed append %d: %v", i, err)
		}
	}
	time.Sleep(100 * time.Millisecond) // let the consumer pick one up and block

	reloaded := make(chan struct{})
	go func() { Setup(); close(reloaded) }()

	// Publish continuously while the reload is in progress.
	var refused atomic.Int64
	var attempts atomic.Int64
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			attempts.Add(1)
			if err := sdk.Runtime.GetQueuePrefix("").Append(swapMsg()); errors.Is(err, corestorage.ErrQueueClosed) {
				refused.Add(1)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(200 * time.Millisecond) // the reload is now inside Shutdown's wait
	close(release)

	select {
	case <-reloaded:
	case <-time.After(10 * time.Second):
		t.Fatal("the reload never finished")
	}
	close(stop)

	if attempts.Load() == 0 {
		t.Fatal("nothing was published during the reload; the test proves nothing")
	}
	if n := refused.Load(); n > 0 {
		t.Errorf("%d of %d publishes during the reload were refused: producers were pointed at the closed queue",
			n, attempts.Load())
	}
}
