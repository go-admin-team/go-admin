package migrate

import (
	"errors"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/app"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
)

func newInstallDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&adminmodels.SysApp{}); err != nil {
		t.Fatalf("automigrate sys_app: %v", err)
	}
	return db
}

// fakeEngine stands in for the migration engine. The real one is a
// package-level singleton with no exported constructor, so a test taking it
// would share one registry with every other test in this process.
type fakeEngine struct {
	entries []migration.StatusEntry
	// failWith, when set, is what MigrateApp returns instead of applying.
	failWith error
	calls    []string
}

func (f *fakeEngine) SetDb(*gorm.DB) {}

func (f *fakeEngine) Status() ([]migration.StatusEntry, error) {
	out := make([]migration.StatusEntry, len(f.entries))
	copy(out, f.entries)
	return out, nil
}

func (f *fakeEngine) MigrateApp(code string) error {
	f.calls = append(f.calls, code)
	if f.failWith != nil {
		return f.failWith
	}
	for i := range f.entries {
		if f.entries[i].AppCode == code && f.entries[i].Registered {
			f.entries[i].Applied = true
		}
	}
	return nil
}

func orderManifest(version string) app.Manifest {
	return app.Manifest{
		Code:        "order",
		Name:        "Orders",
		Version:     version,
		Description: "order management",
		Author:      "go-admin",
		Requires:    []string{"crm"},
		Pricing:     "free",
		License:     "MIT",
	}
}

func loadRow(t *testing.T, db *gorm.DB, code string) adminmodels.SysApp {
	t.Helper()
	var row adminmodels.SysApp
	if err := db.Where("app_code = ?", code).First(&row).Error; err != nil {
		t.Fatalf("sys_app has no row for %q: %v", code, err)
	}
	return row
}

// A1: a first install records the app, at the version the manifest declares,
// with every descriptive column copied from it.
func TestInstallRecordsAFirstInstall(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true},
		{Version: "order-1786800002000", AppCode: "order", Registered: true},
		{Version: "crm-1786800001000", AppCode: "crm", Registered: true},
	}}

	rep, err := install(db, eng, orderManifest("1.0.0"))
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if rep.NoOp {
		t.Error("a first install reported nothing to do")
	}
	if got, want := len(rep.Applied), 2; got != want {
		t.Errorf("applied %v, want %d versions", rep.Applied, want)
	}
	// Only this app's migrations, not every pending one in the process.
	if len(eng.calls) != 1 || eng.calls[0] != "order" {
		t.Errorf("MigrateApp calls = %v", eng.calls)
	}

	row := loadRow(t, db, "order")
	if row.Status != adminmodels.AppInstalled {
		t.Errorf("status = %d, want installed", row.Status)
	}
	if row.Version != "1.0.0" {
		t.Errorf("version = %q", row.Version)
	}
	if row.InstalledAt == nil {
		t.Error("installed_at was not set")
	}
	if row.Name != "Orders" || row.Author != "go-admin" || row.Description != "order management" {
		t.Errorf("descriptive columns not copied from the manifest: %+v", row)
	}
	if row.Requires != "crm" {
		t.Errorf("requires = %q, want the manifest's list as CSV", row.Requires)
	}
	if row.Pricing != "free" || row.License != "MIT" {
		t.Errorf("the reserved fields were not carried through: %+v", row)
	}
}

// A2: installing the same version again is a no-op, and says so.
func TestInstallIsANoOpAtTheSameVersion(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true},
	}}
	if _, err := install(db, eng, orderManifest("1.0.0")); err != nil {
		t.Fatalf("first install: %v", err)
	}
	before := loadRow(t, db, "order")

	rep, err := install(db, eng, orderManifest("1.0.0"))
	if err != nil {
		t.Fatalf("second install: %v", err)
	}
	if !rep.NoOp {
		t.Error("installing the same version again was not reported as a no-op")
	}
	if len(eng.calls) != 1 {
		t.Errorf("the engine was driven again: %v", eng.calls)
	}
	after := loadRow(t, db, "order")
	if !after.UpdatedAt.Equal(before.UpdatedAt) {
		t.Error("a no-op rewrote the row")
	}
	var n int64
	db.Model(&adminmodels.SysApp{}).Count(&n)
	if n != 1 {
		t.Errorf("sys_app has %d rows, want 1", n)
	}
}

// A no-op is only a no-op when nothing is outstanding. A row that says
// installed while a migration of its has never run is the case sys_app must
// not be believed over sys_migration.
func TestInstallRunsWhenTheRowSaysInstalledButAMigrationIsPending(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true, Applied: true},
	}}
	if _, err := install(db, eng, orderManifest("1.0.0")); err != nil {
		t.Fatalf("first install: %v", err)
	}

	// A second version of the same app appears - the app was rebuilt with
	// one more migration file, without its version changing.
	eng.entries = append(eng.entries, migration.StatusEntry{
		Version: "order-1786800002000", AppCode: "order", Registered: true,
	})

	rep, err := install(db, eng, orderManifest("1.0.0"))
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if rep.NoOp {
		t.Fatal("an outstanding migration was reported as nothing to do")
	}
	if len(rep.Applied) != 1 || rep.Applied[0] != "order-1786800002000" {
		t.Errorf("applied = %v", rep.Applied)
	}
}

// A9: an upgrade is in place. installed_at is the first install's, not this
// one's.
func TestInstallUpgradesInPlaceAndKeepsTheFirstInstallTime(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true},
	}}
	if _, err := install(db, eng, orderManifest("1.0.0")); err != nil {
		t.Fatalf("first install: %v", err)
	}
	first := loadRow(t, db, "order")
	if first.InstalledAt == nil {
		t.Fatal("installed_at was not set by the first install")
	}

	eng.entries = append(eng.entries, migration.StatusEntry{
		Version: "order-1786800002000", AppCode: "order", Registered: true,
	})
	rep, err := install(db, eng, orderManifest("2.0.0"))
	if err != nil {
		t.Fatalf("upgrade: %v", err)
	}
	if rep.Previous != "1.0.0" {
		t.Errorf("previous = %q, want 1.0.0", rep.Previous)
	}

	row := loadRow(t, db, "order")
	if row.Version != "2.0.0" {
		t.Errorf("version = %q, want 2.0.0", row.Version)
	}
	if row.Status != adminmodels.AppInstalled {
		t.Errorf("status = %d, want installed", row.Status)
	}
	if !row.InstalledAt.Equal(*first.InstalledAt) {
		t.Errorf("installed_at moved from %v to %v; an upgrade keeps the first install's time",
			first.InstalledAt, row.InstalledAt)
	}
}

// A10: a downgrade is refused, and refused before anything is written.
func TestInstallRefusesADowngrade(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true},
	}}
	if _, err := install(db, eng, orderManifest("2.0.0")); err != nil {
		t.Fatalf("first install: %v", err)
	}
	before := loadRow(t, db, "order")

	_, err := install(db, eng, orderManifest("1.0.0"))
	if err == nil {
		t.Fatal("a downgrade was accepted")
	}
	if !strings.Contains(err.Error(), "downgrade") {
		t.Errorf("error = %q, it has to say what it refused", err)
	}
	after := loadRow(t, db, "order")
	if after.Version != before.Version || after.Status != before.Status {
		t.Errorf("the refused downgrade still wrote to the row: %+v -> %+v", before, after)
	}
}

// A5: a failing migration leaves a row that says so, and says where.
func TestInstallRecordsAFailure(t *testing.T) {
	db := newInstallDB(t)
	boom := errors.New("the seed hit a duplicate")
	eng := &fakeEngine{
		entries: []migration.StatusEntry{
			{Version: "order-1786800001000", AppCode: "order", Registered: true},
		},
		failWith: &migration.VersionFailure{Version: "order-1786800001000", Err: boom},
	}

	_, err := install(db, eng, orderManifest("1.0.0"))
	if err == nil {
		t.Fatal("a failed install reported success")
	}
	if !errors.Is(err, boom) {
		t.Errorf("the cause is not reachable: %v", err)
	}

	row := loadRow(t, db, "order")
	if row.Status != adminmodels.AppFailed {
		t.Errorf("status = %d, want failed", row.Status)
	}
	if row.FailedVersion != "order-1786800001000" {
		t.Errorf("failed_version = %q", row.FailedVersion)
	}
	if !strings.Contains(row.LastError, "duplicate") {
		t.Errorf("last_error = %q", row.LastError)
	}
	if row.InstalledAt != nil {
		t.Error("installed_at was set by an install that failed")
	}
}

// A failed install is retried by running it again - not by any special
// command, and without the previous attempt's diagnostics surviving into a
// row that now says installed.
func TestInstallResumesAfterAFailure(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{
		entries: []migration.StatusEntry{
			{Version: "order-1786800001000", AppCode: "order", Registered: true},
		},
		failWith: &migration.VersionFailure{Version: "order-1786800001000", Err: errors.New("boom")},
	}
	if _, err := install(db, eng, orderManifest("1.0.0")); err == nil {
		t.Fatal("the first attempt did not fail")
	}

	eng.failWith = nil
	rep, err := install(db, eng, orderManifest("1.0.0"))
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if rep.NoOp {
		t.Error("a failed row was treated as installed")
	}

	row := loadRow(t, db, "order")
	if row.Status != adminmodels.AppInstalled {
		t.Errorf("status = %d, want installed", row.Status)
	}
	if row.FailedVersion != "" || row.LastError != "" {
		t.Errorf("the previous failure survived onto a row that now says installed: %q / %q",
			row.FailedVersion, row.LastError)
	}
	if row.InstalledAt == nil {
		t.Error("installed_at was not set by the attempt that succeeded")
	}
}

// A row stuck at installing - the process was killed partway - is not
// installed, and must not be mistaken for it.
func TestInstallRetriesARowStuckAtInstalling(t *testing.T) {
	db := newInstallDB(t)
	if err := db.Create(&adminmodels.SysApp{
		AppCode: "order", Name: "Orders", Version: "1.0.0",
		Status: adminmodels.AppInstalling,
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true, Applied: true},
	}}

	rep, err := install(db, eng, orderManifest("1.0.0"))
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if rep.NoOp {
		t.Fatal("a row stuck at installing was reported as already installed")
	}
	if row := loadRow(t, db, "order"); row.Status != adminmodels.AppInstalled {
		t.Errorf("status = %d, want installed", row.Status)
	}
}

func TestInstallRejectsTheFrameworkCode(t *testing.T) {
	db := newInstallDB(t)
	m := orderManifest("1.0.0")
	m.Code = migration.FrameworkAppCode
	_, err := install(db, &fakeEngine{}, m)
	if err == nil {
		t.Fatal("the framework was installed as an application")
	}
	if !strings.Contains(err.Error(), "migrate") {
		t.Errorf("error = %q, it should point at the command that does this", err)
	}
}

func TestInstallRefusesAnUnparseableRecordedVersion(t *testing.T) {
	db := newInstallDB(t)
	if err := db.Create(&adminmodels.SysApp{
		AppCode: "order", Name: "Orders", Version: "v1.0", Status: adminmodels.AppInstalled,
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	_, err := install(db, &fakeEngine{}, orderManifest("1.0.0"))
	if err == nil {
		t.Fatal("an unparseable recorded version was compared anyway")
	}
	row := loadRow(t, db, "order")
	if row.Status != adminmodels.AppInstalled || row.Version != "v1.0" {
		t.Errorf("the row was overwritten before the comparison failed: %+v", row)
	}
}

// A7: the report has to say the code is not running yet. Menus appearing is
// exactly what makes an operator think it is.
func TestReportInstallSaysTheCodeIsNotRunningYet(t *testing.T) {
	var out strings.Builder
	reportInstall(&out, installReport{Code: "order", Version: "1.0.0", Applied: []string{"order-1786800001000"}})
	got := out.String()
	if !strings.Contains(got, "rebuild") || !strings.Contains(got, "restart") {
		t.Errorf("the report does not say the binary has to be rebuilt: %q", got)
	}
	if !strings.Contains(got, "order-1786800001000") {
		t.Errorf("the report does not name what it applied: %q", got)
	}
}

func TestReportInstallOnANoOp(t *testing.T) {
	var out strings.Builder
	reportInstall(&out, installReport{Code: "order", Version: "1.0.0", NoOp: true})
	if !strings.Contains(out.String(), "already installed") {
		t.Errorf("output = %q", out.String())
	}
}

// last_error is a varchar(255) declared in characters. A message that is
// partly Chinese would be cut mid-rune by a byte-wise truncation and stored
// as an invalid sequence.
func TestTruncateCutsRunesNotBytes(t *testing.T) {
	s := strings.Repeat("迁", 300)
	got := truncate(s, 255)
	if n := len([]rune(got)); n != 255 {
		t.Errorf("kept %d runes, want 255", n)
	}
	if !strings.HasPrefix(s, got) {
		t.Error("truncation did not cut at a rune boundary")
	}
	if short := truncate("ok", 255); short != "ok" {
		t.Errorf("a short message was altered: %q", short)
	}
}

func TestInstallNeedsSysApp(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	_, err = install(db, &fakeEngine{}, orderManifest("1.0.0"))
	if err == nil {
		t.Fatal("install ran against a database with no sys_app")
	}
	if !strings.Contains(err.Error(), "migrate") {
		t.Errorf("error = %q, it should say what to run first", err)
	}
}

// The code written to sys_app and handed to the engine is the normalized one.
// A manifest whose Code was typed with different case or stray spaces has to
// land on the same identity migration.ForApp and seed.SeedMenus already use,
// or the row and the migrations it stands for are filed under two names.
func TestInstallNormalizesTheAppCode(t *testing.T) {
	db := newInstallDB(t)
	eng := &fakeEngine{entries: []migration.StatusEntry{
		{Version: "order-1786800001000", AppCode: "order", Registered: true},
	}}
	m := orderManifest("1.0.0")
	m.Code = "  Order  "

	rep, err := install(db, eng, m)
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if rep.Code != "order" {
		t.Errorf("reported code = %q, want order", rep.Code)
	}
	if len(eng.calls) != 1 || eng.calls[0] != "order" {
		t.Errorf("the engine was asked for %v, want [order]", eng.calls)
	}
	// The row has to be findable by the normalized code, which is what every
	// other table in this batch is keyed by.
	row := loadRow(t, db, "order")
	if row.AppCode != "order" {
		t.Errorf("app_code = %q", row.AppCode)
	}
	if len(rep.Applied) != 1 {
		t.Errorf("applied = %v; the normalized code has to match what Status reports", rep.Applied)
	}
}
