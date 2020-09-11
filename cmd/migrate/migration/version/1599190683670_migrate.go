package version

import (
	"go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
	"gorm.io/gorm"
	"path/filepath"
	"runtime"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	fileName = filepath.Base(fileName)
	fileName = fileName[:len(fileName)-3]
	migration.Migrate.SetVersion(fileName, _1599190683670Test)
}

func _1599190683670Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		if err := models.InitDb(tx); err != nil {

		}


		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
