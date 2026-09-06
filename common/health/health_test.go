package health

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
)

func freshRuntime(t *testing.T) {
	t.Helper()
	previous := sdk.Runtime
	t.Cleanup(func() { sdk.Runtime = previous })
	sdk.Runtime = runtime.NewConfig()
}

// fakeCache answers whatever the test needs it to.
type fakeCache struct {
	mu      sync.Mutex
	setErr  error
	getErr  error
	getBack string // returned instead of what was written, when non-empty
	stored  map[string]string

	// oneSlot makes the cache keep a single value however many keys are
	// written, which is what a shared probe key turns any cache into.
	oneSlot bool
	slot    string

	// setBarrier, when set, holds every writer until all of them have written.
	// Without it the probes are short enough that the scheduler usually runs
	// them one after another, and a shared key survives by luck rather than by
	// design - which would leave the test below asserting nothing.
	setBarrier *barrier
}

// barrier releases every waiter once n of them have arrived.
type barrier struct {
	n   int
	mu  sync.Mutex
	got int
	ch  chan struct{}
}

func newBarrier(n int) *barrier { return &barrier{n: n, ch: make(chan struct{})} }

func (b *barrier) wait() {
	b.mu.Lock()
	b.got++
	if b.got == b.n {
		close(b.ch)
	}
	b.mu.Unlock()
	<-b.ch
}

func (c *fakeCache) String() string { return "fake" }

func (c *fakeCache) Set(key string, val interface{}, _ int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.setErr != nil {
		return c.setErr
	}
	v, _ := val.(string)
	if c.oneSlot {
		c.slot = v
		return nil
	}
	if c.stored == nil {
		c.stored = map[string]string{}
	}
	c.stored[key] = v
	c.mu.Unlock()
	if c.setBarrier != nil {
		// Outside the lock on purpose: waiting while holding it would deadlock
		// every other writer before the barrier could fill.
		c.setBarrier.wait()
	}
	c.mu.Lock()
	return nil
}

func (c *fakeCache) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.getErr != nil {
		return "", c.getErr
	}
	if c.getBack != "" {
		return c.getBack, nil
	}
	if c.oneSlot {
		return c.slot, nil
	}
	return c.stored[key], nil
}

func (c *fakeCache) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.stored, key)
	return nil
}

func (c *fakeCache) HashGet(_, _ string) (string, error) { return "", nil }
func (c *fakeCache) HashDel(_, _ string) error           { return nil }
func (c *fakeCache) Increase(string) error               { return nil }
func (c *fakeCache) Decrease(string) error               { return nil }
func (c *fakeCache) Expire(string, time.Duration) error  { return nil }

var _ corestorage.AdapterCache = (*fakeCache)(nil)

func named(checks []Check, name string) Check {
	for _, c := range checks {
		if c.Name == name {
			return c
		}
	}
	return Check{Name: name, Err: "check not reported at all"}
}

// A cache that accepts writes and answers every read with a different value is
// the failure this probe exists for - a client pointed at the wrong server, or
// one that silently drops everything. A read alone cannot tell that apart from
// a healthy cache with a cold key, which is why the probe writes first.
func TestCacheProbeFailsWhenTheValueDoesNotComeBack(t *testing.T) {
	freshRuntime(t)
	sdk.Runtime.SetCacheAdapter(&fakeCache{getBack: "something else"})

	got := named(Ready(context.Background()), "cache")
	if got.OK {
		t.Error("the cache check passed although the value written was not the value read back")
	}
	if got.Err == "" {
		t.Error("the failing check reported no reason")
	}
}

func TestCacheProbePassesWhenTheValueComesBack(t *testing.T) {
	freshRuntime(t)
	sdk.Runtime.SetCacheAdapter(&fakeCache{})

	if got := named(Ready(context.Background()), "cache"); !got.OK {
		t.Errorf("the cache check failed for a cache that works: %s", got.Err)
	}
}

func TestCacheProbeReportsAWriteFailure(t *testing.T) {
	freshRuntime(t)
	sdk.Runtime.SetCacheAdapter(&fakeCache{setErr: errors.New("connection refused")})

	got := named(Ready(context.Background()), "cache")
	if got.OK {
		t.Error("the cache check passed although the write failed")
	}
}

// Nothing configured is the case that used to take the process down rather
// than answer. GetCacheAdapter builds a wrapper around whatever is configured
// and returns it even when nothing is, so the value is not nil, the cache
// inside it is, and Set dereferences it - a probe that panics is the worst
// possible answer to "are you well".
//
// Every check has to be reported, passing or not. A probe that omits what it
// could not reach reads as a shorter list of healthy things.
func TestEveryDependencyIsReportedEvenWithNothingConfigured(t *testing.T) {
	freshRuntime(t)

	checks := Ready(context.Background())
	for _, name := range []string{"database", "cache"} {
		c := named(checks, name)
		if c.Err == "check not reported at all" {
			t.Errorf("%s was not reported", name)
		}
		if c.OK {
			t.Errorf("%s passed with nothing configured", name)
		}
	}
	if Healthy(checks) {
		t.Error("Healthy said yes for a process with no database and no cache")
	}
}

// Draining is what makes the shutdown graceful from the outside: it has to be
// observable before the server stops accepting, or the load balancer learns
// about the shutdown by having its connections cut.
func TestDrainingIsObservableOnceItBegins(t *testing.T) {
	previous := draining.Load()
	t.Cleanup(func() { draining.Store(previous) })

	draining.Store(false)
	if Draining() {
		t.Fatal("Draining reported true before shutdown began")
	}
	BeginDraining()
	if !Draining() {
		t.Error("Draining still reported false after BeginDraining")
	}
}

// Two probes at once must both pass. With one fixed key they overwrite each
// other's value between the write and the read, and a readiness probe that
// reports false negatives under load takes healthy instances out of the pool -
// which is worse than not probing at all.
func TestConcurrentProbesDoNotOverwriteEachOther(t *testing.T) {
	freshRuntime(t)
	const probes = 16
	sdk.Runtime.SetCacheAdapter(&fakeCache{setBarrier: newBarrier(probes)})

	var wg sync.WaitGroup
	failures := make(chan string, probes)
	for i := 0; i < probes; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c := named(Ready(context.Background()), "cache"); !c.OK {
				failures <- c.Err
			}
		}()
	}
	wg.Wait()
	close(failures)

	var n int
	var first string
	for err := range failures {
		if n == 0 {
			first = err
		}
		n++
	}
	if n > 0 {
		t.Errorf("%d of %d concurrent probes called a healthy cache broken; first: %s", n, probes, first)
	}
}
