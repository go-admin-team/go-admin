package migration

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	contractmigration "github.com/go-admin-team/go-admin-core/v2/sdk/contract/migration"

	common "go-admin/common/models"
)

// withContractRegistry points contractSnapshot at an isolated
// *contractmigration.Registry for the duration of one test, instead of
// go-admin-core's single process-wide one - see contractSnapshot's doc
// comment for why that indirection exists. Restored on cleanup so other
// tests in this package keep seeing an empty contract registry regardless of
// run order.
func withContractRegistry(t *testing.T) *contractmigration.Registry {
	t.Helper()
	reg := contractmigration.NewRegistry()
	orig := contractSnapshot
	contractSnapshot = reg.Snapshot
	t.Cleanup(func() { contractSnapshot = orig })
	return reg
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err = db.AutoMigrate(&common.Migration{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

// recordFor is what an app's migration is expected to do: write its own
// completion row, with the version it was handed and the app code it was told
// it belongs to.
func recordFor(db *gorm.DB, version, appCode string) error {
	return db.Create(&common.Migration{Version: version, AppCode: appCode}).Error
}

func rowsByVersion(t *testing.T, db *gorm.DB) map[string]common.Migration {
	t.Helper()
	var rows []common.Migration
	if err := db.Find(&rows).Error; err != nil {
		t.Fatalf("read sys_migration: %v", err)
	}
	out := make(map[string]common.Migration, len(rows))
	for _, r := range rows {
		out[r.Version] = r
	}
	return out
}

// Acceptance 9: a migration registered through ForApp("x") lands in
// sys_migration with app_code "x".
//
// The registry cannot write that row for the migration, because the row is the
// migration's own last statement inside its own transaction. So the only thing
// that can make this true is handing the code to the function - which is why
// AppMigrationFunc takes three parameters.
func TestForAppRecordsItsAppCode(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	m.ForApp("x").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		return recordFor(db, version, appCode)
	})
	m.Migrate()

	rows := rowsByVersion(t, db)
	row, ok := rows["x-1786800001000"]
	if !ok {
		t.Fatalf("no row for x-1786800001000; got %v", rows)
	}
	if row.AppCode != "x" {
		t.Errorf("app_code = %q, want %q", row.AppCode, "x")
	}
}

// The framework path is untouched: same signature, and an empty app code, which
// is what the column defaults to and what every row written before this field
// existed reads back as.
func TestSetVersionStillRecordsTheFrameworkAsEmpty(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error {
		return db.Create(&common.Migration{Version: version}).Error
	})
	m.Migrate()

	rows := rowsByVersion(t, db)
	row, ok := rows["1786700009000"]
	if !ok {
		t.Fatalf("no row for 1786700009000; got %v", rows)
	}
	if row.AppCode != "" {
		t.Errorf("app_code = %q, want empty (framework)", row.AppCode)
	}
}

// Acceptance 12: --app x runs x's migrations and touches nothing else.
func TestMigrateAppRunsOnlyThatApp(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	ran := map[string]bool{}
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error {
		ran["core"] = true
		return db.Create(&common.Migration{Version: version}).Error
	})
	m.ForApp("x").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		ran["x"] = true
		return recordFor(db, version, appCode)
	})
	m.ForApp("y").SetVersion("1786800002000", func(db *gorm.DB, version, appCode string) error {
		ran["y"] = true
		return recordFor(db, version, appCode)
	})

	m.MigrateApp("x")

	if !ran["x"] {
		t.Error("x did not run")
	}
	if ran["y"] || ran["core"] {
		t.Errorf("MigrateApp(x) also ran %v", ran)
	}
	rows := rowsByVersion(t, db)
	if len(rows) != 1 {
		t.Fatalf("sys_migration has %d rows, want 1: %v", len(rows), rows)
	}
}

// "core" is what status prints for the framework, so --app core has to select
// it. The stored code is the empty string; AppFilter is the translation.
func TestMigrateAppCoreSelectsTheFramework(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	ran := map[string]bool{}
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error {
		ran["core"] = true
		return db.Create(&common.Migration{Version: version}).Error
	})
	m.ForApp("x").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		ran["x"] = true
		return recordFor(db, version, appCode)
	})

	m.MigrateApp(FrameworkAppCode)

	if !ran["core"] {
		t.Error("framework migration did not run")
	}
	if ran["x"] {
		t.Error("--app core also ran x")
	}
}

// Zero-argument Migrate keeps meaning "everything", which is what every
// existing caller relies on.
func TestMigrateRunsEveryApp(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	var order []string
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error {
		order = append(order, version)
		return db.Create(&common.Migration{Version: version}).Error
	})
	m.ForApp("bbb").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		order = append(order, version)
		return recordFor(db, version, appCode)
	})
	m.ForApp("aaa").SetVersion("1786800002000", func(db *gorm.DB, version, appCode string) error {
		order = append(order, version)
		return recordFor(db, version, appCode)
	})

	m.Migrate()

	// Namespacing puts every framework migration - bare digits - ahead of every
	// app migration, and orders apps by code rather than by whose timestamp
	// happened to be smaller. aaa's file is the newer of the two and still runs
	// first. Cross-app order is not promised, but this is the order, and it is
	// the one to notice changed.
	want := []string{"1786700009000", "aaa-1786800002000", "bbb-1786800001000"}
	if len(order) != len(want) {
		t.Fatalf("ran %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("ran %v, want %v", order, want)
		}
	}
}

// Two apps minting the same millisecond timestamp used to mean one of them was
// read as already applied and silently skipped. The namespace prefix is what
// makes that impossible without changing the primary key.
func TestNamespacingKeepsTwoAppsWithTheSameTimestampApart(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	const sameTimestamp = "1786800001000"
	ran := 0
	for _, app := range []string{"crm", "oms"} {
		m.ForApp(app).SetVersion(sameTimestamp, func(db *gorm.DB, version, appCode string) error {
			ran++
			return recordFor(db, version, appCode)
		})
	}
	m.Migrate()

	if ran != 2 {
		t.Errorf("ran %d migrations, want 2", ran)
	}
	rows := rowsByVersion(t, db)
	for _, want := range []string{"crm-" + sameTimestamp, "oms-" + sameTimestamp} {
		if _, ok := rows[want]; !ok {
			t.Errorf("missing %s; got %v", want, rows)
		}
	}
}

func TestNamespacedKeyLeavesFrameworkVersionsBare(t *testing.T) {
	if got := namespacedKey("", "1786700009000"); got != "1786700009000" {
		t.Errorf("framework version was rewritten to %q", got)
	}
	if got := namespacedKey("crm", "1786800001000"); got != "crm-1786800001000" {
		t.Errorf("namespacedKey = %q", got)
	}
}

// An app code differing only in case would group as two apps in status and sort
// before every lower-case one, for no reason a reader could guess.
func TestForAppNormalisesTheCode(t *testing.T) {
	m := newMigration()
	if got := m.ForApp("  CRM  ").AppCode(); got != "crm" {
		t.Errorf("AppCode = %q, want crm", got)
	}
}

func TestForAppRejectsReservedCodes(t *testing.T) {
	for _, code := range []string{"", "   ", FrameworkAppCode, "CORE"} {
		t.Run("code="+code, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("ForApp(%q) did not panic", code)
				}
			}()
			newMigration().ForApp(code)
		})
	}
}

func TestStatusReportsPendingAppliedAndOrphaned(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	applied := time.Date(2026, 8, 25, 14, 3, 11, 0, time.UTC)
	if err := db.Create(&common.Migration{Version: "1786700009000", ApplyTime: applied}).Error; err != nil {
		t.Fatal(err)
	}
	// Recorded, but nothing registers it any more.
	if err := db.Create(&common.Migration{Version: "gone-1786800000000", ApplyTime: applied, AppCode: "gone"}).Error; err != nil {
		t.Fatal(err)
	}
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error { return nil })
	m.ForApp("crm").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error { return nil })

	entries, err := m.Status()
	if err != nil {
		t.Fatal(err)
	}
	byVersion := map[string]StatusEntry{}
	for _, e := range entries {
		byVersion[e.Version] = e
	}

	if e := byVersion["1786700009000"]; !e.Applied || !e.Registered || e.AppCode != "" {
		t.Errorf("framework entry = %+v", e)
	} else if e.ApplyTime == nil || !e.ApplyTime.Equal(applied) {
		t.Errorf("framework apply time = %v, want %v", e.ApplyTime, applied)
	}
	if e := byVersion["crm-1786800001000"]; e.Applied || !e.Registered || e.AppCode != "crm" {
		t.Errorf("crm entry = %+v", e)
	}
	if e := byVersion["gone-1786800000000"]; !e.Applied || e.Registered || e.AppCode != "gone" {
		t.Errorf("orphaned entry = %+v", e)
	}
}

// Acceptance 11 rests on this: status and --dry-run both go through Status, and
// Status must not create the table it reads.
func TestStatusDoesNotCreateItsTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	m := newMigration()
	m.SetDb(db)
	m.ForApp("crm").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error { return nil })

	entries, err := m.Status()
	if err != nil {
		t.Fatalf("Status on a database with no sys_migration: %v", err)
	}
	if len(entries) != 1 || entries[0].Applied {
		t.Errorf("entries = %+v, want one pending", entries)
	}
	if db.Migrator().HasTable(&common.Migration{}) {
		t.Error("Status created sys_migration; it must only read")
	}
}

// The completion row is the migration's own last statement, inside its own
// transaction. A migration that fails must leave no record of having run, or
// the next run skips it and the schema stays half-changed with nothing to say
// so.
func TestFailedMigrationLeavesNoRecord(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	m.ForApp("crm").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		return db.Transaction(func(tx *gorm.DB) error {
			if err := recordFor(tx, version, appCode); err != nil {
				return err
			}
			return errTestMigrationFailed
		})
	})

	// run() calls log.Fatal on failure, which would take the test binary with
	// it, so drive the registered function directly - the point here is the
	// transaction boundary, not the scheduler.
	entry := m.version["crm-1786800001000"]
	if err := entry.fn(db, "crm-1786800001000"); err == nil {
		t.Fatal("migration reported success")
	}
	if rows := rowsByVersion(t, db); len(rows) != 0 {
		t.Errorf("sys_migration has %v after a failed migration", rows)
	}
}

var errTestMigrationFailed = &testError{"boom"}

type testError struct{ s string }

func (e *testError) Error() string { return e.s }

// A mistyped --app used to select nothing and print "no migrations to apply",
// which reads as "already up to date" - the command reports success and does
// nothing, which is the failure mode this whole batch exists to remove.
func TestMigrateAppOnAnUnknownCodeSaysSo(t *testing.T) {
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error {
		return db.Create(&common.Migration{Version: version}).Error
	})
	m.ForApp("crm").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		return recordFor(db, version, appCode)
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	m.MigrateApp("crmm")

	if !strings.Contains(buf.String(), `no migrations are registered for app "crmm"`) {
		t.Errorf("output = %q", buf.String())
	}
	if !strings.Contains(buf.String(), "registered: core, crm") {
		t.Errorf("the message must list what is registered; got %q", buf.String())
	}
	if rows := rowsByVersion(t, db); len(rows) != 0 {
		t.Errorf("a typo ran %v", rows)
	}
}

// This is the acceptance test for PRD 006's host-wiring gap: a migration
// registered through contract/migration.ForApp - the only door open to a
// third-party application - must actually run, be recorded under its app
// code, and show up in AppCodes/Status/--app the same as one registered
// through the host's own m.ForApp. Before mergedEntries existed, m.Migrate()
// never looked at contract/migration's registry at all, so this compiled,
// registered, and silently never ran.
func TestMergedEntriesRunsAContractRegisteredAppMigration(t *testing.T) {
	reg := withContractRegistry(t)
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	ran := false
	reg.ForApp("order").SetVersion("1793800000000", func(db *gorm.DB, version, appCode string) error {
		ran = true
		return recordFor(db, version, appCode)
	})

	m.Migrate()

	if !ran {
		t.Fatal("contract-registered migration did not run")
	}
	rows := rowsByVersion(t, db)
	row, ok := rows["order-1793800000000"]
	if !ok {
		t.Fatalf("no row for order-1793800000000; got %v", rows)
	}
	if row.AppCode != "order" {
		t.Errorf("app_code = %q, want %q", row.AppCode, "order")
	}
}

// migrate status and --dry-run both read Status; a contract-registered
// migration has to appear there under its app code exactly like a
// host-registered one, both before and after it is applied.
func TestMergedEntriesStatusIncludesContractRegisteredMigrations(t *testing.T) {
	reg := withContractRegistry(t)
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	reg.ForApp("order").SetVersion("1793800000000", func(db *gorm.DB, version, appCode string) error {
		return recordFor(db, version, appCode)
	})

	entries, err := m.Status()
	if err != nil {
		t.Fatal(err)
	}
	byVersion := map[string]StatusEntry{}
	for _, e := range entries {
		byVersion[e.Version] = e
	}
	e, ok := byVersion["order-1793800000000"]
	if !ok || !e.Registered || e.Applied || e.AppCode != "order" {
		t.Fatalf("pending contract entry = %+v (ok=%v)", e, ok)
	}

	m.Migrate()

	entries, err = m.Status()
	if err != nil {
		t.Fatal(err)
	}
	byVersion = map[string]StatusEntry{}
	for _, e := range entries {
		byVersion[e.Version] = e
	}
	if e := byVersion["order-1793800000000"]; !e.Applied {
		t.Errorf("applied contract entry = %+v", e)
	}
}

// AppCodes feeds both --app's typo detection (appRegistrationError) and the
// group headings status prints; a contract-registered app has to appear
// there or a real "go-admin migrate --app order" would be told the app does
// not exist.
func TestMergedEntriesAppCodesIncludesContractRegisteredApps(t *testing.T) {
	reg := withContractRegistry(t)
	m := newMigration()
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error { return nil })
	reg.ForApp("order").SetVersion("1793800000000", func(db *gorm.DB, version, appCode string) error { return nil })

	got := m.AppCodes()
	want := []string{"core", "order"}
	if len(got) != len(want) {
		t.Fatalf("AppCodes = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AppCodes = %v, want %v", got, want)
		}
	}
}

// --app order has to actually run only order's migrations - the same
// per-app isolation MigrateApp already gives host-registered apps - even
// though order is registered in a different registry entirely.
func TestMergedEntriesMigrateAppRunsOnlyThatContractApp(t *testing.T) {
	reg := withContractRegistry(t)
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	ran := map[string]bool{}
	m.SetVersion("1786700009000", func(db *gorm.DB, version string) error {
		ran["core"] = true
		return db.Create(&common.Migration{Version: version}).Error
	})
	reg.ForApp("order").SetVersion("1793800000000", func(db *gorm.DB, version, appCode string) error {
		ran["order"] = true
		return recordFor(db, version, appCode)
	})

	m.MigrateApp("order")

	if !ran["order"] {
		t.Error("order did not run")
	}
	if ran["core"] {
		t.Errorf("MigrateApp(order) also ran %v", ran)
	}
}

// A host-registered key is not supposed to collide with a namespaced
// contract key (see mergedEntries' doc comment), but if it somehow did, the
// host's own registration must win rather than a third-party application
// silently overwriting a framework migration under the same key.
func TestMergedEntriesHostRegistrationWinsOnKeyCollision(t *testing.T) {
	reg := withContractRegistry(t)
	db := newTestDB(t)
	m := newMigration()
	m.SetDb(db)

	hostRan, contractRan := false, false
	m.ForApp("dup").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		hostRan = true
		return recordFor(db, version, appCode)
	})
	reg.ForApp("dup").SetVersion("1786800001000", func(db *gorm.DB, version, appCode string) error {
		contractRan = true
		return recordFor(db, version, appCode)
	})

	m.Migrate()

	if !hostRan {
		t.Error("host registration did not run")
	}
	if contractRan {
		t.Error("contract registration ran; host registration should have won the collision")
	}
}

// GetFilename must stay the same rule the contract package applies, since an
// application registering through contract/migration names its files by that
// convention and has to land on the same version string. Pinning the reject
// case is what catches a re-divergence: a local copy that only sliced would
// return "add_orders.go" here and register a migration under a key that never
// matches anything.
func TestGetFilenameDelegatesToTheContractRule(t *testing.T) {
	if got := GetFilename("version/1786700001000_demo_menu.go"); got != "1786700001000" {
		t.Fatalf("GetFilename = %q, want %q", got, "1786700001000")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("a file name carrying no version did not panic")
		}
	}()
	GetFilename("version/add_orders.go")
}
