package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis"
	"go-admin/app/admin/apis/tools"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, sysNoCheckRoleRouter, registerDBRouter, registerSysTableRouter)
}

func sysNoCheckRoleRouter(v1 *gin.RouterGroup) {
	r := v1.Group("")
	{
		sys := apis.System{}
		r.GET("/getCaptcha", sys.GenerateCaptchaHandler)
		gen := tools.Gen{}
		r.GET("/gen/preview/:tableId", gen.Preview)
		r.GET("/gen/toproject/:tableId", gen.GenCode)
		r.GET("/gen/apitofile/:tableId", gen.GenApiToFile)
		r.GET("/gen/todb/:tableId", gen.GenMenuAndApi)
		sysTable := tools.SysTable{}
		r.GET("/gen/tabletree", sysTable.GetSysTablesTree)
	}
}

func registerDBRouter(v1 *gin.RouterGroup) {
	db := v1.Group("/db")
	{
		gen := tools.Gen{}
		db.GET("/tables/page", gen.GetDBTableList)
		db.GET("/columns/page", gen.GetDBColumnList)
	}
}

func registerSysTableRouter(v1 *gin.RouterGroup) {
	tables := v1.Group("/sys/tables")
	{
		sysTable := tools.SysTable{}
		tables.GET("/page", sysTable.GetSysTableList)
		tablesInfo := tables.Group("/info")
		{
			tablesInfo.POST("", sysTable.InsertSysTable)
			tablesInfo.PUT("", sysTable.UpdateSysTable)
			tablesInfo.DELETE("/:tableId", sysTable.DeleteSysTables)
			tablesInfo.GET("/:tableId", sysTable.GetSysTables)
			tablesInfo.GET("", sysTable.GetSysTablesInfo)
		}
	}
}
