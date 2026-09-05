/*
 * @Author: zhangwenjian
 * @Date: 2025/04/13 22:03
 * @Last Modified by: zhangwenjian
 * @Last Modified time: 2025/04/13 22:03
 */

package storage

import (
	"log"
	"sync"

	"github.com/go-admin-team/go-admin-core/v2/captcha"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
)

// Setup 配置storage组件
func Setup() {
	setupCache()
	setupCaptcha()
	setupQueue()
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
)

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

	// Only an adapter this package installed. GetQueueAdapter never returns
	// nil - with nothing configured the runtime falls back to its own memory
	// queue and wraps that - so the `if q := GetQueueAdapter(); q != nil` this
	// replaces was always true, and shut down the fallback queue on the very
	// first start, before anything had used it.
	if installed != nil {
		installed.Shutdown()
	}

	queueAdapter, err := config.QueueConfig.Setup()
	if err != nil {
		log.Fatalf("queue setup error, %s\n", err.Error())
	}
	sdk.Runtime.SetQueueAdapter(queueAdapter)
	installed = queueAdapter
	installedGen++

	// Deliberately not started here. Run has to come after the consumers have
	// registered: the contract implementations refuse a registration once the
	// queue is running (storage.ErrQueueAlreadyStarted), and the legacy
	// adapter this repository still goes through swallows that error rather
	// than reporting it - its own comment says the interface gives it no way
	// to tell the caller. Starting here and registering afterwards is
	// therefore a race that loses consumers in silence. Whoever registers is
	// the one that starts it.
}
