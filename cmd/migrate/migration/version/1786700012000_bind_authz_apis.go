package version

import (
	"errors"
	"runtime"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Bind the APIs that were mounted without AuthCheckRole to the menus that
// expose them, and grant them to every role that already holds one of those
// menus.
//
// /job/start, /job/remove and /role-status now go through AuthCheckRole. The
// seed data registered them in sys_api but bound none to a menu, so a role
// saved through the menu tree could never be granted them. Binding only
// reaches roles saved after this migration; roles that already hold a menu
// are granted here so an upgrade does not take the buttons away from them.
// The admin role is let through by AuthCheckRole and needs no policy.
//
// Menus are found by permission, not id, for the same reason as in
// 1786700011000. start has no menu of its own in the seed, so it follows
// job:sysJob:edit as well as job:sysJob:start where a deployment added one.
//
// Ordered after 1786700003000, so the runtime models are used.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700012000BindAuthzApis)
}

var authzMenuApis = []struct {
	permissions []string
	api         models.SysApi
}{
	{[]string{"admin:sysRole:update"}, models.SysApi{Title: "角色状态修改", Path: "/api/v1/role-status", Action: "PUT"}},
	{[]string{"job:sysJob:remove"}, models.SysApi{Title: "job移除", Path: "/api/v1/job/remove/:id", Action: "GET"}},
	{[]string{"job:sysJob:start", "job:sysJob:edit"}, models.SysApi{Title: "job启动", Path: "/api/v1/job/start/:id", Action: "GET"}},
}

func _1786700012000BindAuthzApis(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := bindAuthzApis(tx); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

// bindAuthzApis is split out so tests can run it without sys_migration.
func bindAuthzApis(tx *gorm.DB) error {
	for _, b := range authzMenuApis {
		var menus []models.SysMenu
		if err := tx.Where("permission IN ?", b.permissions).Find(&menus).Error; err != nil {
			return err
		}
		if len(menus) == 0 {
			continue
		}

		var api models.SysApi
		err := tx.Where("path = ? AND action = ?", b.api.Path, b.api.Action).First(&api).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api = b.api
			api.Type = "BUS"
			err = tx.Create(&api).Error
		}
		if err != nil {
			return err
		}

		for i := range menus {
			if err := tx.Model(&menus[i]).Association("SysApi").Append(&api); err != nil {
				return err
			}

			var roleKeys []string
			if err := tx.Model(&models.SysRole{}).
				Joins("JOIN sys_role_menu ON sys_role_menu.role_id = sys_role.role_id").
				Where("sys_role_menu.menu_id = ?", menus[i].MenuId).
				Distinct().Pluck("sys_role.role_key", &roleKeys).Error; err != nil {
				return err
			}
			for _, key := range roleKeys {
				if err := tx.Exec(
					"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) SELECT 'p', ?, ?, ?, '', '', '' WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=? AND v1=? AND v2=?)",
					key, api.Path, api.Action, key, api.Path, api.Action,
				).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
