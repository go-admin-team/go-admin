package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"

	"go-admin/app/demo/models"
	"go-admin/app/demo/service/dto"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

// 路由通过 init 自注册，无需在任何中心文件登记。
// 新建应用时用 `go run main.go app -n <名称>` 生成骨架，
// 它会同时产出 cmd/api/<名称>.go 完成注册。
func init() {
	routerCheckRole = append(routerCheckRole, registerDemoProductRouter)
}

// registerDemoProductRouter 标准 CRUD 的推荐写法。
//
// 五个通用 Action 覆盖了增删改查的全部样板逻辑——参数绑定、数据权限过滤、
// 操作人注入、分页、错误响应，因此本模块没有 apis 与 service 文件。
//
// 仅当业务逻辑超出单表 CRUD（如跨表事务、外部调用、复杂校验）时，才需要
// 自行编写 Handler 与 Service，写法参照 app/admin/apis/sys_post.go。
func registerDemoProductRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/demo-product").
		Use(authMiddleware.MiddlewareFunc()). // JWT 认证
		Use(middleware.AuthCheckRole())       // Casbin 鉴权
	{
		m := &models.DemoProduct{}

		// actions.PermissionAction() 注入数据权限上下文，
		// 列表与详情缺少它会绕过 DataScope 过滤
		r.GET("", actions.PermissionAction(), actions.IndexAction(m, new(dto.DemoProductSearch), func() interface{} {
			list := make([]models.DemoProduct, 0)
			return &list
		}))

		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto.DemoProductById), func() interface{} {
			return &models.DemoProduct{}
		}))

		r.POST("", actions.CreateAction(new(dto.DemoProductControl)))
		r.PUT("/:id", actions.PermissionAction(), actions.UpdateAction(new(dto.DemoProductControl)))
		r.DELETE("", actions.PermissionAction(), actions.DeleteAction(new(dto.DemoProductById)))
	}
}
