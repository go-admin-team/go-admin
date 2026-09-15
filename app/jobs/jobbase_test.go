package jobs

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg/cronjob"
)

// The scheduler had never been stopped. `defer crontab.Stop()` sat directly
// above a `select {}` that never returned, so the deferred call was
// unreachable for the life of the process.
//
// There is one test rather than several because BeforeExit closes to further
// registration once it has run: a second RunShutdown in this binary would find
// an empty registry and pass while proving nothing.
func TestTheSchedulerIsStoppedOnTheWayOut(t *testing.T) {
	var ticks atomic.Int64

	c := cronjob.NewWithSeconds()
	if _, err := c.AddFunc("* * * * * *", func() { ticks.Add(1) }); err != nil {
		t.Fatalf("AddFunc: %v", err)
	}

	startCrontab(c)

	// It has to be running before stopping it can mean anything.
	deadline := time.Now().Add(5 * time.Second)
	for ticks.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if ticks.Load() == 0 {
		t.Fatal("the scheduler never ran the job, so this test cannot show it was stopped")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sdk.Runtime.RunShutdown(ctx); err != nil {
		t.Fatalf("RunShutdown: %v", err)
	}

	// Two and a half seconds is two more firings of a job that runs every
	// second, so silence here is the assertion.
	at := ticks.Load()
	time.Sleep(2500 * time.Millisecond)
	if n := ticks.Load() - at; n > 0 {
		t.Errorf("the job fired %d more times after shutdown: the scheduler is still running", n)
	}
}
