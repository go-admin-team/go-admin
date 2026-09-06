// Package health answers whether this process should be sent traffic.
//
// The two questions an orchestrator asks are not the same one, and go-admin
// answers them at two endpoints:
//
//   - /health is liveness: is the process there at all. It stays a bare 200,
//     because the honest answer to "should I restart you" is almost always no.
//     Restarting a process because its database is unreachable turns one
//     outage into a crash loop that also loses the connection pool, the cache
//     and every in-flight request.
//   - /ready is readiness: should this instance receive requests now. It fails
//     while the dependencies are unreachable, and - the part that only exists
//     because of the life-cycle phases - it fails as soon as shutdown begins,
//     before the server stops accepting.
//
// # What the draining answer is worth
//
// Order alone does not produce a window. Answering before the server stops
// accepting is the right order - the reverse reports the state after the
// connections are already cut - but with nothing between the two they are
// microseconds apart, and a poller on a multi-second interval never sees the
// 503.
//
// The delay between them is extend.shutdown.drain, which is zero unless it is
// configured. On the shipped defaults this is therefore still an answer that
// can be read rather than one anything acts on; a deployment that sets a drain
// window is the one that gets a window to act in.
//
// What acts on it depends on who does the removing. A load balancer that polls
// /ready takes this instance out when it reads the 503, and the window has to
// cover its check interval times its failure threshold, plus however long the
// removal takes to apply. On Kubernetes the endpoint is withdrawn when the Pod
// receives a deletionTimestamp, concurrently with SIGTERM and regardless of
// what the probe returns - there the window covers the delay in that removal
// reaching every node, and the 503 is what makes the state observable.
package health

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
)

// draining is set when the process starts shutting down.
//
// It is kept here rather than read back from core: BeginShutdown sets a flag on
// the Application, but nothing exports it, and one host wanting to know is not
// yet a reason to widen that interface.
var draining atomic.Bool

// BeginDraining records that shutdown has started, so readiness fails from now
// on. It is called with BeginShutdown, before anything is taken apart.
func BeginDraining() { draining.Store(true) }

// Draining reports whether shutdown has begun.
func Draining() bool { return draining.Load() }

// Check is one dependency and what asking it produced.
type Check struct {
	Name string `json:"name"`
	OK   bool   `json:"ok"`
	Err  string `json:"error,omitempty"`
}

// Ready asks every dependency this process cannot serve a request without.
//
// The queue is deliberately absent. Nothing on AdapterQueue answers "are you
// reachable" without publishing something, the memory backend cannot fail, and
// a queue that is down degrades logging rather than stopping requests - which
// is a reason to alert, not a reason to leave the load balancer pool.
func Ready(ctx context.Context) []Check {
	return []Check{
		safely("database", func() error { return pingDB(ctx) }),
		safely("cache", probeCache),
	}
}

// safely turns a panic into a failed check.
//
// Not defensive habit: the accessors hand back wrappers, not the resources.
// sdk.Runtime.GetCacheAdapter builds a runtime.Cache around whatever is
// configured and returns it even when nothing is - so the value is not nil, the
// cache inside it is, and the first call dereferences it. A nil check cannot
// see that, and the same is true of GetQueueAdapter.
//
// Whatever the reason, a probe is the last thing that should be able to take
// the process down: the caller is asking whether this instance is well, and
// killing it to answer is the wrong reply.
func safely(name string, fn func() error) (c Check) {
	c = Check{Name: name}
	defer func() {
		if r := recover(); r != nil {
			c.OK, c.Err = false, fmt.Sprintf("the check panicked: %v", r)
		}
	}()
	if err := fn(); err != nil {
		c.Err = err.Error()
		return c
	}
	c.OK = true
	return c
}

// Healthy reports whether every check passed.
func Healthy(checks []Check) bool {
	for _, c := range checks {
		if !c.OK {
			return false
		}
	}
	return true
}

func pingDB(ctx context.Context) error {
	db := sdk.Runtime.GetDb()
	if db == nil {
		return errors.New("no database configured")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// cacheProbePrefix names the probe's keys. The key itself is per probe, not
// fixed: two /ready requests arriving together - or two instances sharing one
// redis, which is the normal deployment - would otherwise overwrite each
// other's value between the write and the read and each conclude the cache was
// broken. A readiness probe that reports false negatives under load takes
// healthy instances out of the pool, which is worse than not probing.
const cacheProbePrefix = "go-admin:health:"

// cacheProbeTTL is short because these keys are write-once and never read
// again by anyone else; it only has to outlive the read that follows.
const cacheProbeTTL = 30

func probeCache() error {
	adapter := sdk.Runtime.GetCacheAdapter()
	if adapter == nil {
		return errors.New("no cache configured")
	}

	suffix := make([]byte, 8)
	if _, err := rand.Read(suffix); err != nil {
		return fmt.Errorf("could not build a probe key: %w", err)
	}
	key := cacheProbePrefix + hex.EncodeToString(suffix)

	// Written and read back rather than only read: a cache that answers "miss"
	// for every key - a client pointed at the wrong server - is
	// indistinguishable from a healthy one on a read alone.
	want := time.Now().Format(time.RFC3339Nano)
	if err := adapter.Set(key, want, cacheProbeTTL); err != nil {
		return err
	}
	// Best effort, and its error is deliberately dropped: the verdict is
	// already decided by the read below, and a cache that cannot delete a key
	// it just wrote is not a reason to refuse traffic. The TTL is the real
	// cleanup.
	defer func() { _ = adapter.Del(key) }()

	got, err := adapter.Get(key)
	if err != nil {
		return err
	}
	if got != want {
		return errors.New("the cache returned a different value than was written")
	}
	return nil
}
