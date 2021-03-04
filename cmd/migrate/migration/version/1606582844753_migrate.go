package version

import (
	"go-admin/app/admin/models/system"

	//"go-admin/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1606582844753Test)
}

func _1606582844753Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		dept := system.SysDept{}
		err := tx.Model(&dept).Where("status = ?", 0).Update("status", 2).Error
		if err != nil {
			return err
		}

		dictData := DictData{}
		err = tx.Model(&dictData).Where("status = ?", 0).Update("status", 2).Error
		if err != nil {
			return err
		}
		sql1 := "INSERT INTO `sys_menu`(`menu_id`, `menu_name`, `title`, `icon`, `path`, `paths`, `menu_type`, `action`, `permission`, `parent_id`, `no_cache`, `breadcrumb`, `component`, `sort`, `visible`, `create_by`, `update_by`, `is_frame`, `created_at`, `updated_at`, `deleted_at`) VALUES (522, 'SysContentCreate', '新增', 'form', 'syscontent/create', '/0/498/522', 'C', '', 'syscontent:syscontent:add', 498, 0, '', '/syscontent/create.vue', 0, '1', '1', '1', '1', '2020-11-08 17:00:35.259', '2020-11-12 22:55:50.739', NULL);"
		sql2 := "INSERT INTO `sys_menu`(`menu_id`, `menu_name`, `title`, `icon`, `path`, `paths`, `menu_type`, `action`, `permission`, `parent_id`, `no_cache`, `breadcrumb`, `component`, `sort`, `visible`, `create_by`, `update_by`, `is_frame`, `created_at`, `updated_at`, `deleted_at`) VALUES (523, 'SysContentEdit', '编辑', 'edit', 'syscontent/edit:id', '/0/498/523', 'C', '', 'syscontent:syscontent:edit', 498, 0, '', '/syscontent/edit.vue', 0, '1', '1', '1', '1', '2020-11-12 22:53:42.643', '2020-11-12 22:54:56.852', NULL);"

		err = tx.Exec(sql1).Error
		if err != nil {
			return err
		}

		err = tx.Exec(sql2).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
