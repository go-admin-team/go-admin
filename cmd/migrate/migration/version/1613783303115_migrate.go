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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613783303115Test)
}

func _1613783303115Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		//修改字段类型
		err := tx.Model(&models.DictType{}).Select("create_by").Not("create_by  > 0").Update("create_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.DictType{}).Select("update_by").Not("update_by > 0").Update("update_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.DictType{}).Select("create_by").Not("create_by > 0").Update("create_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.DictType{}).Select("update_by").Not("update_by > 0").Update("update_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.DictType{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.DictType{}, "update_by")
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&models.DictType{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.DictType{}, "update_by")
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
