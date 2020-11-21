package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/sys_file"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysFileDirRouter)
}

// 需认证的路由代码
func registerSysFileDirRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &sys_file.SysFileDir{}
	r := v1.Group("/sysfiledir").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysFileDirList)
		r.GET("/:id", api.GetSysFileDir)
		r.POST("", api.InsertSysFileDir)
		r.PUT("/:id", api.UpdateSysFileDir)
		r.DELETE("/:id", api.DeleteSysFileDir)
	}
}
