package version

import (
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1615961953379Test)
}

func _1615961953379Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Exec("INSERT INTO `sys_config` (`id`, `config_name`, `config_key`, `config_value`, `config_type`, `is_frontend`, `remark`, `create_by`, `update_by`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '系统名称', 'sys_app_name', 'go-admin管理系统', 'Y', 1, '', 1, 0, '2021-03-17 08:52:06.067', '2021-03-17 08:52:06.067', NULL);").Error
		if err != nil {
			return err
		}

		err = tx.Exec("INSERT INTO `sys_config` (`id`, `config_name`, `config_key`, `config_value`, `config_type`, `is_frontend`, `remark`, `create_by`, `update_by`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, '系统logo', 'sys_app_logo', 'https://gitee.com/mydearzwj/image/raw/master/img/go-admin.png', 'Y', 1, '', 1, 0, '2021-03-17 08:53:19.462', '2021-03-17 08:53:19.462', NULL);").Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
