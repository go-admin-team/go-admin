package service

import (
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

// An update the data permission excludes has to be refused, and refused in a
// way that does not tell the caller whether the row exists. First reports both
// cases the same way - no rows - so the message has to come from there rather
// than from a RowsAffected check the error return has already skipped past.
func TestSysApiUpdateRefusesARowOutsideTheDataPermission(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:sysapi-perm?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Skipf("sqlite unavailable: %v", err)
	}
	if err := db.AutoMigrate(&models.SysApi{}); err != nil {
		t.Skipf("automigrate: %v", err)
	}

	prev := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() { config.ApplicationConfig.EnableDP = prev })

	// Owned by user 1.
	row := models.SysApi{Handle: "h", Title: "t", Path: "/api/v1/probe", Type: "BUS", Action: "GET"}
	row.CreateBy = 1
	if err := db.Create(&row).Error; err != nil {
		t.Fatal(err)
	}

	e := &SysApi{Service: service.Service{Orm: db, Log: logger.NewHelper(logger.DefaultLogger)}}
	req := &dto.SysApiUpdateReq{Id: row.Id, Title: "changed"}

	// User 2, scope 5: only rows they created.
	outsider := &actions.DataPermission{DataScope: "5", UserId: 2, DeptId: 1, RoleId: 2}
	err = e.Update(req, outsider)
	if err == nil {
		t.Fatal("the update was allowed on a row the data permission excludes")
	}
	if !strings.Contains(err.Error(), "无权更新该数据") {
		t.Errorf("refused with %q, want the permission message; a raw database error tells the "+
			"caller the row exists", err)
	}

	var after models.SysApi
	if err := db.First(&after, row.Id).Error; err != nil {
		t.Fatal(err)
	}
	if after.Title != "t" {
		t.Errorf("the row was modified: title is now %q", after.Title)
	}

	// The owner still gets through, so the scope is refusing rather than
	// everything failing.
	owner := &actions.DataPermission{DataScope: "5", UserId: 1, DeptId: 1, RoleId: 1}
	if err := e.Update(&dto.SysApiUpdateReq{Id: row.Id, Title: "by owner"}, owner); err != nil {
		t.Fatalf("the owner could not update their own row: %v", err)
	}
}
