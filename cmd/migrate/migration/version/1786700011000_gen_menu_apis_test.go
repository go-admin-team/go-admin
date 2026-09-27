package version

import (
	"sort"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-admin/app/admin/models"
)

const (
	genMenuId  = 261
	editMenuId = 262
)

// openGenMenuDB builds the tables bindGenMenuApis reads and writes, with the
// two generator menus under the ids config/db.sql gives them.
func openGenMenuDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(new(models.SysApi), new(models.SysMenu), new(models.SysRole)); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE casbin_rule (id integer primary key autoincrement,
		ptype text, v0 text, v1 text, v2 text, v3 text, v4 text, v5 text)`).Error; err != nil {
		t.Fatal(err)
	}
	for _, m := range []models.SysMenu{
		{MenuId: genMenuId, MenuName: "Gen", Component: "/dev-tools/gen/index", MenuType: "C"},
		{MenuId: editMenuId, MenuName: "EditTable", Component: "/dev-tools/gen/editTable", MenuType: "C"},
	} {
		if err := db.Create(&m).Error; err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func addRole(t *testing.T, db *gorm.DB, id int, key string, menus ...int) {
	t.Helper()
	if err := db.Create(&models.SysRole{RoleId: id, RoleKey: key, RoleName: key}).Error; err != nil {
		t.Fatal(err)
	}
	for _, m := range menus {
		if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", id, m).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func policies(t *testing.T, db *gorm.DB, role string) []string {
	t.Helper()
	var rows []struct{ V1, V2 string }
	if err := db.Raw("SELECT v1, v2 FROM casbin_rule WHERE ptype = 'p' AND v0 = ?", role).Scan(&rows).Error; err != nil {
		t.Fatal(err)
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.V2+" "+r.V1)
	}
	sort.Strings(out)
	return out
}

func boundApis(t *testing.T, db *gorm.DB, menuId int) int {
	t.Helper()
	menu := models.SysMenu{MenuId: menuId}
	return int(db.Model(&menu).Association("SysApi").Count())
}

func TestGenMenuApisGrantsWhatARoleAlreadyHolds(t *testing.T) {
	db := openGenMenuDB(t)
	addRole(t, db, 2, "developer", genMenuId, editMenuId)
	addRole(t, db, 3, "editor", editMenuId)
	addRole(t, db, 4, "clerk")

	if err := bindGenMenuApis(db); err != nil {
		t.Fatal(err)
	}

	if n := boundApis(t, db, genMenuId); n != 9 {
		t.Errorf("代码生成 is bound to %d APIs, want 9", n)
	}
	if n := boundApis(t, db, editMenuId); n != 4 {
		t.Errorf("代码生成修改 is bound to %d APIs, want 4", n)
	}
	if got := policies(t, db, "developer"); len(got) != 13 {
		t.Errorf("developer holds %d policies, want 13: %v", len(got), got)
	}
	want := []string{
		"GET /api/v1/gen/tabletree",
		"GET /api/v1/sys/tables/info",
		"GET /api/v1/sys/tables/info/:tableId",
		"PUT /api/v1/sys/tables/info",
	}
	if got := policies(t, db, "editor"); !equal(got, want) {
		t.Errorf("editor holds %v, want %v", got, want)
	}
	if got := policies(t, db, "clerk"); len(got) != 0 {
		t.Errorf("clerk, who holds neither menu, was granted %v", got)
	}
}

// The seed data already registers most of these APIs. Binding must point at
// those rows, not add a second row per API next to each.
func TestGenMenuApisReusesRegisteredApis(t *testing.T) {
	db := openGenMenuDB(t)
	seeded := models.SysApi{Id: 32, Title: "数据库表生成到项目", Path: "/api/v1/gen/toproject/:tableId", Action: "GET", Type: "SYS"}
	if err := db.Create(&seeded).Error; err != nil {
		t.Fatal(err)
	}

	if err := bindGenMenuApis(db); err != nil {
		t.Fatal(err)
	}

	var n int64
	db.Model(&models.SysApi{}).Where("path = ? AND action = ?", seeded.Path, seeded.Action).Count(&n)
	if n != 1 {
		t.Errorf("%d rows for %s, want the seeded one only", n, seeded.Path)
	}
	var total int64
	db.Model(&models.SysApi{}).Count(&total)
	if total != 13 {
		t.Errorf("sys_api holds %d rows, want 13: the 12 missing created and the seeded one reused", total)
	}
}

func TestGenMenuApisRunTwiceAddsNothing(t *testing.T) {
	db := openGenMenuDB(t)
	addRole(t, db, 2, "developer", genMenuId, editMenuId)

	for i := 0; i < 2; i++ {
		if err := bindGenMenuApis(db); err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
	}
	if n := boundApis(t, db, genMenuId); n != 9 {
		t.Errorf("代码生成 is bound to %d APIs after two runs, want 9", n)
	}
	if got := policies(t, db, "developer"); len(got) != 13 {
		t.Errorf("developer holds %d policies after two runs, want 13", len(got))
	}
}

func TestGenMenuApisSkipsADeletedMenuAndADeletedRole(t *testing.T) {
	db := openGenMenuDB(t)
	if err := db.Delete(&models.SysMenu{MenuId: editMenuId}).Error; err != nil {
		t.Fatal(err)
	}
	addRole(t, db, 5, "gone", genMenuId)
	if err := db.Model(&models.SysRole{}).Where("role_id = ?", 5).
		Update("deleted_at", time.Now().UnixMilli()).Error; err != nil {
		t.Fatal(err)
	}

	if err := bindGenMenuApis(db); err != nil {
		t.Fatal(err)
	}
	if n := boundApis(t, db, editMenuId); n != 0 {
		t.Errorf("a deleted menu was bound to %d APIs", n)
	}
	if got := policies(t, db, "gone"); len(got) != 0 {
		t.Errorf("a deleted role was granted %v", got)
	}
}

func equal(a, b []string) bool {
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
