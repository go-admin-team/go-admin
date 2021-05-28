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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613697961697Test)
}

func _1613697961697Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		//修改字段类型
		err := tx.Migrator().AlterColumn(&models.SysMenu{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.SysMenu{}, "update_by")
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&models.SysRole{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.SysRole{}, "update_by")
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
