package service

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/seed"

	"go-admin/app/admin/models"
)

// newSeedTestDB builds the tables adminSeeder.SeedMenus writes to. sys_menu,
// sys_api, sys_role and sys_role_menu (GORM's own join table for
// SysRole.SysMenu) come from AutoMigrate; casbin_rule does not have a GORM
// model anywhere in this codebase - see 1786700001000_demo_menu.go's own
// comment on why models.CasbinRule (-> sys_casbin_rule) is the wrong table -
// so it is created directly, matching the columns grantToAdminRole's INSERT
// addresses.
func newSeedTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.AutoMigrate(&models.SysMenu{}, &models.SysApi{}, &models.SysRole{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	if err := db.Exec(`CREATE TABLE casbin_rule (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ptype TEXT, v0 TEXT, v1 TEXT, v2 TEXT, v3 TEXT, v4 TEXT, v5 TEXT
	)`).Error; err != nil {
		t.Fatalf("create casbin_rule: %v", err)
	}
	return db
}

func seedAdminRole(t *testing.T, db *gorm.DB) models.SysRole {
	t.Helper()
	role := models.SysRole{RoleName: "Administrator", RoleKey: adminRoleKey}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("seed admin role: %v", err)
	}
	return role
}

// This is the acceptance case go-admin-core's docs/contract.md requires: one
// SeedMenus call populates all four tables a visible, working menu entry
// needs, every row tagged with the appCode it was called with, and the
// parent/child tree resolved into sys_menu's parent_id/paths.
func TestSeedMenusPopulatesAllFourTables(t *testing.T) {
	db := newSeedTestDB(t)
	seedAdminRole(t, db)

	menus := []seed.MenuSpec{
		{Code: "dir", Kind: contractmodels.Directory, Title: "Order Example", Path: "/apps/order", Component: "Layout", Sort: 10},
		{Code: "list", Parent: "dir", Kind: contractmodels.Menu, Title: "Orders", Path: "list", Component: "apps/order/order/index", Sort: 1, ApiCodes: []string{"list"}},
		{Code: "btn-create", Parent: "list", Kind: contractmodels.Button, Title: "Create", Permission: "order:order:create", Sort: 1},
	}
	apis := []seed.ApiSpec{
		{Code: "list", Title: "Order list", Path: "/api/v1/order", Method: "GET", Handle: "apis.Order.GetPage-fm"},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		return adminSeeder{}.SeedMenus(tx, "order", menus, apis)
	})
	if err != nil {
		t.Fatalf("SeedMenus: %v", err)
	}

	var apiRows []models.SysApi
	if err := db.Find(&apiRows).Error; err != nil {
		t.Fatal(err)
	}
	if len(apiRows) != 1 || apiRows[0].AppCode != "order" || apiRows[0].Path != "/api/v1/order" {
		t.Fatalf("sys_api = %+v", apiRows)
	}

	var menuRows []models.SysMenu
	if err := db.Order("sort").Find(&menuRows).Error; err != nil {
		t.Fatal(err)
	}
	if len(menuRows) != 3 {
		t.Fatalf("sys_menu has %d rows, want 3: %+v", len(menuRows), menuRows)
	}
	byName := map[string]models.SysMenu{}
	for _, m := range menuRows {
		if m.AppCode != "order" {
			t.Errorf("menu %q app_code = %q, want order", m.MenuName, m.AppCode)
		}
		byName[m.MenuName] = m
	}
	dir, ok := byName[menuName("order", "dir")]
	if !ok || dir.ParentId != 0 || dir.Paths != "/0/"+strconv.Itoa(dir.MenuId) {
		t.Fatalf("dir menu = %+v", dir)
	}
	list, ok := byName[menuName("order", "list")]
	if !ok || list.ParentId != dir.MenuId || list.Paths != dir.Paths+"/"+strconv.Itoa(list.MenuId) {
		t.Fatalf("list menu = %+v (dir=%+v)", list, dir)
	}
	btn, ok := byName[menuName("order", "btn-create")]
	if !ok || btn.ParentId != list.MenuId {
		t.Fatalf("btn menu = %+v (list=%+v)", btn, list)
	}

	// sys_menu_api_rule: gorm's own many2many join table for SysMenu.SysApi.
	var joinCount int64
	if err := db.Table("sys_menu_api_rule").
		Where("sys_menu_menu_id = ? AND sys_api_id = ?", list.MenuId, apiRows[0].Id).
		Count(&joinCount).Error; err != nil {
		t.Fatal(err)
	}
	if joinCount != 1 {
		t.Errorf("sys_menu_api_rule has %d row(s) linking list to its api, want 1", joinCount)
	}

	// sys_role_menu: every seeded menu granted to the admin role.
	var roleMenuCount int64
	if err := db.Table("sys_role_menu").Count(&roleMenuCount).Error; err != nil {
		t.Fatal(err)
	}
	if roleMenuCount != 3 {
		t.Errorf("sys_role_menu has %d row(s), want 3 (one per seeded menu)", roleMenuCount)
	}

	// casbin_rule: the api's path/method granted to the admin role.
	var casbinCount int64
	if err := db.Table("casbin_rule").
		Where("ptype = 'p' AND v0 = ? AND v1 = ? AND v2 = ?", adminRoleKey, "/api/v1/order", "GET").
		Count(&casbinCount).Error; err != nil {
		t.Fatal(err)
	}
	if casbinCount != 1 {
		t.Errorf("casbin_rule has %d matching row(s), want 1", casbinCount)
	}
}

// A database that has not run the framework's own seed data yet (no admin
// role) must not fail SeedMenus - 1786700001000_demo_menu.go tolerates
// exactly the same condition for the host's own demo module.
func TestSeedMenusToleratesMissingAdminRole(t *testing.T) {
	db := newSeedTestDB(t)

	err := db.Transaction(func(tx *gorm.DB) error {
		return adminSeeder{}.SeedMenus(tx, "order", []seed.MenuSpec{
			{Code: "dir", Kind: contractmodels.Directory, Title: "Order"},
		}, nil)
	})
	if err != nil {
		t.Fatalf("SeedMenus: %v", err)
	}

	var roleMenuCount int64
	if err := db.Table("sys_role_menu").Count(&roleMenuCount).Error; err != nil {
		t.Fatal(err)
	}
	if roleMenuCount != 0 {
		t.Errorf("sys_role_menu has %d row(s) with no role to grant to", roleMenuCount)
	}
}

func TestSeedMenusRejectsMalformedSpecs(t *testing.T) {
	cases := []struct {
		name  string
		menus []seed.MenuSpec
		apis  []seed.ApiSpec
		want  string
	}{
		{
			name:  "duplicate menu code",
			menus: []seed.MenuSpec{{Code: "a", Kind: contractmodels.Directory}, {Code: "a", Kind: contractmodels.Directory}},
			want:  `duplicate MenuSpec.Code "a"`,
		},
		{
			name:  "unresolved parent",
			menus: []seed.MenuSpec{{Code: "a", Parent: "missing", Kind: contractmodels.Menu}},
			want:  `Parent "missing" is not a Code in this call`,
		},
		{
			name:  "unresolved api code",
			menus: []seed.MenuSpec{{Code: "a", Kind: contractmodels.Menu, ApiCodes: []string{"missing"}}},
			want:  `ApiCodes references "missing"`,
		},
		{
			name:  "unknown kind",
			menus: []seed.MenuSpec{{Code: "a", Kind: "X"}},
			want:  `Kind "X" is not one of Directory/Menu/Button`,
		},
		{
			name:  "sort overflows a tinyint",
			menus: []seed.MenuSpec{{Code: "a", Kind: contractmodels.Directory, Sort: 900}},
			want:  `Sort 900 does not fit sys_menu.sort's tinyint column`,
		},
		{
			name: "duplicate api code",
			apis: []seed.ApiSpec{{Code: "x"}, {Code: "x"}},
			want: `duplicate ApiSpec.Code "x"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := newSeedTestDB(t)
			err := db.Transaction(func(tx *gorm.DB) error {
				return adminSeeder{}.SeedMenus(tx, "order", tc.menus, tc.apis)
			})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tc.want)
			}
		})
	}
}

// TestSeederIsRegistered pins the registration itself, not the behaviour.
//
// Every other test here calls adminSeeder{}.SeedMenus directly, which proves
// the implementation is right and proves nothing about whether anything ever
// reaches it: delete the RegisterSeeder call in init() and they all stay
// green, while a real migrate fails with ErrNoSeeder and no menu is written.
// Going through the package-level SeedMenus is what closes that gap - it is
// the door an application actually knocks on.
func TestSeederIsRegistered(t *testing.T) {
	db := newSeedTestDB(t)
	err := db.Transaction(func(tx *gorm.DB) error {
		return seed.SeedMenus(tx, "probe", []seed.MenuSpec{{
			Code: "root", Kind: contractmodels.Directory, Title: "Probe", Sort: 1,
		}}, nil)
	})
	if errors.Is(err, seed.ErrNoSeeder) {
		t.Fatal("no Seeder is registered: an application's SeedMenus would write no menu at all")
	}
	if err != nil {
		t.Fatalf("SeedMenus through the package-level entry point: %v", err)
	}
}

// An application is free to register apis with no menus at all - endpoints
// another service calls, or a UI mounted somewhere else. Skipping
// grantToAdminRole on an empty menu list wrote the sys_api rows and then no
// casbin rule for them, so every one of those endpoints was denied to
// everyone including admin, from a migration that reported success.
func TestSeedMenusGrantsApisWhenThereAreNoMenus(t *testing.T) {
	db := newSeedTestDB(t)
	role := seedAdminRole(t, db)

	apis := []seed.ApiSpec{
		{Code: "hook", Title: "Inbound hook", Path: "/api/v1/hook", Method: "POST", Handle: "hook.Receive"},
		{Code: "sync", Title: "Sync", Path: "/api/v1/sync", Method: "GET", Handle: "hook.Sync"},
	}
	if err := (adminSeeder{}).SeedMenus(db, "hooks", nil, apis); err != nil {
		t.Fatalf("SeedMenus: %v", err)
	}

	var apiCount int64
	db.Model(&models.SysApi{}).Where("app_code = ?", "hooks").Count(&apiCount)
	if apiCount != int64(len(apis)) {
		t.Fatalf("sys_api rows = %d, want %d", apiCount, len(apis))
	}

	for _, a := range apis {
		var n int64
		db.Table("casbin_rule").
			Where("ptype = 'p' AND v0 = ? AND v1 = ? AND v2 = ?", role.RoleKey, a.Path, a.Method).
			Count(&n)
		if n != 1 {
			t.Errorf("casbin_rule for %s %s = %d rows, want 1: the endpoint is denied to admin", a.Method, a.Path, n)
		}
	}
}

// The other half of the same guard: nothing registered at all must stay a
// no-op rather than start touching sys_role_menu or casbin_rule.
func TestSeedMenusWithNothingRegisteredWritesNothing(t *testing.T) {
	db := newSeedTestDB(t)
	seedAdminRole(t, db)

	if err := (adminSeeder{}).SeedMenus(db, "empty", nil, nil); err != nil {
		t.Fatalf("SeedMenus: %v", err)
	}
	for _, table := range []string{"casbin_rule", "sys_role_menu"} {
		var n int64
		db.Table(table).Count(&n)
		if n != 0 {
			t.Errorf("%s has %d rows, want 0", table, n)
		}
	}
}

// newSeedTestDB's AutoMigrate builds a unique index on seed_code alone,
// because SysMenu.SeedCode is the only field in the struct carrying the
// uk_sys_menu_app_seed_code_del tag - app_code already carries a different,
// non-unique index name of its own, and the embedded ModelTime's
// DeletedAt (aliased from go-admin-core) cannot be given a third one. The
// real migration (cmd/migrate/migration/version/1786700008000_seed_natural_keys.go)
// never lets AutoMigrate touch this table for exactly that reason: it
// builds the composite (app_code, seed_code, deleted_at) index by hand
// instead. Reproduce that by hand here too, so a test that seeds two rows
// sharing a seed_code under different deleted_at values sees what a real
// install would, not gorm's narrower default.
func useCompositeSeedCodeIndex(t *testing.T, db *gorm.DB) {
	t.Helper()
	if db.Migrator().HasIndex(&models.SysMenu{}, "uk_sys_menu_app_seed_code_del") {
		if err := db.Migrator().DropIndex(&models.SysMenu{}, "uk_sys_menu_app_seed_code_del"); err != nil {
			t.Fatalf("drop the single-column seed_code index: %v", err)
		}
	}
	if err := db.Exec(
		"CREATE UNIQUE INDEX uk_sys_menu_app_seed_code_del ON sys_menu (app_code, seed_code, deleted_at)",
	).Error; err != nil {
		t.Fatalf("create the composite seed_code index: %v", err)
	}
}

// A retried migration - one that failed partway through and is run again,
// or simply run twice by mistake - must not create a second sys_api or
// sys_menu row for the same (appCode, natural key). This is the defect the
// demo site hit in production: duplicate sys_menu/casbin_rule rows from a
// bare tx.Create on a natural key nothing was checking.
func TestSeedMenusIsIdempotentAcrossARetry(t *testing.T) {
	db := newSeedTestDB(t)
	useCompositeSeedCodeIndex(t, db)
	seedAdminRole(t, db)

	menus := []seed.MenuSpec{
		{Code: "dir", Kind: contractmodels.Directory, Title: "Order", Sort: 10},
		{Code: "list", Parent: "dir", Kind: contractmodels.Menu, Title: "Orders", Sort: 1, ApiCodes: []string{"list"}},
	}
	apis := []seed.ApiSpec{
		{Code: "list", Title: "Order list", Path: "/api/v1/order", Method: "GET"},
	}

	run := func() {
		t.Helper()
		if err := db.Transaction(func(tx *gorm.DB) error {
			return adminSeeder{}.SeedMenus(tx, "order", menus, apis)
		}); err != nil {
			t.Fatalf("SeedMenus: %v", err)
		}
	}
	run()
	firstMenuIDs := allMenuIDs(t, db, "order")
	firstApiIDs := allApiIDs(t, db, "order")

	run() // the retry

	if got := allMenuIDs(t, db, "order"); !sameIDs(got, firstMenuIDs) {
		t.Errorf("sys_menu ids after retry = %v, want unchanged %v (a second call inserted new rows)", got, firstMenuIDs)
	}
	if got := allApiIDs(t, db, "order"); !sameIDs(got, firstApiIDs) {
		t.Errorf("sys_api ids after retry = %v, want unchanged %v (a second call inserted new rows)", got, firstApiIDs)
	}

	assertRowCount(t, db, "sys_api", 1)
	assertRowCount(t, db, "sys_menu", 2)
	assertRowCount(t, db, "sys_menu_api_rule", 1)
	assertRowCount(t, db, "sys_role_menu", 2)
	assertRowCount(t, db, "casbin_rule", 1)
}

// Only a live row counts as "already written". A row a prior, unrelated
// soft-delete already retired must not be reused - seedApis/seedMenuTree
// have to insert a fresh one under the same natural key, the same way the
// unique indexes 1786700008000_seed_natural_keys.go builds only bind live
// rows.
func TestSeedMenusOnlyReusesLiveRows(t *testing.T) {
	db := newSeedTestDB(t)
	useCompositeSeedCodeIndex(t, db)
	seedAdminRole(t, db)

	menus := []seed.MenuSpec{{Code: "dir", Kind: contractmodels.Directory, Title: "Order", Sort: 10}}
	apis := []seed.ApiSpec{{Code: "list", Title: "Order list", Path: "/api/v1/order", Method: "GET"}}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return adminSeeder{}.SeedMenus(tx, "order", menus, apis)
	}); err != nil {
		t.Fatalf("SeedMenus: %v", err)
	}

	// Soft-delete both rows this first call wrote, as if an operator (or an
	// earlier uninstall) had retired them, independently of this migration
	// ever running again.
	if err := db.Exec("UPDATE sys_menu SET deleted_at = 1").Error; err != nil {
		t.Fatalf("soft-delete sys_menu: %v", err)
	}
	if err := db.Exec("UPDATE sys_api SET deleted_at = 1").Error; err != nil {
		t.Fatalf("soft-delete sys_api: %v", err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return adminSeeder{}.SeedMenus(tx, "order", menus, apis)
	}); err != nil {
		t.Fatalf("SeedMenus after soft-delete: %v", err)
	}

	// Two rows total: the soft-deleted original, plus a fresh one - not the
	// dead row resurrected in place, and not left with zero live rows.
	assertRowCount(t, db, "sys_menu", 2)
	assertRowCount(t, db, "sys_api", 2)

	var liveMenus, liveApis int64
	db.Model(&models.SysMenu{}).Where("app_code = ?", "order").Count(&liveMenus)
	db.Model(&models.SysApi{}).Where("app_code = ?", "order").Count(&liveApis)
	if liveMenus != 1 {
		t.Errorf("live sys_menu rows = %d, want 1", liveMenus)
	}
	if liveApis != 1 {
		t.Errorf("live sys_api rows = %d, want 1", liveApis)
	}
}

// app_code is part of the natural key, not a descriptive column alongside
// it. Two applications that happen to register an identical (path, action)
// or seed_code must each get their own row - reusing one app's row for
// another's install would make an uninstall of the first delete a row the
// second considers its own.
func TestSeedMenusScopesTheNaturalKeyByAppCode(t *testing.T) {
	db := newSeedTestDB(t)
	useCompositeSeedCodeIndex(t, db)
	seedAdminRole(t, db)

	menus := []seed.MenuSpec{{Code: "dir", Kind: contractmodels.Directory, Title: "Dir", Sort: 10}}
	apis := []seed.ApiSpec{{Code: "list", Title: "Shared endpoint", Path: "/api/v1/shared", Method: "GET"}}

	for _, appCode := range []string{"order", "billing"} {
		if err := db.Transaction(func(tx *gorm.DB) error {
			return adminSeeder{}.SeedMenus(tx, appCode, menus, apis)
		}); err != nil {
			t.Fatalf("SeedMenus(%q): %v", appCode, err)
		}
	}

	var apiRows []models.SysApi
	if err := db.Where("path = ? AND action = ?", "/api/v1/shared", "GET").
		Order("app_code").Find(&apiRows).Error; err != nil {
		t.Fatalf("read sys_api: %v", err)
	}
	if len(apiRows) != 2 {
		t.Fatalf("sys_api has %d row(s) for the shared (path, action), want 2 - one per app", len(apiRows))
	}
	if apiRows[0].AppCode != "billing" || apiRows[1].AppCode != "order" {
		t.Errorf("sys_api app_codes = [%s %s], want [billing order]", apiRows[0].AppCode, apiRows[1].AppCode)
	}

	var menuRows []models.SysMenu
	if err := db.Where("seed_code = ?", "dir").Order("app_code").Find(&menuRows).Error; err != nil {
		t.Fatalf("read sys_menu: %v", err)
	}
	if len(menuRows) != 2 {
		t.Fatalf("sys_menu has %d row(s) for the shared seed_code, want 2 - one per app", len(menuRows))
	}
	if menuRows[0].AppCode != "billing" || menuRows[1].AppCode != "order" {
		t.Errorf("sys_menu app_codes = [%s %s], want [billing order]", menuRows[0].AppCode, menuRows[1].AppCode)
	}
}

func allMenuIDs(t *testing.T, db *gorm.DB, appCode string) []int {
	t.Helper()
	var rows []models.SysMenu
	if err := db.Where("app_code = ?", appCode).Order("menu_id").Find(&rows).Error; err != nil {
		t.Fatalf("read sys_menu: %v", err)
	}
	ids := make([]int, len(rows))
	for i, r := range rows {
		ids[i] = r.MenuId
	}
	return ids
}

func allApiIDs(t *testing.T, db *gorm.DB, appCode string) []int {
	t.Helper()
	var rows []models.SysApi
	if err := db.Where("app_code = ?", appCode).Order("id").Find(&rows).Error; err != nil {
		t.Fatalf("read sys_api: %v", err)
	}
	ids := make([]int, len(rows))
	for i, r := range rows {
		ids[i] = r.Id
	}
	return ids
}

func sameIDs(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func assertRowCount(t *testing.T, db *gorm.DB, table string, want int64) {
	t.Helper()
	var n int64
	if err := db.Table(table).Count(&n).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if n != want {
		t.Errorf("%s has %d row(s), want %d", table, n, want)
	}
}
