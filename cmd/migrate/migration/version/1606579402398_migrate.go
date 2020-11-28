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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1606579402398Test)
}

func _1606579402398Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		dicData := models.DictData{}
		err := tx.Model(&dicData).Where("dict_code = ?", 1).Update("dict_value", 2).Error
		if err != nil {
			return err
		}

		user := models.SysUser{}
		err = tx.Model(&user).Where("user_id = ?", 1).Update("status", 2).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
