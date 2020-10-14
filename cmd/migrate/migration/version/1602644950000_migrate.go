package version

import (
	"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1602644950000Test)
}

func _1602644950000Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := db.Migrator().RenameColumn(&models.SysConfig{}, "config_id", "id")
		if err != nil {
			return err
		}
		return db.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
