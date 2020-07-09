package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	log2 "go-admin/apis/log"
	"go-admin/apis/monitor"
	"go-admin/apis/system"
	"go-admin/apis/system/dict"
	. "go-admin/apis/tools"
	_ "go-admin/docs"
	"go-admin/handler"
	"go-admin/middleware"
	"go-admin/pkg/jwtauth"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwtauth.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	sysBaseRouter(g)
	// 静态文件
	sysStaticFileRouter(g)

	// swagger；注意：生产环境可以注释掉
	sysSwaggerRouter(g)

	// 无需认证
	sysNoCheckRoleRouter(g)
	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)

	return g
}

func sysBaseRouter(r *gin.RouterGroup) {
	r.GET("/", system.HelloWorld)
	r.GET("/info", handler.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	r.Static("/static", "./static")
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func sysNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")

	v1.GET("/monitor/server", monitor.ServerInfo)
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)
	v1.GET("/gen/preview/:tableId", Preview)
	v1.GET("/gen/toproject/:tableId", GenCode)

	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
	v1.GET("/dict/databytype/:dictType", dict.GetDictDataByDictType)

	registerDBRouter(v1)

	registerSysTableRouter(v1)

}

func registerDBRouter(api *gin.RouterGroup) {
	db := api.Group("/db")
	{
		db.GET("/tables/page", GetDBTableList)
		db.GET("/columns/page", GetDBColumnList)
	}
}

func registerSysTableRouter(v1 *gin.RouterGroup) {
	systables := v1.Group("/sys/tables")
	{
		systables.GET("/page", GetSysTableList)
		tablesinfo := systables.Group("/info")
		{
			tablesinfo.POST("", InsertSysTable)
			tablesinfo.PUT("", UpdateSysTable)
			tablesinfo.DELETE("/:tableId", DeleteSysTables)
			tablesinfo.GET("/:tableId", GetSysTables)
		}
	}
}

func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwtauth.GinJWTMiddleware) {
	r.POST("/login", authMiddleware.LoginHandler)
	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	v1 := r.Group("/api/v1").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		registerPageRouter(v1)
		registerBaseRouter(v1)
		registerDeptRouter(v1)
		registerDictRouter(v1)
		registerSysUserRouter(v1)
		registerRoleRouter(v1)
		registerConfigRouter(v1)
		registerUserCenterRouter(v1)
		registerPostRouter(v1)
		registerMenuRouter(v1)
		registerLoginLogRouter(v1)
		registerOperLogRouter(v1)
	}
}

func registerBaseRouter(v1 gin.IRoutes) {

	v1.GET("/getinfo", system.GetInfo)
	v1.GET("/menurole", system.GetMenuRole)
	v1.PUT("/roledatascope", system.UpdateRoleDataScope)
	v1.GET("/roleMenuTreeselect/:roleId", system.GetMenuTreeRoleselect)
	v1.GET("/roleDeptTreeselect/:roleId", system.GetDeptTreeRoleselect)

	v1.POST("/logout", handler.LogOut)
	v1.GET("/menuids", system.GetMenuIDS)

	v1.GET("/operloglist", log2.GetOperLogList)
	v1.GET("/configKey/:configKey", system.GetConfigByConfigKey)

}

func registerPageRouter(v1 gin.IRoutes) {
	v1.GET("/deptList", system.GetDeptList)
	v1.GET("/deptTree", system.GetDeptTree)
	v1.GET("/sysUserList", system.GetSysUserList)
	v1.GET("/rolelist", system.GetRoleList)
	v1.GET("/configList", system.GetConfigList)
	v1.GET("/postlist", system.GetPostList)
	v1.GET("/menulist", system.GetMenuList)
	v1.GET("/loginloglist", log2.GetLoginLogList)

}

func registerUserCenterRouter(v1 gin.IRoutes) {
	v1.GET("/user/profile", system.GetSysUserProfile)
	v1.POST("/user/avatar", system.InsetSysUserAvatar)
	v1.PUT("/user/pwd", system.SysUserUpdatePwd)
}

func registerOperLogRouter(v1 gin.IRoutes) {

	v1.GET("/operlog/:operId", log2.GetOperLog)
	v1.DELETE("/operlog/:operId", log2.DeleteOperLog)

}

func registerLoginLogRouter(v1 gin.IRoutes) {

	v1.GET("/loginlog/:infoId", log2.GetLoginLog)
	v1.POST("/loginlog", log2.InsertLoginLog)
	v1.PUT("/loginlog", log2.UpdateLoginLog)
	v1.DELETE("/loginlog/:infoId", log2.DeleteLoginLog)

}

func registerPostRouter(v1 gin.IRoutes) {

	v1.GET("/post/:postId", system.GetPost)
	v1.POST("/post", system.InsertPost)
	v1.PUT("/post", system.UpdatePost)
	v1.DELETE("/post/:postId", system.DeletePost)

}

func registerMenuRouter(v1 gin.IRoutes) {

	v1.GET("/menu/:id", system.GetMenu)
	v1.POST("/menu", system.InsertMenu)
	v1.PUT("/menu", system.UpdateMenu)
	v1.DELETE("/menu/:id", system.DeleteMenu)

}

func registerConfigRouter(v1 gin.IRoutes) {
	v1.GET("/config/:configId", system.GetConfig)
	v1.POST("/config", system.InsertConfig)
	v1.PUT("/config", system.UpdateConfig)
	v1.DELETE("/config/:configId", system.DeleteConfig)

}

func registerRoleRouter(v1 gin.IRoutes) {

	v1.GET("/role/:roleId", system.GetRole)
	v1.POST("/role", system.InsertRole)
	v1.PUT("/role", system.UpdateRole)
	v1.DELETE("/role/:roleId", system.DeleteRole)

}

func registerSysUserRouter(v1 gin.IRoutes) {

	v1.GET("/sysUser/:userId", system.GetSysUser)
	v1.GET("/sysUser/", system.GetSysUserInit)
	v1.POST("/sysUser", system.InsertSysUser)
	v1.PUT("/sysUser", system.UpdateSysUser)
	v1.DELETE("/sysUser/:userId", system.DeleteSysUser)

}

func registerDictRouter(v1 gin.IRoutes) {

	v1.GET("/dict/datalist", dict.GetDictDataList)
	v1.GET("/dict/typelist", dict.GetDictTypeList)
	v1.GET("/dict/typeoptionselect", dict.GetDictTypeOptionSelect)

	v1.GET("/dict/data/:dictCode", dict.GetDictData)
	v1.POST("/dict/data", dict.InsertDictData)
	v1.PUT("/dict/data/", dict.UpdateDictData)
	v1.DELETE("/dict/data/:dictCode", dict.DeleteDictData)

	v1.GET("/dict/type/:dictId", dict.GetDictType)
	v1.POST("/dict/type", dict.InsertDictType)
	v1.PUT("/dict/type", dict.UpdateDictType)
	v1.DELETE("/dict/type/:dictId", dict.DeleteDictType)

}

func registerDeptRouter(v1 gin.IRoutes) {

	v1.GET("/dept/:deptId", system.GetDept)
	v1.POST("/dept", system.InsertDept)
	v1.PUT("/dept", system.UpdateDept)
	v1.DELETE("/dept/:id", system.DeleteDept)

}
