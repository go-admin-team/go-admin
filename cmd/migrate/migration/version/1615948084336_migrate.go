package version

import (
	"runtime"
	"strings"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1615948084336Test)
}

func _1615948084336Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Exec("alter table sys_config add is_frontend int default 1 null comment '是否前台参数' after config_type").Error
		if err != nil && !strings.Contains(err.Error(), "Duplicate column name 'is_frontend'") {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
