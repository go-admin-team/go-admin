package storage

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
	corestorage "github.com/go-admin-team/go-admin-core/v2/storage"
)

// countingQueue stands in for an installed adapter. Only Shutdown is exercised
// - the drain callback never publishes or consumes - so the rest of
// AdapterQueue is deliberately absent: `installed` is typed on Shutdown alone,
// and widening the fake would only invite it to be used for something else.
type countingQueue struct {
	calls atomic.Int32
	block chan struct{}

	// started is closed on the way into Shutdown, so a test can wait for the
	// call rather than assume the goroutine that makes it was scheduled. The
	// caller returns on its own deadline while Shutdown is still running, so
	// reading calls straight after that return is a race with the increment.
	startOnce sync.Once
	started   chan struct{}
}

func newCountingQueue() *countingQueue {
	return &countingQueue{started: make(chan struct{})}
}

func (q *countingQueue) Shutdown() {
	q.startOnce.Do(func() { close(q.started) })
	q.calls.Add(1)
	if q.block != nil {
		<-q.block
	}
}

// waitStarted blocks until Shutdown has been entered, or fails the test.
func (q *countingQueue) waitStarted(t *testing.T) {
	t.Helper()
	select {
	case <-q.started:
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown was never called")
	}
}

// isolate gives the test its own runtime and its own view of what this package
// has installed, and puts the process-wide state back afterwards.
//
// Same isolation as TestSetupBumpsTheQueueGenerationOnEveryReload, plus
// drainRegistered: it is what stops a reload registering a second callback, so
// leaving it set would make every later test in this binary see a package that
// has already registered.
func isolate(t *testing.T) {
	t.Helper()

	prevQ, prevC := config.QueueConfig, config.CacheConfig
	prevRuntime := sdk.Runtime
	queueMu.Lock()
	prevInstalled, prevGen, prevRegistered := installed, installedGen, drainRegistered
	queueMu.Unlock()

	t.Cleanup(func() {
		config.QueueConfig, config.CacheConfig = prevQ, prevC
		sdk.Runtime = prevRuntime
		queueMu.Lock()
		installed, installedGen, drainRegistered = prevInstalled, prevGen, prevRegistered
		queueMu.Unlock()
	})

	sdk.Runtime = runtime.NewConfig()
	queueMu.Lock()
	installed, drainRegistered = nil, false
	queueMu.Unlock()
}

// setInstalled puts a fake where setupQueue would have left the real adapter.
//
// Legitimate because the callback reads `installed` when it runs rather than
// capturing it at registration - that is the property that lets a reload
// replace the adapter and still have the right one drained.
func setInstalled(q interface{ Shutdown() }) {
	queueMu.Lock()
	installed = q
	queueMu.Unlock()
}

func currentInstalled() interface{ Shutdown() } {
	queueMu.Lock()
	defer queueMu.Unlock()
	return installed
}

// Issue #911: nothing shut the queue down at exit, so whatever was buffered
// went with the process.
//
// This is the half that matters most - a callback is on BeforeExit and it
// reaches the adapter this package installed. It says nothing about how many
// times the callback was registered; see the test below for that.
func TestSetupPutsTheQueueDrainOnBeforeExit(t *testing.T) {
	isolate(t)
	config.CacheConfig = &config.Cache{Memory: struct{}{}}
	config.QueueConfig = &config.Queue{Memory: &config.QueueMemory{PoolSize: 10}}

	Setup()
	Setup()
	Setup()

	q := newCountingQueue()
	setInstalled(q)

	if err := sdk.Runtime.RunShutdown(context.Background()); err != nil {
		t.Fatalf("RunShutdown: %v", err)
	}

	if got := q.calls.Load(); got != 1 {
		t.Errorf("Shutdown called %d times, want 1 - 0 means nothing registered the drain", got)
	}
}

// Setup is re-run on every configuration change, so registering from it has to
// be guarded: a callback per reload would leave the shutdown phase holding a
// row of identical entries, each timed and each eligible to be named as the one
// that overran the budget.
//
// Counted at the seam rather than through the effect. The test above cannot see
// this - shutdownQueue takes the adapter on its first run, so the second and
// third callbacks find nothing and return, and three registrations produce
// exactly the same observable result as one. That is a good property of the
// callback and a blind spot for any test that goes through it.
func TestSetupRegistersTheDrainOncePerProcessHoweverManyReloads(t *testing.T) {
	isolate(t)

	previous := setShutdown
	t.Cleanup(func() { setShutdown = previous })
	registrations := 0
	setShutdown = func(func(context.Context)) { registrations++ }

	config.CacheConfig = &config.Cache{Memory: struct{}{}}
	config.QueueConfig = &config.Queue{Memory: &config.QueueMemory{PoolSize: 10}}

	Setup()
	Setup()
	Setup()

	if registrations != 1 {
		t.Errorf("three reloads registered the drain %d times, want 1", registrations)
	}
}

// The drain reaches the adapter that is current when the signal arrives, not
// one captured while wiring up. A reload replaces the adapter, and draining the
// one that was installed at start-up would drain something nobody has published
// to since.
func TestTheDrainRunsAgainstTheAdapterInstalledLast(t *testing.T) {
	isolate(t)
	config.CacheConfig = &config.Cache{Memory: struct{}{}}
	config.QueueConfig = &config.Queue{Memory: &config.QueueMemory{PoolSize: 10}}
	Setup()

	first, second := newCountingQueue(), newCountingQueue()
	setInstalled(first)
	setInstalled(second)

	if err := sdk.Runtime.RunShutdown(context.Background()); err != nil {
		t.Fatalf("RunShutdown: %v", err)
	}

	if first.calls.Load() != 0 {
		t.Error("the adapter that was replaced was shut down; the callback captured it instead of " +
			"reading it when it ran")
	}
	if second.calls.Load() != 1 {
		t.Errorf("the current adapter was shut down %d times, want 1", second.calls.Load())
	}
}

// Nothing installed is the shipped default: settings.yml has no queue section,
// so setupQueue returns early and the runtime's own fallback queue is what
// callers get. Shutting that down would close a queue this package neither
// built nor started.
func TestTheDrainDoesNothingWhenThisPackageInstalledNothing(t *testing.T) {
	isolate(t)
	config.CacheConfig = &config.Cache{Memory: struct{}{}}
	config.QueueConfig = &config.Queue{}
	Setup()

	if currentInstalled() != nil {
		t.Fatal("an empty queue configuration installed an adapter, so this test asserts nothing")
	}
	if err := sdk.Runtime.RunShutdown(context.Background()); err != nil {
		t.Fatalf("RunShutdown: %v", err)
	}
}

// The budget bounds the wait, not the work. A consumer that never finishes must
// not hold the process past its grace period - SIGKILL would arrive mid-write
// instead of at a point of the process's choosing.
func TestTheDrainStopsWaitingWhenTheBudgetIsGone(t *testing.T) {
	isolate(t)

	blocked := newCountingQueue()
	blocked.block = make(chan struct{})
	t.Cleanup(func() { close(blocked.block) })
	setInstalled(blocked)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	returned := make(chan struct{})
	go func() {
		defer close(returned)
		shutdownQueue(ctx)
	}()

	select {
	case <-returned:
	case <-time.After(5 * time.Second):
		t.Fatal("shutdownQueue did not return after its context expired - it waits on a Shutdown " +
			"that takes no context, so the wait has to be bounded here")
	}

	// Waited for rather than read straight after the return: Shutdown runs on a
	// goroutine that the caller does not join, so the increment is not ordered
	// against the caller giving up on its deadline.
	blocked.waitStarted(t)
	if blocked.calls.Load() != 1 {
		t.Errorf("Shutdown called %d times, want 1 - the drain has to be attempted even when it "+
			"cannot be waited out", blocked.calls.Load())
	}
}

// After the drain this package owns nothing. A reload arriving mid-shutdown
// then builds a new adapter rather than being handed a closed one as its
// `previous` to shut down again.
func TestTheDrainGivesUpOwnershipOfTheAdapter(t *testing.T) {
	isolate(t)
	setInstalled(newCountingQueue())

	shutdownQueue(context.Background())

	if got := currentInstalled(); got != nil {
		t.Errorf("installed is %T after the drain, want nil", got)
	}
}

// runtimeQueue is a full AdapterQueue, so it can be handed to the runtime
// rather than only to this package's own record of what it installed.
type runtimeQueue struct {
	countingQueue
}

func (q *runtimeQueue) String() string                            { return "runtime-fake" }
func (q *runtimeQueue) Append(corestorage.Messager) error         { return nil }
func (q *runtimeQueue) Register(string, corestorage.ConsumerFunc) {}
func (q *runtimeQueue) Run()                                      {}

// The drain must not reach a queue this package did not install.
//
// sdk.Runtime.GetQueueAdapter never returns nil: with no queue section
// configured it wraps the runtime's own fallback, and the wrapper's Shutdown
// forwards. Reaching for the accessor would therefore look like it worked and
// would close a queue this package neither built nor started - the same shape
// as the `if q := GetQueueAdapter(); q != nil` that setupQueue already had to
// drop.
//
// The previous test's empty-configuration case cannot see this: it only checks
// that RunShutdown returns, which it would either way.
func TestTheDrainNeverReachesTheRuntimesOwnQueue(t *testing.T) {
	isolate(t)

	onTheRuntime := &runtimeQueue{countingQueue: *newCountingQueue()}
	sdk.Runtime.SetQueueAdapter(onTheRuntime)

	if currentInstalled() != nil {
		t.Fatal("this package installed something, so the distinction under test is not set up")
	}
	if sdk.Runtime.GetQueueAdapter() == nil {
		t.Fatal("the accessor returned nil, so it is no longer the trap this guards")
	}

	shutdownQueue(context.Background())

	if got := onTheRuntime.calls.Load(); got != 0 {
		t.Errorf("the runtime's queue was shut down %d times - the drain went through "+
			"GetQueueAdapter instead of the adapter this package installed", got)
	}
}

// A drain that finishes in the same instant the budget expires counts as
// finished.
//
// Both channels are ready when the select runs, and select picks at random
// among ready cases, so a single look reports an overrun for a drain that
// completed - roughly half the times it lands here. The repetition is what
// makes that visible: one iteration passes either way.
func TestATieBetweenTheDeadlineAndTheDrainGoesToTheDrain(t *testing.T) {
	expired, cancel := context.WithCancel(context.Background())
	cancel()
	<-expired.Done()

	done := make(chan struct{})
	close(done)

	for i := 0; i < 1000; i++ {
		if !finishedBeforeDeadline(expired, done) {
			t.Fatalf("iteration %d of 1000: both the deadline and the drain were ready and the "+
				"deadline won - a drain that completed is being reported as an overrun", i)
		}
	}
}

// The other side of it. A drain that really has not finished has to be
// reported, or the warning never fires and the tie-break above has quietly
// turned into "always say it drained".
func TestADrainThatHasNotFinishedIsReportedAsAnOverrun(t *testing.T) {
	expired, cancel := context.WithCancel(context.Background())
	cancel()
	<-expired.Done()

	stillRunning := make(chan struct{}) // never closed

	if finishedBeforeDeadline(expired, stillRunning) {
		t.Error("an unfinished drain was reported as having finished in time")
	}
}
