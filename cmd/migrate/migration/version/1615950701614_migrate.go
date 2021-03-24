package version

import (
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1615950701614Test)
}

func _1615950701614Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Exec("UPDATE sys_dict_data t SET t.dict_value = '2' WHERE t.dict_code = 16").Error
		if err != nil {
			return err
		}

		err = tx.Exec("UPDATE sys_opera_log t SET t.status = '2' WHERE t.status = '0'").Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
