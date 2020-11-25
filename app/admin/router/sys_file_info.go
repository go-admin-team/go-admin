package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/sys_file"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysFileInfoRouter)
}

// 需认证的路由代码
func registerSysFileInfoRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &sys_file.SysFileInfo{}
	r := v1.Group("/sysfileinfo").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysFileInfoList)
		r.GET("/:id", api.GetSysFileInfo)
		r.POST("", api.InsertSysFileInfo)
		r.PUT("/:id", api.UpdateSysFileInfo)
		r.DELETE("/:id", api.DeleteSysFileInfo)
	}
}
