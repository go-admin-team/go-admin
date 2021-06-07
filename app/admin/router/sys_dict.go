package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerDictRouter)
}

func registerDictRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	dictApi := apis.SysDictType{}
	dataApi := apis.SysDictData{}
	dicts := v1.Group("/dict").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{

		dicts.GET("/data", dataApi.GetList)
		dicts.GET("/data/:dictCode", dataApi.Get)
		dicts.POST("/data", dataApi.Insert)
		dicts.PUT("/data/:dictCode", dataApi.Update)
		dicts.DELETE("/data", dataApi.Delete)

		dicts.GET("/type-option-select", dictApi.GetSysDictTypeAll)
		dicts.GET("/type", dictApi.GetList)
		dicts.GET("/type/:id", dictApi.Get)
		dicts.POST("/type", dictApi.Insert)
		dicts.PUT("/type/:id", dictApi.Update)
		dicts.DELETE("/type", dictApi.Delete)
	}
	v1.Group("/dict-data").Use(authMiddleware.MiddlewareFunc()).GET("/option-select", dataApi.GetSysDictDataAll)
}
