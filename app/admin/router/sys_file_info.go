package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis/sys_file"
	middleware2 "go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysFileInfoRouter)
}

// 需认证的路由代码
func registerSysFileInfoRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := sys_file.SysFileInfo{}
	r := v1.Group("/sysfileinfo").Use(authMiddleware.MiddlewareFunc()).Use(middleware2.AuthCheckRole())
	{
		r.GET("", api.GetSysFileInfoList)
		r.GET("/:id", api.GetSysFileInfo)
		r.POST("", api.InsertSysFileInfo)
		r.PUT("/:id", api.UpdateSysFileInfo)
		r.DELETE("/:id", api.DeleteSysFileInfo)
	}
}
