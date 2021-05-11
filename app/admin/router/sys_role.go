package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysRoleRouter)
}

// 需认证的路由代码
func registerSysRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysRole{}
	r := v1.Group("/role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysRoleList)
		r.GET("/:id", api.GetSysRole)
		r.POST("", api.InsertSysRole)
		r.PUT("/:id", api.UpdateSysRole)
		r.DELETE("", api.DeleteSysRole)
	}
	v1.PUT("/roledatascope", api.UpdateRoleDataScope)
}
