package migrate

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/seed"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	adminmodels "go-admin/app/admin/models"
	_ "go-admin/app/admin/service" // registers the seeder SeedMenus dispatches to
	"go-admin/cmd/migrate/migration"
	commonmodels "go-admin/common/models"
)

const adminRoleKey = "admin"

// newUninstallDB builds every table an install writes to, plus one table
// standing in for the application's own data, which an uninstall must not
// touch.
func newUninstallDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&adminmodels.SysMenu{}, &adminmodels.SysApi{}, &adminmodels.SysRole{},
		&adminmodels.SysApp{}, &adminmodels.SysAppCasbinGrant{}, &commonmodels.Migration{},
	); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	// casbin_rule has no GORM model in this repository; the columns are the
	// ones grantToAdminRole's INSERT addresses.
	if err := db.Exec(`CREATE TABLE casbin_rule (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ptype TEXT, v0 TEXT, v1 TEXT, v2 TEXT, v3 TEXT, v4 TEXT, v5 TEXT
	)`).Error; err != nil {
		t.Fatalf("create casbin_rule: %v", err)
	}
	if err := db.Exec(`CREATE TABLE app_order (id INTEGER PRIMARY KEY, note TEXT)`).Error; err != nil {
		t.Fatalf("create app_order: %v", err)
	}
	if err := db.Exec(`INSERT INTO app_order (id, note) VALUES (1, 'a real order')`).Error; err != nil {
		t.Fatalf("seed app_order: %v", err)
	}
	if err := db.Create(&adminmodels.SysRole{RoleName: "Administrator", RoleKey: adminRoleKey}).Error; err != nil {
		t.Fatalf("seed admin role: %v", err)
	}
	return db
}

// specsFor builds one application's menus and apis. The paths carry the app
// code, because two applications do not share an endpoint - and if a fixture
// let them, the second one's policies would already exist and its ledger
// would legitimately come out empty, which would make it a useless control.
func specsFor(code string) ([]seed.MenuSpec, []seed.ApiSpec) {
	menus := []seed.MenuSpec{
		{Code: "dir", Kind: "M", Title: code + " example", Path: "/apps/" + code, Component: "Layout", Sort: 10},
		{Code: "list", Parent: "dir", Kind: "C", Title: code, Path: "list", Component: "apps/" + code + "/index", Sort: 1, ApiCodes: []string{"list"}},
	}
	apis := []seed.ApiSpec{
		{Code: "list", Title: code + " list", Path: "/api/v1/" + code, Method: "GET", Handle: "apis." + code + ".GetPage-fm"},
		{Code: "create", Title: "create " + code, Path: "/api/v1/" + code, Method: "POST", Handle: "apis." + code + ".Insert-fm"},
	}
	return menus, apis
}

// seedApp runs the real seeding path, so what the uninstaller has to undo is
// what an install actually writes rather than a hand-built approximation.
func seedApp(t *testing.T, db *gorm.DB, code string) {
	t.Helper()
	menus, apis := specsFor(code)
	if err := db.Transaction(func(tx *gorm.DB) error {
		return seed.SeedMenus(tx, code, menus, apis)
	}); err != nil {
		t.Fatalf("seeding %q: %v", code, err)
	}
	if err := db.Create(&commonmodels.Migration{
		Version: code + "-1786800001000", AppCode: code, ApplyTime: time.Now(),
	}).Error; err != nil {
		t.Fatalf("recording the migration for %q: %v", code, err)
	}
	if err := db.Create(&adminmodels.SysApp{
		AppCode: code, Name: code, Version: "1.0.0", Status: adminmodels.AppInstalled,
	}).Error; err != nil {
		t.Fatalf("recording sys_app for %q: %v", code, err)
	}
}

func count(t *testing.T, db *gorm.DB, table, where string, args ...any) int64 {
	t.Helper()
	var n int64
	q := db.Table(table)
	if where != "" {
		q = q.Where(where, args...)
	}
	if err := q.Count(&n).Error; err != nil {
		t.Fatalf("counting %s: %v", table, err)
	}
	return n
}

// A3: everything the install wrote goes, and the application's own table does
// not.
func TestUninstallRemovesWhatWasSeededAndNothingElse(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")

	if count(t, db, "sys_menu", "app_code = ?", "order") == 0 {
		t.Fatal("nothing was seeded, so this test proves nothing")
	}

	rep, err := uninstall(db, "order")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if !rep.Found {
		t.Error("the sys_app row was not found")
	}

	for _, c := range []struct {
		table, where string
		args         []any
	}{
		{"sys_menu", "app_code = ?", []any{"order"}},
		{"sys_api", "app_code = ?", []any{"order"}},
		{"sys_menu_api_rule", "", nil},
		{"sys_role_menu", "", nil},
		{"casbin_rule", "v1 = ?", []any{"/api/v1/order"}},
		{"sys_app_casbin_grant", "app_code = ?", []any{"order"}},
		{"sys_migration", "app_code = ?", []any{"order"}},
		{"sys_app", "app_code = ?", []any{"order"}},
	} {
		if n := count(t, db, c.table, c.where, c.args...); n != 0 {
			t.Errorf("%s still has %d row(s)", c.table, n)
		}
	}
	if n := count(t, db, "app_order", "", nil); n != 1 {
		t.Errorf("app_order has %d row(s); the application's own data is not the uninstaller's to remove", n)
	}
	if len(rep.Skipped) != 0 || len(rep.Orphans) != 0 {
		t.Errorf("a clean uninstall reported skipped=%v orphans=%v", rep.Skipped, rep.Orphans)
	}
	if rep.Menus == 0 || rep.Apis == 0 || rep.Policies == 0 || rep.Migrations == 0 {
		t.Errorf("the report says nothing was removed: %+v", rep)
	}
}

// Uninstalling one application must not reach into another's rows. Every
// delete here is filtered, and a missing filter is invisible on a database
// with only one application in it.
func TestUninstallLeavesAnotherApplicationAlone(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")
	seedApp(t, db, "crm")

	before := map[string]int64{
		"sys_menu":             count(t, db, "sys_menu", "app_code = ?", "crm"),
		"sys_api":              count(t, db, "sys_api", "app_code = ?", "crm"),
		"sys_app_casbin_grant": count(t, db, "sys_app_casbin_grant", "app_code = ?", "crm"),
		"sys_migration":        count(t, db, "sys_migration", "app_code = ?", "crm"),
		"sys_app":              count(t, db, "sys_app", "app_code = ?", "crm"),
	}
	for k, v := range before {
		if v == 0 {
			t.Fatalf("crm has no rows in %s, so this test proves nothing", k)
		}
	}
	crmBindings := count(t, db, "sys_menu_api_rule", "", nil)
	crmRoleMenus := count(t, db, "sys_role_menu", "", nil)

	if _, err := uninstall(db, "order"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}

	for k, v := range before {
		if n := count(t, db, k, "app_code = ?", "crm"); n != v {
			t.Errorf("%s for crm went from %d to %d", k, v, n)
		}
	}
	// crm's own bindings and role rows are half of each total, and must be
	// exactly what is left.
	if n := count(t, db, "sys_menu_api_rule", "", nil); n != crmBindings/2 {
		t.Errorf("sys_menu_api_rule = %d, want %d (crm's half)", n, crmBindings/2)
	}
	if n := count(t, db, "sys_role_menu", "", nil); n != crmRoleMenus/2 {
		t.Errorf("sys_role_menu = %d, want %d (crm's half)", n, crmRoleMenus/2)
	}
	// crm's policies name a different path, so they are untouched.
	if n := count(t, db, "casbin_rule", "v1 = ?", "/api/v1/order"); n != 0 {
		t.Errorf("order's policies survived: %d", n)
	}
}

// A6b: somebody granted this application's API to another role by hand. That
// grant is not in the ledger, is not this uninstall's to remove, and would
// otherwise vanish from view entirely.
func TestUninstallReportsAGrantSomebodyElseMade(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")

	if err := db.Exec(
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ('p', 'ops', '/api/v1/order', 'GET', '', '', '')",
	).Error; err != nil {
		t.Fatalf("hand-made grant: %v", err)
	}

	rep, err := uninstall(db, "order")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if n := count(t, db, "casbin_rule", "v0 = ?", "ops"); n != 1 {
		t.Errorf("somebody else's grant was deleted (%d rows left)", n)
	}
	if len(rep.Orphans) != 1 {
		t.Fatalf("orphans = %v, want the one hand-made grant", rep.Orphans)
	}
	if rep.Orphans[0].V0 != "ops" {
		t.Errorf("orphan = %+v", rep.Orphans[0])
	}
	// The admin grants it did own are gone.
	if n := count(t, db, "casbin_rule", "v0 = ?", adminRoleKey); n != 0 {
		t.Errorf("%d of this app's own policies survived", n)
	}

	var out strings.Builder
	reportUninstall(&out, rep)
	if !strings.Contains(out.String(), "ops") || !strings.Contains(out.String(), "left alone") {
		t.Errorf("the report does not say what was left behind: %q", out.String())
	}
}

// A6a: a policy this install created is not there any more. Not an error -
// the uninstall wanted it gone and it is - but reported, because something
// else rewrote casbin_rule.
func TestUninstallReportsALedgerEntryWhosePolicyIsGone(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")

	if err := db.Exec("DELETE FROM casbin_rule WHERE v2 = 'POST'").Error; err != nil {
		t.Fatalf("removing a policy: %v", err)
	}

	rep, err := uninstall(db, "order")
	if err != nil {
		t.Fatalf("a missing policy made the uninstall fail: %v", err)
	}
	if len(rep.Skipped) != 1 {
		t.Fatalf("skipped = %v, want the one that had gone", rep.Skipped)
	}
	if rep.Skipped[0].V2 != "POST" {
		t.Errorf("skipped = %+v", rep.Skipped[0])
	}
	if n := count(t, db, "sys_app_casbin_grant", "", nil); n != 0 {
		t.Errorf("the ledger kept %d row(s); its job ends with the uninstall", n)
	}
	// It still committed: a skip is a reported branch, not a failure.
	if n := count(t, db, "sys_menu", "app_code = ?", "order"); n != 0 {
		t.Errorf("the transaction rolled back over a skip: sys_menu has %d row(s)", n)
	}
}

// G5/A4: without this the reinstall finds every version applied, runs no
// migration, seeds nothing, and reports success.
func TestUninstallClearsThisAppsMigrationRecordsOnly(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")
	if err := db.Create(&commonmodels.Migration{
		Version: "1786700001000", AppCode: "", ApplyTime: time.Now(),
	}).Error; err != nil {
		t.Fatalf("framework migration row: %v", err)
	}

	if _, err := uninstall(db, "order"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if n := count(t, db, "sys_migration", "app_code = ?", "order"); n != 0 {
		t.Errorf("sys_migration still has %d row(s) for order; a reinstall would seed nothing", n)
	}
	if n := count(t, db, "sys_migration", "app_code = ?", ""); n != 1 {
		t.Errorf("the framework's own migration record was removed (%d left)", n)
	}
}

// A11: sys_role_menu is found by menu id, not by a column on it. A column
// would have been blanked the first time somebody edited a role, because
// SysRole.Update deletes the role's rows and writes them back through GORM's
// many2many, which does not carry extra columns. This reproduces that edit.
func TestUninstallSurvivesARoleMenuRewrite(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")

	var roleID int
	if err := db.Model(&adminmodels.SysRole{}).Where("role_key = ?", adminRoleKey).
		Pluck("role_id", &roleID).Error; err != nil {
		t.Fatalf("reading the admin role: %v", err)
	}
	var menuIDs []int
	if err := db.Model(&adminmodels.SysMenu{}).Where("app_code = ?", "order").
		Pluck("menu_id", &menuIDs).Error; err != nil {
		t.Fatalf("reading menus: %v", err)
	}
	if len(menuIDs) == 0 {
		t.Fatal("no menus were seeded")
	}
	// What SysRole.Update does: drop every row for the role, then write them
	// back with nothing but the two keys.
	if err := db.Exec("DELETE FROM sys_role_menu WHERE role_id = ?", roleID).Error; err != nil {
		t.Fatalf("clearing role menus: %v", err)
	}
	for _, id := range menuIDs {
		if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", roleID, id).Error; err != nil {
			t.Fatalf("rewriting role menus: %v", err)
		}
	}

	rep, err := uninstall(db, "order")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if rep.RoleMenus != int64(len(menuIDs)) {
		t.Errorf("removed %d role assignment(s), want %d", rep.RoleMenus, len(menuIDs))
	}
	if n := count(t, db, "sys_role_menu", "", nil); n != 0 {
		t.Errorf("sys_role_menu still has %d row(s) after a role edit", n)
	}
}

// `migrate` with no subcommand applies every registered migration, an
// application's included, so an application can have all of its rows and
// never have had a sys_app row. Refusing to clean that up would leave the
// only case where nothing else can.
func TestUninstallWorksWithoutASysAppRow(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")
	if err := db.Where("app_code = ?", "order").Delete(&adminmodels.SysApp{}).Error; err != nil {
		t.Fatalf("removing the sys_app row: %v", err)
	}

	rep, err := uninstall(db, "order")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if rep.Found {
		t.Error("the report claims a sys_app row that was not there")
	}
	if n := count(t, db, "sys_menu", "app_code = ?", "order"); n != 0 {
		t.Errorf("sys_menu still has %d row(s)", n)
	}
	var out strings.Builder
	reportUninstall(&out, rep)
	if !strings.Contains(out.String(), "no sys_app row") {
		t.Errorf("the report does not say the row was missing: %q", out.String())
	}
}

// One transaction, and it really is one: nothing here runs DDL, so unlike an
// install there is nothing to commit it out from under itself.
func TestUninstallRollsBackAsAWhole(t *testing.T) {
	db := newUninstallDB(t)
	seedApp(t, db, "order")
	menusBefore := count(t, db, "sys_menu", "app_code = ?", "order")
	policiesBefore := count(t, db, "casbin_rule", "", nil)

	// Step 8's table is gone, so the uninstall fails after it has already
	// deleted menus, apis, bindings and policies.
	if err := db.Migrator().DropTable(&commonmodels.Migration{}); err != nil {
		t.Fatalf("dropping sys_migration: %v", err)
	}

	if _, err := uninstall(db, "order"); err == nil {
		t.Fatal("the uninstall reported success with sys_migration missing")
	}
	if n := count(t, db, "sys_menu", "app_code = ?", "order"); n != menusBefore {
		t.Errorf("sys_menu = %d, want %d: the failed uninstall did not roll back", n, menusBefore)
	}
	if n := count(t, db, "casbin_rule", "", nil); n != policiesBefore {
		t.Errorf("casbin_rule = %d, want %d: the failed uninstall did not roll back", n, policiesBefore)
	}
}

func TestUninstallRefusesTheFrameworkCode(t *testing.T) {
	db := newUninstallDB(t)
	if _, err := uninstall(db, migration.FrameworkAppCode); err == nil {
		t.Fatal("the framework was uninstalled")
	}
}

// findOrphanPolicies chunks its OR chain because a driver runs out of
// placeholders long before an application runs out of endpoints. The
// boundary is where an off-by-one hides: a chunk size that drops the last
// element of each batch, or one that never advances, both leave policies
// unreported and nothing says so.
func TestFindOrphanPoliciesCoversEveryPathAcrossChunks(t *testing.T) {
	db := newUninstallDB(t)

	// Deliberately not a multiple of the chunk size, so the last batch is
	// short, and large enough to need three of them.
	const n = 205
	keys := make([]policyKey, 0, n)
	for i := 0; i < n; i++ {
		path := "/api/v1/thing" + strconv.Itoa(i)
		keys = append(keys, policyKey{V1: path, V2: "GET"})
		if err := db.Exec(
			"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ('p', 'ops', ?, 'GET', '', '', '')",
			path,
		).Error; err != nil {
			t.Fatalf("seeding policy %d: %v", i, err)
		}
	}
	// One policy that must not match: a path no key names.
	if err := db.Exec(
		"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES ('p', 'ops', '/api/v1/elsewhere', 'GET', '', '', '')",
	).Error; err != nil {
		t.Fatalf("seeding the control policy: %v", err)
	}

	found, err := findOrphanPolicies(db, keys)
	if err != nil {
		t.Fatalf("findOrphanPolicies: %v", err)
	}
	if len(found) != n {
		t.Fatalf("found %d policies, want %d", len(found), n)
	}
	seen := make(map[string]bool, len(found))
	for _, f := range found {
		seen[f.V1] = true
		if f.V1 == "/api/v1/elsewhere" {
			t.Error("a path no key names was reported")
		}
	}
	for _, k := range keys {
		if !seen[k.V1] {
			t.Errorf("%s was not reported", k.V1)
		}
	}
}

// An application may register apis with no menus at all - endpoints another
// service calls - so either of the id lists an uninstall reads can be empty.
// The guard in front of the join-table delete turns out not to be what makes
// this work: GORM renders IN with an empty slice as a condition that matches
// nothing, rather than the empty IN list that would be a syntax error in raw
// SQL, and removing the guard leaves this test green. It stays as an explicit
// statement of intent rather than a reliance on that rendering.
func TestUninstallWithApisButNoMenus(t *testing.T) {
	db := newUninstallDB(t)
	apis := []seed.ApiSpec{
		{Code: "hook", Title: "Inbound hook", Path: "/api/v1/hook", Method: "POST", Handle: "hook.Receive"},
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return seed.SeedMenus(tx, "hooks", nil, apis)
	}); err != nil {
		t.Fatalf("seeding: %v", err)
	}

	rep, err := uninstall(db, "hooks")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if rep.Apis != 1 {
		t.Errorf("removed %d api(s), want 1", rep.Apis)
	}
	if rep.Policies != 1 {
		t.Errorf("removed %d policy(ies), want 1", rep.Policies)
	}
	if n := count(t, db, "casbin_rule", "", nil); n != 0 {
		t.Errorf("casbin_rule has %d row(s)", n)
	}
}

// The mirror case: menus and no apis at all.
func TestUninstallWithMenusButNoApis(t *testing.T) {
	db := newUninstallDB(t)
	menus := []seed.MenuSpec{
		{Code: "dir", Kind: "M", Title: "Reports", Path: "/apps/reports", Component: "Layout", Sort: 10},
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		return seed.SeedMenus(tx, "reports", menus, nil)
	}); err != nil {
		t.Fatalf("seeding: %v", err)
	}

	rep, err := uninstall(db, "reports")
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if rep.Menus != 1 {
		t.Errorf("removed %d menu(s), want 1", rep.Menus)
	}
	if n := count(t, db, "sys_menu", "app_code = ?", "reports"); n != 0 {
		t.Errorf("sys_menu has %d row(s)", n)
	}
	if len(rep.Orphans) != 0 {
		t.Errorf("an application with no apis reported orphans: %v", rep.Orphans)
	}
}
