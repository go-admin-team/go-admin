package version

import (
	//"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1614159531111Test)
}

func _1614159531111Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		sql1 := "INSERT INTO `sys_menu`(`menu_id`, `menu_name`, `title`, `icon`, `path`, `paths`, `menu_type`, `action`, `permission`, `parent_id`, `no_cache`, `breadcrumb`, `component`, `sort`, `visible`, `is_frame`, `create_by`, `update_by`, `created_at`, `updated_at`, `deleted_at`) VALUES (525, '', '公共资源目录获取', 'bug', '/api/v1/sysfiledir', '', 'A', 'GET', '', 256, 0, '', '', 0, '1', '1', 1, 1, '2021-02-24 17:36:02.900', '2021-02-24 17:36:16.600', NULL);"
		sql2 := "INSERT INTO `sys_menu`(`menu_id`, `menu_name`, `title`, `icon`, `path`, `paths`, `menu_type`, `action`, `permission`, `parent_id`, `no_cache`, `breadcrumb`, `component`, `sort`, `visible`, `is_frame`, `create_by`, `update_by`, `created_at`, `updated_at`, `deleted_at`) VALUES (526, '', '公共资源获取', 'bug', '/api/v1/sysfileinfo', '', 'A', 'GET', '', 256, 0, '', '', 0, '1', '1', 1, 0, '2021-02-24 17:36:49.575', '2021-02-24 17:36:49.575', NULL);"
		err := tx.Exec(sql1).Error
		if err != nil {
			return err
		}

		err = tx.Exec(sql2).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
