package version

import (
	"go-admin/app/admin/models"
	//"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1606582844753Test)
}

func _1606582844753Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {


		dept := models.SysDept{}
		err := tx.Model(&dept).Where("status = ?", 0).Update("status", 2).Error
		if err != nil {
			return err
		}

		dictData := models.DictData{}
		err = tx.Model(&dictData).Where("status = ?", 0).Update("status", 2).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
