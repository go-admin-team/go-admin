package version

import (
	"runtime"
	"strconv"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
)

// 为 app/demo 模块写入菜单、接口与权限种子数据。
//
// 一个业务模块要在界面上可用，需要四类数据：
//
//  1. sys_api        —— 后端路由的登记，Casbin 据此判定权限
//  2. sys_menu       —— 侧边栏菜单（目录 M / 菜单 C / 按钮 F）
//  3. sys_menu_api_rule —— 菜单与接口的多对多关联，角色保存时据此生成策略
//  4. casbin_rule    —— 实际生效的权限策略
//
// 注意 casbin_rule 才是 adapter 实际使用的表；models.CasbinRule 对应的
// sys_casbin_rule 是历史遗留，其 7 列 size:512 唯一索引在 MySQL 下会超出
// 索引长度限制，不要将其纳入 AutoMigrate。
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700001000DemoMenu)
}

const (
	demoMenuId    = 9000 // 目录：示例模块
	demoProductId = 9001 // 菜单：商品管理
	demoApiBaseId = 9000 // sys_api 起始 id
	adminRoleKey  = "admin"
)

func _1786700001000DemoMenu(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// ---- 1. 接口登记 ----
		apis := []models.SysApi{
			{Id: demoApiBaseId + 1, Handle: "go-admin/common/actions.IndexAction.func1", Title: "示例商品列表", Path: "/api/v1/demo-product", Type: "SYS", Action: "GET"},
			{Id: demoApiBaseId + 2, Handle: "go-admin/common/actions.ViewAction.func1", Title: "示例商品详情", Path: "/api/v1/demo-product/:id", Type: "SYS", Action: "GET"},
			{Id: demoApiBaseId + 3, Handle: "go-admin/common/actions.CreateAction.func1", Title: "示例商品新增", Path: "/api/v1/demo-product", Type: "SYS", Action: "POST"},
			{Id: demoApiBaseId + 4, Handle: "go-admin/common/actions.UpdateAction.func1", Title: "示例商品修改", Path: "/api/v1/demo-product/:id", Type: "SYS", Action: "PUT"},
			{Id: demoApiBaseId + 5, Handle: "go-admin/common/actions.DeleteAction.func1", Title: "示例商品删除", Path: "/api/v1/demo-product", Type: "SYS", Action: "DELETE"},
		}
		for i := range apis {
			if err := upsert(tx, &models.SysApi{}, "id = ?", apis[i].Id, &apis[i]); err != nil {
				return err
			}
		}

		// ---- 2. 菜单 ----
		// 目录本身不对应页面，component 固定为 Layout
		dir := models.SysMenu{
			MenuId: demoMenuId, MenuName: "Demo", Title: "示例模块", Icon: "example",
			Path: "/demo", Paths: "/0/9000", MenuType: "M", ParentId: 0,
			// sort is `gorm:"size:4"`, which MySQL builds as a tinyint - anything
			// over 127 is rejected outright. The seeded menus run to 100, so 110
			// still puts this last.
			Component: "Layout", Sort: 110, Visible: "0", IsFrame: "1",
		}
		if err := upsert(tx, &models.SysMenu{}, "menu_id = ?", dir.MenuId, &dir); err != nil {
			return err
		}

		// 菜单的 menu_name 需与前端组件 name 一致，否则 keep-alive 缓存无法命中
		page := models.SysMenu{
			MenuId: demoProductId, MenuName: "DemoProduct", Title: "商品管理", Icon: "documentation",
			Path: "/demo/product", Paths: "/0/9000/9001", MenuType: "C", ParentId: demoMenuId,
			Component: "/demo/product/index", Sort: 1, Visible: "0", IsFrame: "1",
			SysApi: apis,
		}
		if err := upsert(tx, &models.SysMenu{}, "menu_id = ?", page.MenuId, &page); err != nil {
			return err
		}

		// 按钮级权限，permission 需与前端 v-permisaction 中的标识一致
		buttons := []models.SysMenu{
			{MenuId: 9002, MenuName: "DemoProductAdd", Title: "新增", MenuType: "F", ParentId: demoProductId, Permission: "demo:product:add", Sort: 1, Visible: "0", IsFrame: "1"},
			{MenuId: 9003, MenuName: "DemoProductEdit", Title: "修改", MenuType: "F", ParentId: demoProductId, Permission: "demo:product:edit", Sort: 2, Visible: "0", IsFrame: "1"},
			{MenuId: 9004, MenuName: "DemoProductDel", Title: "删除", MenuType: "F", ParentId: demoProductId, Permission: "demo:product:delete", Sort: 3, Visible: "0", IsFrame: "1"},
		}
		for i := range buttons {
			buttons[i].Paths = "/0/9000/9001/" + strconv.Itoa(buttons[i].MenuId)
			if err := upsert(tx, &models.SysMenu{}, "menu_id = ?", buttons[i].MenuId, &buttons[i]); err != nil {
				return err
			}
		}

		// ---- 3. 授权给 admin 角色 ----
		var role models.SysRole
		err := tx.Where("role_key = ?", adminRoleKey).First(&role).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 尚未初始化角色数据时跳过授权，不阻断迁移
				return tx.Create(&common.Migration{Version: version}).Error
			}
			return err
		}

		for _, id := range []int{demoMenuId, demoProductId, 9002, 9003, 9004} {
			if err = tx.Exec(
				"INSERT INTO sys_role_menu (role_id, menu_id) SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM sys_role_menu WHERE role_id = ? AND menu_id = ?)",
				role.RoleId, id, role.RoleId, id,
			).Error; err != nil {
				return err
			}
		}

		// ---- 4. Casbin 策略 ----
		// admin 角色在中间件中直接放行，此处仍写入策略以便复制该角色配置时可继承
		for _, a := range apis {
			if err = tx.Exec(
				"INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) SELECT 'p', ?, ?, ?, '', '', '' WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0=? AND v1=? AND v2=?)",
				role.RoleKey, a.Path, a.Action, role.RoleKey, a.Path, a.Action,
			).Error; err != nil {
				return err
			}
		}

		return tx.Create(&common.Migration{Version: version}).Error
	})
}

// upsert 存在则更新、不存在则插入。迁移可能在已有数据的库上运行，
// 直接 Create 会因主键冲突失败。
func upsert(tx *gorm.DB, model interface{}, where string, id int, value interface{}) error {
	var count int64
	if err := tx.Model(model).Where(where, id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return tx.Model(model).Where(where, id).Updates(value).Error
	}
	return tx.Create(value).Error
}
