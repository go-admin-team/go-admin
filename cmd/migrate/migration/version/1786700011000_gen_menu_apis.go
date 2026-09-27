package version

import (
	"errors"
	"runtime"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Bind the code generator's APIs to its two menus, and grant them to every
// role that already holds one of those menus.
//
// The generator's routes were in CasbinExclude, or mounted without
// AuthCheckRole, so any account that could log in reached them. They now go
// through AuthCheckRole. A role is granted an API through the menus it is
// given: SysRole's update writes a casbin_rule for every API bound to each of
// its menus in sys_menu_api_rule. The seed data bound none to the
// generator's menus, so without this a role given 代码生成 would be refused
// by every endpoint the page calls.
//
// Binding alone only reaches roles saved after this migration. Roles that
// already hold a generator menu are granted here, so a deployment that gave
// the generator to a role keeps it working across the upgrade. The admin
// role is let through by AuthCheckRole without a policy and needs none.
//
// Menus are found by component, not id: a deployment may have renumbered
// them, and one that deleted a menu has nobody to grant and nothing to bind.
// An API missing from sys_api is created, since its row is what the binding
// points at.
//
// Ordered after 1786700003000 (the soft-delete conversion), so the runtime
// models under app/admin/models are used, not cmd/migrate/migration/models -
// see schema_coverage_test.go's TestPostConversionMigrationsAvoidFrozenSeedModels.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700011000GenMenuApis)
}

// genMenuApis is which page calls which endpoint, read off go-admin-ui's
// src/api/tools/gen.ts and the views under dev-tools/gen.
var genMenuApis = []struct {
	component string
	apis      []models.SysApi
}{
	{"/dev-tools/gen/index", []models.SysApi{
		{Title: "代码生成表列表", Path: "/api/v1/sys/tables/page", Action: "GET"},
		{Title: "数据库表列表", Path: "/api/v1/db/tables/page", Action: "GET"},
		{Title: "数据表列列表", Path: "/api/v1/db/columns/page", Action: "GET"},
		{Title: "导入表", Path: "/api/v1/sys/tables/info", Action: "POST"},
		{Title: "删除表配置", Path: "/api/v1/sys/tables/info/:tableId", Action: "DELETE"},
		{Title: "生成预览通过id获取", Path: "/api/v1/gen/preview/:tableId", Action: "GET"},
		{Title: "数据库表生成到项目", Path: "/api/v1/gen/toproject/:tableId", Action: "GET"},
		{Title: "生成api带文件", Path: "/api/v1/gen/apitofile/:tableId", Action: "GET"},
		{Title: "数据库表生成到DB", Path: "/api/v1/gen/todb/:tableId", Action: "GET"},
	}},
	{"/dev-tools/gen/editTable", []models.SysApi{
		{Title: "表配置详情", Path: "/api/v1/sys/tables/info/:tableId", Action: "GET"},
		{Title: "按表名查询表配置", Path: "/api/v1/sys/tables/info", Action: "GET"},
		{Title: "修改表配置", Path: "/api/v1/sys/tables/info", Action: "PUT"},
		{Title: "关系表数据【代码生成】", Path: "/api/v1/gen/tabletree", Action: "GET"},
	}},
}

func _1786700011000GenMenuApis(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := bindGenMenuApis(tx); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

// bindGenMenuApis is split out so tests can run it without sys_migration.
func bindGenMenuApis(tx *gorm.DB) error {
	for _, m := range genMenuApis {
		var menu models.SysMenu
		err := tx.Where("component = ?", m.component).First(&menu).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return err
		}

		apis := make([]models.SysApi, 0, len(m.apis))
		for _, want := range m.apis {
			var api models.SysApi
			err := tx.Where("path = ? AND action = ?", want.Path, want.Action).First(&api).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				api = want
				api.Type = "SYS"
				err = tx.Create(&api).Error
			}
			if err != nil {
				return err
			}
			apis = append(apis, api)
		}
		if err := tx.Model(&menu).Association("SysApi").Append(apis); err != nil {
			return err
		}

		var roleKeys []string
		if err := tx.Model(&models.SysRole{}).
			Joins("JOIN sys_role_menu ON sys_role_menu.role_id = sys_role.role_id").
			Where("sys_role_menu.menu_id = ?", menu.MenuId).
			Distinct().Pluck("sys_role.role_key", &roleKeys).Error; err != nil {
			return err
		}
		for _, key := range roleKeys {
			for _, a := range apis {
				if err := tx.Exec(
					"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) SELECT 'p', ?, ?, ?, '', '', '' WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=? AND v1=? AND v2=?)",
					key, a.Path, a.Action, key, a.Path, a.Action,
				).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
