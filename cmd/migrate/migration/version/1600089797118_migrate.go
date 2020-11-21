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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1600089797118Migrate)
}

func _1600089797118Migrate(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		f := &models.SysFileDir{
			Label: "根目录",
			PId:   0,
			//Sort:  "0",
			Path:  "",
			ControlBy: common.ControlBy{
				CreateBy: 1,
				UpdateBy: 1,
			},
		}
		err := tx.Create(f).Error
		if err != nil {
			return err
		}
		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
