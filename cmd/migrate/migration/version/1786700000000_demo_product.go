package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/app/demo/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// 示例模块的建表迁移。
//
// 迁移文件名前 13 位是版本号（Unix 毫秒时间戳），框架按文件名升序执行，
// 已执行的版本记录在 sys_migration 表中，不会重复运行。
//
// 本文件放在 version/ 是因为它随框架一起分发；**业务项目自己的迁移应放
// version-local/**，该目录已被 .gitignore 忽略，不会与上游冲突。
//
// 生成新的迁移骨架：go run main.go migrate -c config/settings.yml -g
//
// 注意：已执行过的迁移文件不可再修改——版本号已入表，改动不会重跑。
// 需要调整时新建一个迁移。
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700000000DemoProduct)
}

func _1786700000000DemoProduct(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 建表
		if err := tx.Migrator().AutoMigrate(new(models.DemoProduct)); err != nil {
			return err
		}

		// 2. 初始化数据（可选）。此处留空，示例模块不写入业务数据。

		// 3. 记录版本号——必须，否则该迁移每次启动都会重跑
		return tx.Create(&common.Migration{Version: version}).Error
	})
}
