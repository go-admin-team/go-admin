package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis/system/sys_menu"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysMenuRouter)
}

// 需认证的路由代码
func registerSysMenuRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := sys_menu.SysMenu{}
	r := v1.Group("/menu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysMenuList)
		r.GET("/:id", api.GetSysMenu)
		r.POST("", api.InsertSysMenu)
		r.PUT("/:id", api.UpdateSysMenu)
		r.DELETE("/:id", api.DeleteSysMenu)
	}

	r1 := v1.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		r1.GET("/menurole", api.GetMenuRole)
		r1.GET("/menuids", api.GetMenuIDS)
	}

}