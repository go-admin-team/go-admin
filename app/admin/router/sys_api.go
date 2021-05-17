package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysApiRouter)
}

// registerSysApiRouter
func registerSysApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysApi{}
	r := v1.Group("/sys-api").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysApiList)
		r.GET("/:id", api.GetSysApi)
		//r.POST("", api.InsertSysApi)
		r.PUT("/:id", api.UpdateSysApi)
		//r.DELETE("", api.DeleteSysApi)
	}
}