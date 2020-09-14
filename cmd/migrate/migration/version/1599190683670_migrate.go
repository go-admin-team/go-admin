package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683670Test)
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
