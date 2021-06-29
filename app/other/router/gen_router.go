package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/app/other/apis/tools"
)

func init() {
	routerCheckRole = append(routerCheckRole, sysNoCheckRoleRouter, registerDBRouter, registerSysTableRouter)
}

func sysNoCheckRoleRouter(v1 *gin.RouterGroup ,authMiddleware *jwt.GinJWTMiddleware) {
	r1 := v1.Group("")
	{
		sys := apis.System{}
		r1.GET("/captcha", sys.GenerateCaptchaHandler)
	}

	r := v1.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		gen := tools.Gen{}
		r.GET("/gen/preview/:tableId", gen.Preview)
		r.GET("/gen/toproject/:tableId", gen.GenCode)
		r.GET("/gen/apitofile/:tableId", gen.GenApiToFile)
		r.GET("/gen/todb/:tableId", gen.GenMenuAndApi)
		sysTable := tools.SysTable{}
		r.GET("/gen/tabletree", sysTable.GetSysTablesTree)
	}
}

func registerDBRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	db := v1.Group("/db").Use(authMiddleware.MiddlewareFunc())
	{
		gen := tools.Gen{}
		db.GET("/tables/page", gen.GetDBTableList)
		db.GET("/columns/page", gen.GetDBColumnList)
	}
}

func registerSysTableRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	tables := v1.Group("/sys/tables")
	{
		sysTable := tools.SysTable{}
		tables.Group("").Use(authMiddleware.MiddlewareFunc()).GET("/page", sysTable.GetPage)
		tablesInfo := tables.Group("/info").Use(authMiddleware.MiddlewareFunc())
		{
			tablesInfo.POST("", sysTable.Insert)
			tablesInfo.PUT("", sysTable.Update)
			tablesInfo.DELETE("/:tableId", sysTable.Delete)
			tablesInfo.GET("/:tableId", sysTable.Get)
			tablesInfo.GET("", sysTable.GetSysTablesInfo)
		}
	}
}