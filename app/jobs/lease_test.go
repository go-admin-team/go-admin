package jobs

import (
	"os"
	"testing"
	"time"

	glebarez "github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	models2 "go-admin/app/jobs/models"
)

// The lease is the one thing in this package whose correctness is a property
// of the database rather than of this process, so these run against every
// dialect that can be reached. SQLite always; the others when their DSN is
// set, and they must be set in CI - a suite that quietly skipped them would
// report success for a lease that cannot be taken at all on the dialect most
// installations actually run.
const (
	mysqlDSNEnv     = "GO_ADMIN_TEST_MYSQL_DSN"
	postgresDSNEnv  = "GO_ADMIN_TEST_POSTGRES_DSN"
	sqlserverDSNEnv = "GO_ADMIN_TEST_SQLSERVER_DSN"
)

type dialectDB struct {
	name string
	open func(string) gorm.Dialector
	env  string
}

var optionalDialects = []dialectDB{
	{"mysql", func(dsn string) gorm.Dialector { return mysql.Open(dsn) }, mysqlDSNEnv},
	{"postgres", func(dsn string) gorm.Dialector { return postgres.Open(dsn) }, postgresDSNEnv},
	{"sqlserver", func(dsn string) gorm.Dialector { return sqlserver.Open(dsn) }, sqlserverDSNEnv},
}

// eachDialect runs body against SQLite and against every optional dialect
// whose DSN is set.
func eachDialect(t *testing.T, body func(t *testing.T, db *gorm.DB)) {
	t.Helper()

	t.Run("sqlite", func(t *testing.T) {
		db, err := gorm.Open(glebarez.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatalf("opening sqlite: %v", err)
		}
		body(t, seedLeaseTable(t, db))
	})

	for _, d := range optionalDialects {
		t.Run(d.name, func(t *testing.T) {
			dsn := os.Getenv(d.env)
			if dsn == "" {
				if os.Getenv("CI") != "" {
					t.Fatalf("%s is not set while CI is: the lease must not go untested on %s", d.env, d.name)
				}
				t.Skipf("%s is not set; skipping %s", d.env, d.name)
			}
			db, err := gorm.Open(d.open(dsn), &gorm.Config{})
			if err != nil {
				t.Fatalf("connecting to %s: %v", d.env, err)
			}
			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("sql.DB: %v", err)
			}
			t.Cleanup(func() { _ = sqlDB.Close() })
			body(t, seedLeaseTable(t, db))
		})
	}
}

// seedLeaseTable builds the shape 1786700009000 leaves behind: the table,
// and exactly one free row.
func seedLeaseTable(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()

	if err := db.Migrator().DropTable(&models2.SysJobLease{}); err != nil {
		t.Fatalf("dropping sys_job_lease: %v", err)
	}
	if err := db.AutoMigrate(&models2.SysJobLease{}); err != nil {
		t.Fatalf("creating sys_job_lease: %v", err)
	}
	row := models2.SysJobLease{Name: models2.SchedulerLeaseName, AcquiredAtMs: 0, ExpiresAtMs: 0}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("seeding the lease row: %v", err)
	}
	return db
}

func TestTheDatabaseClockIsReadableAndIsNotTheZeroTime(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		nowMs, err := dbNowMs(db)
		if err != nil {
			t.Fatalf("dbNowMs: %v", err)
		}
		if nowMs <= 0 {
			t.Fatal("the database clock read as zero, which would read every lease as expired")
		}
		// Not an assertion about either clock's accuracy - a container's
		// clock and this one can drift. An hour is far wider than drift
		// and far narrower than a timezone offset, which is the mistake
		// this catches: reading MySQL's UTC_TIMESTAMP over a loc=Local
		// DSN lands exactly one zone offset away and is invisible to
		// every assertion that only compares the lease against itself.
		drift := time.Duration(time.Now().UnixMilli()-nowMs) * time.Millisecond
		if drift > time.Hour || drift < -time.Hour {
			t.Errorf("the database clock is %v away from this process's; a timezone mistake looks exactly like this", drift)
		}
	})
}

func TestOnlyOneOfTwoInstancesTakesTheLease(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		a := &lease{db: db, owner: "instance-a", ttl: time.Minute}
		b := &lease{db: db, owner: "instance-b", ttl: time.Minute}

		heldA, err := a.acquire()
		if err != nil {
			t.Fatalf("a.acquire: %v", err)
		}
		if !heldA {
			t.Fatal("the first instance did not take a free lease")
		}

		heldB, err := b.acquire()
		if err != nil {
			t.Fatalf("b.acquire: %v", err)
		}
		if heldB {
			t.Error("the second instance took a lease the first one holds: both would schedule")
		}
	})
}

func TestTheHolderRenewsAndTheOtherStillCannotTakeIt(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		a := &lease{db: db, owner: "instance-a", ttl: time.Minute}
		b := &lease{db: db, owner: "instance-b", ttl: time.Minute}

		if held, err := a.acquire(); err != nil || !held {
			t.Fatalf("a.acquire: held=%v err=%v", held, err)
		}
		before := readLease(t, db)

		if held, err := a.acquire(); err != nil || !held {
			t.Fatalf("a renewing: held=%v err=%v", held, err)
		}
		after := readLease(t, db)

		if after.ExpiresAtMs < before.ExpiresAtMs {
			t.Errorf("renewal moved the expiry backwards: %d then %d", before.ExpiresAtMs, after.ExpiresAtMs)
		}
		if after.AcquiredAtMs != before.AcquiredAtMs {
			t.Errorf("renewal moved acquired_at_ms (%d then %d); it must say when the lease was taken, not when it was last renewed",
				before.AcquiredAtMs, after.AcquiredAtMs)
		}
		if held, err := b.acquire(); err != nil || held {
			t.Errorf("the other instance took a renewed lease: held=%v err=%v", held, err)
		}
	})
}

func TestAnExpiredLeaseIsTakenOver(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		a := &lease{db: db, owner: "instance-a", ttl: time.Minute}
		b := &lease{db: db, owner: "instance-b", ttl: time.Minute}

		if held, err := a.acquire(); err != nil || !held {
			t.Fatalf("a.acquire: held=%v err=%v", held, err)
		}

		// What a dead leader leaves behind: its row, unrenewed, past its
		// expiry. Forced rather than waited out, so the test does not
		// trade a second of sleep for the same assertion.
		expire(t, db)

		if held, err := b.acquire(); err != nil || !held {
			t.Fatalf("the successor did not take an expired lease: held=%v err=%v", held, err)
		}
		if got := readLease(t, db).Owner; got != "instance-b" {
			t.Errorf("owner is %q after takeover, want instance-b", got)
		}

		// And the instance that lost it must not get it back by renewing:
		// renewal is scoped to the owner column it no longer matches.
		if held, err := a.acquire(); err != nil || held {
			t.Errorf("the dead leader renewed a lease it had lost: held=%v err=%v", held, err)
		}
	})
}

func TestReleaseHandsTheLeaseOnWithoutWaitingOutTheTTL(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		a := &lease{db: db, owner: "instance-a", ttl: time.Hour}
		b := &lease{db: db, owner: "instance-b", ttl: time.Minute}

		if held, err := a.acquire(); err != nil || !held {
			t.Fatalf("a.acquire: held=%v err=%v", held, err)
		}
		if held, err := b.acquire(); err != nil || held {
			t.Fatalf("precondition: b must not hold it yet (held=%v err=%v)", held, err)
		}

		if err := a.release(); err != nil {
			t.Fatalf("a.release: %v", err)
		}
		if held, err := b.acquire(); err != nil || !held {
			t.Errorf("a released a lease with an hour left and the successor still could not take it: held=%v err=%v", held, err)
		}
	})
}

func TestReleasingALeaseSomebodyElseHoldsDoesNothing(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		a := &lease{db: db, owner: "instance-a", ttl: time.Minute}
		stale := &lease{db: db, owner: "instance-gone", ttl: time.Minute}

		if held, err := a.acquire(); err != nil || !held {
			t.Fatalf("a.acquire: held=%v err=%v", held, err)
		}
		if err := stale.release(); err != nil {
			t.Fatalf("stale.release: %v", err)
		}
		if got := readLease(t, db).Owner; got != "instance-a" {
			t.Errorf("owner is %q; an instance that already lost the lease cleared its successor's row", got)
		}
	})
}

func TestAMissingLeaseRowIsReportedRatherThanSilentlyNeverScheduling(t *testing.T) {
	eachDialect(t, func(t *testing.T, db *gorm.DB) {
		if err := db.Where("name = ?", models2.SchedulerLeaseName).
			Delete(&models2.SysJobLease{}).Error; err != nil {
			t.Fatalf("deleting the lease row: %v", err)
		}
		a := &lease{db: db, owner: "instance-a", ttl: time.Minute}
		held, err := a.acquire()
		if held {
			t.Fatal("acquire reported the lease held with no row to hold")
		}
		if err == nil {
			t.Error("a missing lease row was reported as an ordinary 'someone else holds it': jobs would never run anywhere and nothing would say why")
		}
	})
}

func readLease(t *testing.T, db *gorm.DB) models2.SysJobLease {
	t.Helper()
	var row models2.SysJobLease
	if err := db.Where("name = ?", models2.SchedulerLeaseName).First(&row).Error; err != nil {
		t.Fatalf("reading the lease row: %v", err)
	}
	return row
}

func expire(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Model(&models2.SysJobLease{}).
		Where("name = ?", models2.SchedulerLeaseName).
		Update("expires_at_ms", 0).Error; err != nil {
		t.Fatalf("expiring the lease: %v", err)
	}
}
