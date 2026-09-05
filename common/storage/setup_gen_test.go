package storage

import (
	"testing"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/runtime"
)

// Issue #892: a configuration reload replaces the queue adapter and the
// consumers registered against the previous one are attached to a queue nobody
// publishes to any more.
//
// The fix has two halves. attachQueueConsumers gives a new queue its own
// consumers and the same queue none, which cmd/api covers against a queue the
// test controls. This is the other half: that a reload actually produces a new
// queue for it to notice. Setup is what config re-runs on every change, so
// calling it twice is what a reload does to this package.
func TestSetupBumpsTheQueueGenerationOnEveryReload(t *testing.T) {
	// Setup writes the process-wide sdk.Runtime - the cache and queue adapters -
	// and this package's own record of what it installed. Restoring all of it
	// keeps the test from deciding what a later test in this binary sees,
	// which is the same isolation cmd/api's freshRuntime provides.
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
	config.QueueConfig = &config.Queue{Memory: &config.QueueMemory{PoolSize: 10}}

	before := QueueGeneration()
	Setup()
	first := QueueGeneration()
	Setup()
	second := QueueGeneration()

	t.Logf("before=%d first=%d second=%d", before, first, second)
	if first == before {
		t.Fatal("the first Setup did not install a queue")
	}
	if second == first {
		t.Fatal("a second Setup - which is what a configuration reload does - did not install a new one")
	}
}
