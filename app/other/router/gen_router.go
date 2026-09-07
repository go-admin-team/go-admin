package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"

	"go-admin/app/admin/apis"
	"go-admin/app/other/apis/tools"
)

// GenWriteRoutesEnabled reports whether the code generator's writing endpoints
// are registered in this process.
//
// Three of the generator's endpoints do not read. /gen/toproject writes seven
// Go and Vue source files onto this host, one of them under the path
// gen.frontpath names; /gen/apitofile writes a migration; /gen/todb inserts
// menus and APIs. All three are GET, all three are listed in CasbinExclude,
// and AuthCheckRole skips what is on that list - so Enforce never runs for
// them and any account that can log in may call them. That is a bargain a
// workstation can make and a deployment cannot.
//
// dev is the shipped default and is where the generator is meant to be used.
// demo keeps them because demo mode already has a better answer than a 404:
// DemoEvn refuses these three by name and explains itself, which is the thing
// the demo exists to show. test and prod get nothing.
//
// The mode is read once, while the routes are being built. Changing
// application.mode in a running process adds and removes nothing - a
// configuration reload rebuilds neither the engine nor its routes.
//
// core has constants for dev, test and prod but none for demo, which this
// repository spells as a literal in common/middleware/demo.go. Both are
// literals here so that the two read as one set rather than two conventions.
func GenWriteRoutesEnabled() bool {
	switch config.ApplicationConfig.Mode {
	case "dev", "demo":
		return true
	default:
		return false
	}
}

func init() {
	routerCheckRole = append(routerCheckRole, sysNoCheckRoleRouter, registerDBRouter, registerSysTableRouter)
}

func sysNoCheckRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r1 := v1.Group("")
	{
		sys := apis.System{}
		r1.GET("/captcha", sys.GenerateCaptchaHandler)
	}

	r := v1.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		gen := tools.Gen{}
		r.GET("/gen/preview/:tableId", gen.Preview)
		if GenWriteRoutesEnabled() {
			r.GET("/gen/toproject/:tableId", gen.GenCode)
			r.GET("/gen/apitofile/:tableId", gen.GenApiToFile)
			r.GET("/gen/todb/:tableId", gen.GenMenuAndApi)
		}
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
