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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613998164781Test)
}

func _1613998164781Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/deptList").Update("path", "/api/v1/dept").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/roleList").Update("path", "/api/v1/role").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/menuList").Update("path", "/api/v1/menu").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/postList").Update("path", "/api/v1/post").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/dict/typelist").Update("path", "/api/v1/dict/type").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/dict/datalist").Update("path", "/api/v1/dict/data").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path like ?", "/api/v1/dict/data").
			Where("action = ?", "PUT").Update("path", "/api/v1/dict/data/:id").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path like ?", "/api/v1/dict/data/:id").
			Where("action = ?", "DELETE").Update("path", "/api/v1/dict/data").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path like ?", "/api/v1/dict/type").
			Where("action = ?", "PUT").Update("path", "/api/v1/dict/type/:id").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path like ?", "/api/v1/dict/type/:id").
			Where("action = ?", "DELETE").Update("path", "/api/v1/dict/type").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/operloglist").Update("path", "/api/v1/sys-opera-log").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/loginloglist").Update("path", "/api/v1/sys-login-log").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path = ?", "/api/v1/loginlog/:id").Where("action = ?", "DELETE").Update("path", "/api/v1/sys-opera-log").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path like ?", "/api/v1/operloglist%").Where("action = ?", "DELETE").Update("path", "/api/v1/sys-opera-log").Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.Menu{}).Where("path like ?", "/api/v1/role/:id").Where("action = ?", "DELETE").Update("path", "/api/v1/role").Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
