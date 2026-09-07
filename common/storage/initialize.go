/*
 * @Author: zhangwenjian
 * @Date: 2025/04/13 22:03
 * @Last Modified by: zhangwenjian
 * @Last Modified time: 2025/04/13 22:03
 */

package storage

import (
	"context"
	"log"
	"sync"

	"github.com/go-admin-team/go-admin-core/v2/captcha"
	corelog "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// Setup 配置storage组件
func Setup() {
	setupCache()
	setupCaptcha()
	setupQueue()
	registerQueueDrain()
}

func setupCache() {
	cacheAdapter, err := config.CacheConfig.Setup()
	if err != nil {
		log.Fatalf("cache setup error, %s\n", err.Error())
	}
	sdk.Runtime.SetCacheAdapter(cacheAdapter)
}

func setupCaptcha() {
	captcha.SetStore(captcha.NewCacheStore(sdk.Runtime.GetCacheAdapter(), 600))
}

var (
	queueMu sync.Mutex
	// installed is the adapter setupQueue built, kept so the next reload can
	// shut it down, and counted so a consumer can tell one from the next.
	installed    interface{ Shutdown() }
	installedGen uint64
	// drainRegistered records that the BeforeExit callback is on the runtime,
	// so that a reload does not add another one.
	drainRegistered bool
)

// setShutdown is sdk.Runtime.SetShutdown, indirected so that registering can
// be observed.
//
// It has to be: the runtime does not report how many callbacks a phase holds,
// and shutdownQueue takes the adapter on its first run, so every registration
// after the first returns immediately and changes nothing anybody can see. A
// reload adding one callback per round would therefore be invisible from the
// outside - which is exactly how it would survive.
var setShutdown = func(f func(context.Context)) { sdk.Runtime.SetShutdown(f) }

// registerQueueDrain puts shutdownQueue on the BeforeExit phase, once.
//
// Setup is one of the callbacks bootstrap.SetupConfig re-runs on every
// configuration change, so registering from it without a guard would leave one
// callback per reload - each shutting down the same adapter, each reported
// separately when the budget runs out.
//
// A flag under the existing mutex rather than a sync.Once: the tests in this
// package already save and restore installed and installedGen to keep one test
// from deciding what the next one sees, and a sync.Once cannot be put back.
func registerQueueDrain() {
	queueMu.Lock()
	first := !drainRegistered
	drainRegistered = true
	queueMu.Unlock()

	if first {
		setShutdown(shutdownQueue)
	}
}

// drainedInTime shuts q down and reports whether it finished before ctx expired.
func drainedInTime(ctx context.Context, q interface{ Shutdown() }) bool {
	done := make(chan struct{})
	go func() {
		defer close(done)
		q.Shutdown()
	}()
	return finishedBeforeDeadline(ctx, done)
}

// finishedBeforeDeadline waits for done or for ctx, and resolves a tie in
// favour of done.
//
// The tie is the reason this is a function of its own rather than one select
// inline. Both channels can be ready when the select runs, select picks at
// random among ready cases, and so a single look reports an overrun for a
// drain that completed - about half the times it lands there, which is exactly
// often enough to be dismissed as noise. core's own RunShutdown re-checks for
// this reason.
//
// Taking channels rather than a queue is what makes it testable: a closed done
// and an expired ctx can be handed in together, which is the state a race
// would otherwise have to be caught in.
func finishedBeforeDeadline(ctx context.Context, done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	case <-ctx.Done():
	}

	select {
	case <-done:
		return true
	default:
		return false
	}
}

// shutdownQueue drains the queue this package installed, on the way out.
//
// Nothing used to. core's Memory.Shutdown closes the queue and waits for every
// consumer to finish what it is holding, and the legacy adapter cancels its
// context and closes the underlying queue - but neither ran at exit, so the
// process left with the login log, the operation log and the API sync still
// buffered, and left reporting success.
//
// The adapter is read here rather than captured at registration because a
// reload replaces it. Registration happens once per process; this runs against
// whatever is current when the signal arrives.
//
// Only an adapter this package installed. sdk.Runtime.GetQueueAdapter never
// returns nil - with no queue section configured the runtime wraps its own
// fallback queue - so going through that accessor would shut down a queue this
// package neither built nor started.
//
// The adapter is taken, not read: after this the package owns nothing, so a
// reload arriving mid-shutdown builds a new one instead of being handed a
// closed one to shut down again. Both implementations tolerate a second
// Shutdown, so this is about who owns it rather than about a crash.
func shutdownQueue(ctx context.Context) {
	queueMu.Lock()
	q := installed
	installed = nil
	queueMu.Unlock()

	if q == nil {
		return
	}

	if drainedInTime(ctx, q) {
		corelog.Info("queue: drained")
		return
	}
	// The wait is what the budget bounds, not the work: Shutdown takes no
	// context and is still running in that goroutine. Saying so here names what
	// is being lost, which the generic overrun message cannot.
	corelog.Warnf("queue: the shutdown budget ran out while the queue was still draining - " +
		"whatever it had not delivered goes with the process. Raise extend.shutdown.cleanup " +
		"if this recurs.")
}

// QueueGeneration reports how many times this package has installed a queue
// adapter. It changes every time setupQueue builds a new one, which is on
// every configuration reload, and stays 0 for as long as the configuration has
// no queue section at all - in which case nothing is installed and callers are
// working with the runtime's own fallback queue.
//
// It exists because there is no way to ask for the adapter's identity from the
// outside. sdk.Runtime.GetQueueAdapter and GetQueuePrefix build a fresh
// runtime.Queue wrapper on every call, so comparing what two calls return
// compares two wrappers and never matches, however many times the underlying
// adapter has been replaced. This package creates the adapter, so this is the
// only place that knows. A counter rather than the adapter itself keeps the
// comparison on a uint64: an adapter type that is not comparable would panic
// an `==` between two interface values.
func QueueGeneration() uint64 {
	queueMu.Lock()
	defer queueMu.Unlock()
	return installedGen
}

func setupQueue() {
	if config.QueueConfig.Empty() {
		return
	}

	queueMu.Lock()
	defer queueMu.Unlock()

	queueAdapter, err := config.QueueConfig.Setup()
	if err != nil {
		log.Fatalf("queue setup error, %s\n", err.Error())
	}

	previous := installed
	sdk.Runtime.SetQueueAdapter(queueAdapter)
	installed = queueAdapter
	installedGen++

	// The previous adapter goes down after the new one is installed, not
	// before. Shutdown waits for its consumers to deliver what it still holds,
	// and for that whole wait the runtime would otherwise be handing producers
	// a queue that has stopped accepting: every Append in the window comes back
	// ErrQueueClosed, and both call sites in common/middleware log it. Swapping
	// first leaves no such window - a producer gets the new queue or the old
	// one, and both work.
	//
	// Only an adapter this package installed. GetQueueAdapter never returns
	// nil - with nothing configured the runtime falls back to its own memory
	// queue and wraps that - so the `if q := GetQueueAdapter(); q != nil` this
	// replaces was always true, and shut down the fallback queue on the very
	// first start, before anything had used it.
	if previous != nil {
		previous.Shutdown()
	}

	// Deliberately not started here. Run has to come after the consumers have
	// registered: the contract implementations refuse a registration once the
	// queue is running (storage.ErrQueueAlreadyStarted), and the legacy
	// adapter this repository still goes through swallows that error rather
	// than reporting it - its own comment says the interface gives it no way
	// to tell the caller. Starting here and registering afterwards is
	// therefore a race that loses consumers in silence. Whoever registers is
	// the one that starts it.
}
