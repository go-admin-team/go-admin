package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis/sys_file"
	middleware2 "go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysFileDirRouter)
}

// 需认证的路由代码
func registerSysFileDirRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := sys_file.SysFileDir{}
	r := v1.Group("/sysfiledir").Use(authMiddleware.MiddlewareFunc()).Use(middleware2.AuthCheckRole())
	{
		r.GET("", api.GetSysFileDirList)
		r.GET("/:id", api.GetSysFileDir)
		r.POST("", api.InsertSysFileDir)
		r.PUT("/:id", api.UpdateSysFileDir)
		r.DELETE("/:id", api.DeleteSysFileDir)
	}
}
