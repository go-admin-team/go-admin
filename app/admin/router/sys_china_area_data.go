package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis/sys_china_area_data"
	middleware2 "go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysChinaAreaDataRouter)
}

// 需认证的路由代码
func registerSysChinaAreaDataRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := sys_china_area_data.SysChinaAreaData{}
	r := v1.Group("/sys_china_area_data").Use(authMiddleware.MiddlewareFunc()).Use(middleware2.AuthCheckRole())
	{
		r.GET("", api.GetSysChinaAreaDataList)
		r.GET("/:id", api.GetSysChinaAreaData)
		r.POST("", api.InsertSysChinaAreaData)
		r.PUT("/:id", api.UpdateSysChinaAreaData)
		r.DELETE("", api.DeleteSysChinaAreaData)
	}
}
