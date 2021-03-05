package version

import (
	//"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613720892819Test)
}

func _1613720892819Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// TODO: 这里开始写入要变更的内容
		t := &DictType{
			DictName: "内容状态",
			DictType: "sys_content_status",
			Status:   "2",
			CreateBy: "1",
			UpdateBy: "1",
		}
		err := tx.Create(t).Error
		if err != nil {
			return err
		}

		data := make([]DictData, 2)
		data[0] = DictData{
			DictSort:  0,
			DictLabel: "正常",
			DictValue: "1",
			DictType:  "sys_content_status",
			Status:    "1",
			CreateBy:  "1",
			UpdateBy:  "1",
		}
		data[1] = DictData{
			DictSort:  1,
			DictLabel: "禁用",
			DictValue: "2",
			DictType:  "sys_content_status",
			Status:    "1",
			CreateBy:  "1",
			UpdateBy:  "1",
		}
		err = tx.Create(data).Error
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
