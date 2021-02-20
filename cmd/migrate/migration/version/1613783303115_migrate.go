package version

import (
	"go-admin/app/admin/models/system"

	//"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613783303115Test)
}

func _1613783303115Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		//修改字段类型
		err := tx.Migrator().AlterColumn(&system.SysDictType{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysDictType{}, "update_by")
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&system.SysDictData{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysDictData{}, "update_by")
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
