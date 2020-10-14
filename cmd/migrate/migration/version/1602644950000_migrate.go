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
		list2 := []models.CasbinRule{
			{PType: "p", V0: "admin", V1: "/api/v1/config", V2: "GET"},
		}
		err = tx.Create(list2).Error
		if err != nil {
			return err
		}

		menu := models.Menu{MenuId: 86, Path: "/api/v1/config"}
		err = tx.Model(&menu).Where("menu_id = ?", 86).Update("Path", "/api/v1/config").Error
		if err != nil {
			return err
		}

		return db.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
