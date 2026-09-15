package storage

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
	"github.com/go-admin-team/go-admin-core/v2/storage/queue"
)

// sampleSize is how many publishes have to land inside the reload before the
// measurement is taken. Waiting on the count rather than on wall clock keeps
// the window the test covers the same on a loaded runner as on an idle one.
const sampleSize = 200

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
//
// One refusal survives the fix and is not something this ordering can reach.
// GetQueuePrefix hands back a wrapper that captured the adapter, so a producer
// that fetched before the swap and appends after Shutdown has begun is still
// holding the old one. That window is one call wide and closing it means
// resolving the adapter inside Append, which is core's to change. What the
// ordering removes is the sustained window: every producer that fetches during
// the wait. The test publishes from a single goroutine, so at most one of its
// calls can straddle the swap - which is what makes "more than one" the line
// between the two orders rather than a tolerance.
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
	// Sized so the buffer cannot fill while the consumer is held: a full queue
	// returns an error of its own, and this test needs every error other than
	// ErrQueueClosed to mean something it does not model has happened.
	config.QueueConfig = &config.Queue{Memory: &config.QueueMemory{PoolSize: 4096}}

	Setup()

	// A consumer that will not finish until this test lets it, so the reload's
	// Shutdown has something to wait for.
	release := make(chan struct{})
	consuming := make(chan struct{})
	var picked sync.Once
	first := sdk.Runtime.GetQueuePrefix("")
	first.Register("t", func(corestorage.Messager) error {
		picked.Do(func() { close(consuming) })
		<-release
		return nil
	})
	go first.Run()
	for i := 0; i < 4; i++ {
		if err := first.Append(swapMsg()); err != nil {
			t.Fatalf("seed append %d: %v", i, err)
		}
	}
	select {
	case <-consuming:
	case <-time.After(10 * time.Second):
		t.Fatal("the consumer never picked a message up, so the reload has nothing to wait for")
	}

	reloaded := make(chan struct{})
	go func() { Setup(); close(reloaded) }()

	// Publish continuously while the reload is in progress.
	var refused atomic.Int64
	var attempts atomic.Int64
	unexpected := make(chan error, 1)
	stop := make(chan struct{})
	// publishing is closed by the producer on its way out. The test joins on it
	// before returning: t.Cleanup restores sdk.Runtime, and a producer still in
	// flight would be reading the variable that restore writes.
	publishing := make(chan struct{})
	go func() {
		defer close(publishing)
		for {
			select {
			case <-stop:
				return
			default:
			}
			attempts.Add(1)
			err := sdk.Runtime.GetQueuePrefix("").Append(swapMsg())
			switch {
			case err == nil:
			case errors.Is(err, corestorage.ErrQueueClosed):
				refused.Add(1)
			default:
				// Kept rather than counted: an Append refused for some other
				// reason would otherwise leave refused at zero and the test
				// green while nothing was reaching a queue at all.
				select {
				case unexpected <- err:
				default:
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()

	deadline := time.After(30 * time.Second)
	for attempts.Load() < sampleSize {
		select {
		case <-deadline:
			t.Fatalf("only %d publishes landed inside the reload; the window was never sampled", attempts.Load())
		case <-time.After(time.Millisecond):
		}
	}
	close(release)

	select {
	case <-reloaded:
	case <-time.After(30 * time.Second):
		t.Fatal("the reload never finished")
	}
	close(stop)
	<-publishing

	select {
	case err := <-unexpected:
		t.Fatalf("a publish failed for a reason this test does not model: %v", err)
	default:
	}
	if n := refused.Load(); n > 1 {
		t.Errorf("%d of %d publishes during the reload were refused: producers were pointed at the closed queue",
			n, attempts.Load())
	}
}
