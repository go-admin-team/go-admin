package version

import (
	"go-admin/cmd/migrate/migration/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1619418811873Test)
}

func _1619418811873Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// TODO: 这里开始写入要变更的内容
		err := tx.Migrator().AlterColumn(&models.SysRole{}, "CreateBy")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.SysRole{}, "UpdateBy")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.SysJob{}, "CreateBy")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.SysJob{}, "UpdateBy")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.RoleMenu{}, "CreateBy")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&models.RoleMenu{}, "UpdateBy")
		if err != nil {
			return err
		}

		// TODO: 例如 修改表字段 使用过程中请删除此段代码
		//err := tx.Migrator().RenameColumn(&models.SysConfig{}, "config_id", "id")
		//if err != nil {
		// 	return err
		//}

		// TODO: 例如 新增表结构 使用过程中请删除此段代码
		//err = tx.Debug().Migrator().AutoMigrate(
		//		new(models.CasbinRule),
		// 		)
		//if err != nil {
		// 	return err
		//}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
