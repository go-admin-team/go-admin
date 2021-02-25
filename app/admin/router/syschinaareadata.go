package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/syschinaareadata"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysChinaAreaDataRouter)
}

// 需认证的路由代码
func registerSysChinaAreaDataRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &syschinaareadata.SysChinaAreaData{}
	r := v1.Group("/syschinaareadata").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysChinaAreaDataList)
		r.GET("/:id", api.GetSysChinaAreaData)
		r.POST("", api.InsertSysChinaAreaData)
		r.PUT("/:id", api.UpdateSysChinaAreaData)
		r.DELETE("", api.DeleteSysChinaAreaData)
	}
}
