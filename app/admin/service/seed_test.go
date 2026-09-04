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
