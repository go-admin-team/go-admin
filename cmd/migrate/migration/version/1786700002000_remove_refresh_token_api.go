package version

import (
	"runtime"

	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	"go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
)

// 清理 GET /api/v1/refresh_token 的残留权限数据。
//
// 该路由已随 issue #820 的修复移除：它允许用业务 token 换取新 token，而续期
// 上限 MaxRefresh 依据的 orig_iat 每次续期都被重置，上限永远无法到达。
//
// 路由虽已删除（请求会返回 404），但已有部署的库中仍残留三类记录：接口登记、
// 菜单与接口的绑定、Casbin 策略。留着会让「接口管理」列出一个不存在的端点，
// 角色配置里也仍可勾选，造成误解。
//
// 按 path 匹配而非固定 id —— 使用者若执行过 `server -a` 重新注册接口，id 会
// 与官方种子数据不同。
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700002000RemoveRefreshTokenApi)
}

const refreshTokenPath = "/api/v1/refresh_token"

func _1786700002000RemoveRefreshTokenApi(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var apiIds []int
		if err := tx.Model(&models.SysApi{}).
			Where("path = ? AND action = ?", refreshTokenPath, "GET").
			Pluck("id", &apiIds).Error; err != nil {
			return err
		}

		if len(apiIds) > 0 {
			// 先断开绑定，再删接口本身，避免遗留悬空外键
			// 连接表由 GORM many2many 生成，列名为 sys_api_id 而非 api_id
			if err := tx.Exec("DELETE FROM sys_menu_api_rule WHERE sys_api_id IN ?", apiIds).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", apiIds).Delete(&models.SysApi{}).Error; err != nil {
				return err
			}
		}

		// casbin_rule 才是 adapter 实际使用的表，策略按路径存储，与 sys_api 的
		// id 无关，因此即使上面没匹配到接口也要清理
		if err := tx.Exec(
			"DELETE FROM casbin_rule WHERE ptype = 'p' AND v1 = ?", refreshTokenPath,
		).Error; err != nil {
			return err
		}

		return tx.Create(&common.Migration{Version: version}).Error
	})
}
