package jobs

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	glebarez "github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg/cronjob"
	"gorm.io/gorm"

	models2 "go-admin/app/jobs/models"
)

// The scheduler had never been stopped. `defer crontab.Stop()` sat directly
// above a `select {}` that never returned, so the deferred call was
// unreachable for the life of the process.
//
// It now goes through the supervisor, which is what production does and what
// owns the shutdown callback since the lease landed (#915): a scheduler stops
// either because the process is going down or because this instance lost the
// lease, and only the supervisor can tell those apart.
//
// There is one test rather than several because BeforeExit closes to further
// registration once it has run: a second RunShutdown in this binary would
// find an empty registry and pass while proving nothing. The lease-release
// assertion is folded in here for the same reason.
func TestTheSchedulerIsStoppedAndTheLeaseHandedBackOnTheWayOut(t *testing.T) {
	const tenant = "*"

	db := leaseDB(t)
	var ticks atomic.Int64

	c := cronjob.NewWithSeconds()
	if _, err := c.AddFunc("* * * * * *", func() { ticks.Add(1) }); err != nil {
		t.Fatalf("AddFunc: %v", err)
	}
	sdk.Runtime.SetCrontabByTenant(tenant, c)

	s := newSupervisor(tenant, db, "instance-under-test")
	s.start()

	if !s.holdsLease() {
		t.Fatal("the supervisor did not take a free lease, so this test would prove nothing about giving it back")
	}

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

	// And the lease is free, so a successor takes it immediately instead of
	// waiting out a TTL held by a process that has exited.
	var row models2.SysJobLease
	if err := db.Where("name = ?", models2.SchedulerLeaseName).First(&row).Error; err != nil {
		t.Fatalf("reading the lease row: %v", err)
	}
	if row.Owner != "" {
		t.Errorf("the lease is still owned by %q after shutdown; a successor would wait out the TTL", row.Owner)
	}
}

// leaseDB is a database with the two tables jobs.setup touches and one free
// lease row, which is the shape migration 1786700009000 leaves behind.
func leaseDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(glebarez.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models2.SysJob{}, &models2.SysJobLease{}); err != nil {
		t.Fatalf("migrating: %v", err)
	}
	row := models2.SysJobLease{Name: models2.SchedulerLeaseName}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("seeding the lease row: %v", err)
	}
	return db
}
