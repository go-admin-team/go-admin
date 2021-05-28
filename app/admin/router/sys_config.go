package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysConfigRouter)
}

// 需认证的路由代码
func registerSysConfigRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysConfig{}
	r := v1.Group("/config").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysConfigList)
		r.GET("/:id", api.GetSysConfig)
		r.POST("", api.InsertSysConfig)
		r.PUT("/:id", api.UpdateSysConfig)
		r.DELETE("", api.DeleteSysConfig)
	}

	r1 := v1.Group("/configKey").Use(authMiddleware.MiddlewareFunc())
	{
		r1.GET("/:configKey", api.GetSysConfigByKEYForService)
	}

	r2 := v1.Group("/app-config")
	{
		r2.GET("", api.GetSysConfigBySysApp)
	}

	r3 := v1.Group("/set-config")
	{
		r3.PUT("", api.UpdateSetSysConfig)
		r3.GET("", api.GetSetSysConfig)
	}

}