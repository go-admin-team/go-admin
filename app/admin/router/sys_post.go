package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	middleware2 "go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSyPostRouter)
}

// 需认证的路由代码
func registerSyPostRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysPost{}
	r := v1.Group("/post").Use(authMiddleware.MiddlewareFunc()).Use(middleware2.AuthCheckRole())
	{
		r.GET("", api.GetSysPostList)
		r.GET("/:id", api.GetSysPost)
		r.POST("", api.InsertSysPost)
		r.PUT("/:id", api.UpdateSysPost)
		r.DELETE("/:id", api.DeleteSysPost)
	}
}
