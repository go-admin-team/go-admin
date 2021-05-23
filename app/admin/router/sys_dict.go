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

		dicts.GET("/data", dataApi.GetSysDictDataList)
		dicts.GET("/data/:dictCode", dataApi.GetSysDictData)
		dicts.POST("/data", dataApi.InsertSysDictData)
		dicts.PUT("/data/:dictCode", dataApi.UpdateSysDictData)
		dicts.DELETE("/data", dataApi.DeleteSysDictData)

		dicts.GET("/type-option-select", dictApi.GetSysDictTypeAll)
		dicts.GET("/type", dictApi.GetSysDictTypeList)
		dicts.GET("/type/:id", dictApi.GetSysDictType)
		dicts.POST("/type", dictApi.InsertSysDictType)
		dicts.PUT("/type/:id", dictApi.UpdateSysDictType)
		dicts.DELETE("/type", dictApi.DeleteSysDictType)
	}
	v1.Group("/dict-data").Use(authMiddleware.MiddlewareFunc()).GET("/option-select", dataApi.GetSysDictDataAll)
}
