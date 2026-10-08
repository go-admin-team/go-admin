package version

import (
	"testing"

	"go-admin/app/admin/models"
)

const (
	roleUpdateMenuId = 226
	jobEditMenuId    = 463
	jobRemoveMenuId  = 464
)

func TestBindAuthzApisGrantsWhatARoleAlreadyHolds(t *testing.T) {
	db := openGenMenuDB(t)
	for _, m := range []models.SysMenu{
		{MenuId: roleUpdateMenuId, MenuName: "修改角色", Permission: "admin:sysRole:update", MenuType: "F"},
		{MenuId: jobEditMenuId, MenuName: "修改定时任务", Permission: "job:sysJob:edit", MenuType: "F"},
		{MenuId: jobRemoveMenuId, MenuName: "删除定时任务", Permission: "job:sysJob:remove", MenuType: "F"},
	} {
		if err := db.Create(&m).Error; err != nil {
			t.Fatal(err)
		}
	}
	addRole(t, db, 2, "roleAdmin", roleUpdateMenuId)
	addRole(t, db, 3, "jobEditor", jobEditMenuId)
	addRole(t, db, 4, "jobRemover", jobRemoveMenuId)
	addRole(t, db, 5, "clerk")

	if err := bindAuthzApis(db); err != nil {
		t.Fatal(err)
	}

	for role, want := range map[string][]string{
		"roleAdmin":  {"PUT /api/v1/role-status"},
		"jobEditor":  {"GET /api/v1/job/start/:id"},
		"jobRemover": {"GET /api/v1/job/remove/:id"},
		"clerk":      nil,
	} {
		if got := policies(t, db, role); !equal(got, want) {
			t.Errorf("%s holds %v, want %v", role, got, want)
		}
	}
}

func TestBindAuthzApisReusesSeededRowAndIsIdempotent(t *testing.T) {
	db := openGenMenuDB(t)
	if err := db.Create(&models.SysMenu{MenuId: roleUpdateMenuId, Permission: "admin:sysRole:update", MenuType: "F"}).Error; err != nil {
		t.Fatal(err)
	}
	seeded := models.SysApi{Id: 153, Path: "/api/v1/role-status", Action: "PUT"}
	if err := db.Create(&seeded).Error; err != nil {
		t.Fatal(err)
	}
	addRole(t, db, 2, "roleAdmin", roleUpdateMenuId)

	for i := 0; i < 2; i++ {
		if err := bindAuthzApis(db); err != nil {
			t.Fatalf("run %d: %v", i+1, err)
		}
	}

	var n int64
	db.Model(&models.SysApi{}).Where("path = ?", seeded.Path).Count(&n)
	if n != 1 {
		t.Errorf("%d rows for %s, want the seeded one only", n, seeded.Path)
	}
	if got := boundApis(t, db, roleUpdateMenuId); got != 1 {
		t.Errorf("menu bound to %d APIs after two runs, want 1", got)
	}
	if got := policies(t, db, "roleAdmin"); len(got) != 1 {
		t.Errorf("roleAdmin holds %v after two runs, want one policy", got)
	}
}

func TestBindAuthzApisSkipsMenusThatDoNotExist(t *testing.T) {
	db := openGenMenuDB(t)
	if err := bindAuthzApis(db); err != nil {
		t.Fatal(err)
	}
	var n int64
	db.Model(&models.SysApi{}).Count(&n)
	if n != 0 {
		t.Errorf("%d APIs created with no menu to bind, want 0", n)
	}
}
