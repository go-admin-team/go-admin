package version

import (
	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"

	//"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1614604763713Test)
}

func _1614604763713Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/dict/typeoptionselect").Update("path", "/api/v1/dict/type-option-select").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("action = ?", "GET").
			Where("path = ?", "/api/v1/sysUserList").Update("path", "/api/v1/sysUser").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("action = ?", "DELETE").
			Where("path = ?", "/api/v1/sysUser/:id").Update("path", "/api/v1/sysUser").Error
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&system.SysUser{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysUser{}, "update_by")
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&models.SysContent{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.SysContent{}, "update_by")
		if err != nil {
			return err
		}
		err = tx.Model(&models.Menu{}).Where("action = ?", "PUT").
			Where("path = ?", "/api/v1/syscontent").Update("path", "/api/v1/syscontent/:id").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("action = ?", "DELETE").
			Where("path = ?", "/api/v1/syscontent/:id").Update("path", "/api/v1/syscontent").Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
