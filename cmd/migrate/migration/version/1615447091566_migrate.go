package version

import (
	"fmt"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1615447091566Test)
}

func _1615447091566Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// TODO: 这里开始写入要变更的内容

		// TODO: 例如 修改表字段 使用过程中请删除此段代码
		sql1 := "INSERT INTO `sys_menu`(`menu_id`, `menu_name`, `title`, `icon`, `path`, `paths`, `menu_type`, `action`, `permission`, `parent_id`, `no_cache`, `breadcrumb`, `component`, `sort`, `visible`, `create_by`, `update_by`, `is_frame`, `created_at`, `updated_at`, `deleted_at`) VALUES (522, 'SysContentCreate', '新增', 'form', 'syscontent/create', '/0/498/522', 'C', '', 'syscontent:syscontent:add', 498, 0, '', '/syscontent/create.vue', 0, '1', '1', '1', '1', '2020-11-08 17:00:35.259', '2020-11-12 22:55:50.739', NULL);"

		err := tx.Exec(sql1).Error
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
