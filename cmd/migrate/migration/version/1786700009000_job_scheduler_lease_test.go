package version

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	jobmodels "go-admin/app/jobs/models"
	common "go-admin/common/models"
)

func openJobLeaseDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&common.Migration{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

// The runtime has no insert path - two instances starting together would
// race to create the row they are both trying to claim - so the row has to
// exist when the migration finishes or nothing ever schedules anything.
func TestTheSchedulerLeaseMigrationLeavesExactlyOneFreeRow(t *testing.T) {
	db := openJobLeaseDB(t)

	if err := _1786700009000JobSchedulerLease(db, "1786700009000"); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if !db.Migrator().HasTable(&jobmodels.SysJobLease{}) {
		t.Fatal("sys_job_lease was not created")
	}

	var rows []jobmodels.SysJobLease
	if err := db.Find(&rows).Error; err != nil {
		t.Fatalf("reading sys_job_lease: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("sys_job_lease holds %d rows, want exactly 1", len(rows))
	}

	row := rows[0]
	if row.Name != jobmodels.SchedulerLeaseName {
		t.Errorf("the seeded row is named %q, want %q; acquire looks the row up by this name and would find nothing",
			row.Name, jobmodels.SchedulerLeaseName)
	}
	if row.Owner != "" {
		t.Errorf("the seeded lease is owned by %q; a fresh install would wait out a TTL held by nobody", row.Owner)
	}
	// Zero, not "now": the take is `expires_at_ms <= now`, so a seeded
	// expiry in the future is a scheduler that does not start until it
	// passes.
	if row.ExpiresAtMs != 0 {
		t.Errorf("the seeded lease expires at %d, want 0", row.ExpiresAtMs)
	}
}

func TestTheSchedulerLeaseMigrationRecordsItsVersion(t *testing.T) {
	db := openJobLeaseDB(t)

	if err := _1786700009000JobSchedulerLease(db, "1786700009000"); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	var got common.Migration
	if err := db.Where("version = ?", "1786700009000").First(&got).Error; err != nil {
		t.Fatalf("the migration did not record its version, so it would run again on every start: %v", err)
	}
}
