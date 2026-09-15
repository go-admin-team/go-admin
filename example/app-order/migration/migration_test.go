package migration

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	contractmigration "github.com/go-admin-team/go-admin-core/v2/sdk/contract/migration"
	contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/seed"

	"github.com/go-admin-team/example-app-order/models"
)

// fakeSeeder stands in for a host's real Seeder (the one wt-shim, as of
// this writing, never registers - see the accompanying report's gap list).
// It records what it received instead of writing to any table, which is
// enough to check app-order's own MenuSpec/ApiSpec assembly without
// depending on go-admin's sys_menu/sys_api schema.
type fakeSeeder struct {
	appCode string
	menus   []seed.MenuSpec
	apis    []seed.ApiSpec
}

func (f *fakeSeeder) SeedMenus(tx *gorm.DB, appCode string, menus []seed.MenuSpec, apis []seed.ApiSpec) error {
	f.appCode = appCode
	f.menus = menus
	f.apis = apis
	return nil
}

// seed.RegisterSeeder panics on a second call in the same process (see its
// doc comment) - by design, there is no public way to unregister one. This
// package's tests share the one registration below rather than each
// registering their own.
var fake = &fakeSeeder{}

func init() {
	seed.RegisterSeeder(fake)
}

// TestRegistersUnderContractMigrationForApp is this package's core claim:
// that createOrderSchema is reachable through contract/migration's
// package-level Snapshot, the only registry a third-party module can
// register against. It does not confirm any host actually calls Snapshot
// today - see the report.
func TestRegistersUnderContractMigrationForApp(t *testing.T) {
	entries := contractmigration.Snapshot()
	entry, ok := entries[AppCode+"-"+version]
	if !ok {
		t.Fatalf("no entry for %s-%s; registered: %v", AppCode, version, keysOf(entries))
	}
	if entry.AppCode != AppCode {
		t.Errorf("Entry.AppCode = %q, want %q", entry.AppCode, AppCode)
	}
}

func keysOf(m map[string]contractmigration.Entry) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// TestMigrationCreatesTablesSeedsMenusAndRecordsItself runs the registered
// migration function directly against a fresh sqlite database - standing in
// for the host's execution engine, which (see the report) does not exist
// yet for an externally-registered app. It is the closest thing to an
// end-to-end run this example can do without wt-shim's cooperation.
func TestMigrationCreatesTablesSeedsMenusAndRecordsItself(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	// sys_migration itself is created by the framework's own first
	// migration (go-admin's cmd/migrate/migration/version/*_tables.go),
	// which by the time any app's migration runs has always already run -
	// simulate that precondition rather than app-order's own migration
	// creating a table it does not own.
	if err := db.AutoMigrate(&contractmodels.Migration{}); err != nil {
		t.Fatalf("automigrate sys_migration: %v", err)
	}

	entries := contractmigration.Snapshot()
	entry, ok := entries[AppCode+"-"+version]
	if !ok {
		t.Fatalf("no entry for %s-%s", AppCode, version)
	}
	if err := entry.Fn(db, AppCode+"-"+version); err != nil {
		t.Fatalf("running the registered migration: %v", err)
	}

	if !db.Migrator().HasTable(&models.Order{}) {
		t.Error("app_order was not created")
	}
	if !db.Migrator().HasTable(&models.OrderItem{}) {
		t.Error("app_order_item was not created")
	}

	var migrationRow contractmodels.Migration
	if err := db.Where("version = ?", AppCode+"-"+version).First(&migrationRow).Error; err != nil {
		t.Fatalf("sys_migration row: %v", err)
	}
	if migrationRow.AppCode != AppCode {
		t.Errorf("sys_migration.app_code = %q, want %q", migrationRow.AppCode, AppCode)
	}

	if fake.appCode != AppCode {
		t.Errorf("Seeder saw appCode %q, want %q", fake.appCode, AppCode)
	}
	assertMenuGraphIsConsistent(t, fake.menus, fake.apis)
}

// assertMenuGraphIsConsistent checks the two rules that would otherwise
// only surface as a broken admin UI at install time: every Parent
// reference resolves to a Code in the same batch, and the frontend's
// apps/<code>/ convention for a packaged page's Component (documented on
// MenuSpec.Component, enforced by nothing - see the report) is actually
// followed.
func assertMenuGraphIsConsistent(t *testing.T, menus []seed.MenuSpec, apis []seed.ApiSpec) {
	t.Helper()

	codes := make(map[string]seed.MenuSpec, len(menus))
	for _, m := range menus {
		codes[m.Code] = m
	}
	apiCodes := make(map[string]bool, len(apis))
	for _, a := range apis {
		apiCodes[a.Code] = true
	}

	for _, m := range menus {
		if m.Parent != "" {
			if _, ok := codes[m.Parent]; !ok {
				t.Errorf("menu %q has Parent %q, which is not a Code in this batch", m.Code, m.Parent)
			}
		}
		for _, ac := range m.ApiCodes {
			if !apiCodes[ac] {
				t.Errorf("menu %q references ApiCode %q, which is not in this batch's apis", m.Code, ac)
			}
		}
		if m.Kind == contractmodels.Menu && m.Component != "" {
			if !strings.HasPrefix(m.Component, "apps/"+AppCode+"/") {
				t.Errorf("menu %q has Component %q, want it to start with apps/%s/", m.Code, m.Component, AppCode)
			}
		}
	}
}
