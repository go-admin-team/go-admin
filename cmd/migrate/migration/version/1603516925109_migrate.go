package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1603516925109Test)
}

func _1603516925109Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Migrator().RenameTable("sys_operlog", "sys_opera_log")
		_ = tx.Migrator().RenameTable("sys_loginlog", "sys_login_log")

		if tx.Migrator().HasColumn(&models.SysLoginLog{}, "info_id") {
			err := tx.Migrator().RenameColumn(&models.SysLoginLog{}, "info_id", "id")
			if err != nil {
				return err
			}
		}

		if tx.Migrator().HasColumn(&models.SysOperaLog{}, "oper_id") {
			err := tx.Migrator().RenameColumn(&models.SysOperaLog{}, "oper_id", "id")
			if err != nil {
				return err
			}
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
