package router

import (
	log2 "go-admin/apis/log"
	"go-admin/apis/monitor"
	"go-admin/apis/public"
	"go-admin/apis/sysjob"
	"go-admin/apis/system"
	"go-admin/apis/system/dict"
	. "go-admin/apis/tools"
	_ "go-admin/docs"
	"go-admin/handler"
	"go-admin/middleware"
	jwt "go-admin/pkg/jwtauth"
	"go-admin/pkg/ws"
	"mime"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
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

	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	r.GET("/", system.HelloWorld)

	r.GET("/ws", ws.WebsocketManager.WsClient)


	r.GET("/info", handler.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	mime.AddExtensionType(".js", "application/javascript")
	
	r.Static("/static", "./static")
	r.Static("/form-generator", "./static/form-generator")
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
	v1.GET("/gen/todb/:tableId", GenMenuAndApi)
	v1.GET("/menuTreeselect", system.GetMenuTreeelect)
	v1.GET("/dict/databytype/:dictType", dict.GetDictDataByDictType)

	registerDBRouter(v1)

	registerSysTableRouter(v1)

	registerPublicRouter(v1)

	registerSysSettingRouter(v1)

	registerSysJobRouter(v1)
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

func registerSysJobRouter(v1 *gin.RouterGroup) {

	r := v1.Group("/sysjob")
	{
		r.GET("", sysjob.GetSysJobList)
		r.GET("/:jobId", sysjob.GetSysJob)
		r.POST("", sysjob.InsertSysJob)
		r.PUT("", sysjob.UpdateSysJob)
		r.DELETE("/:jobId", sysjob.DeleteSysJob)
	}

	v1.GET("/job/remove/:jobId",sysjob.RemoveJob)
	v1.GET("/job/start/:jobId",sysjob.StartJob)
}


func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r.POST("/login", authMiddleware.LoginHandler)
	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	v1 := r.Group("/api/v1")

	registerPageRouter(v1, authMiddleware)
	registerBaseRouter(v1, authMiddleware)
	registerDeptRouter(v1, authMiddleware)
	registerDictRouter(v1, authMiddleware)
	registerSysUserRouter(v1, authMiddleware)
	registerRoleRouter(v1, authMiddleware)
	registerConfigRouter(v1, authMiddleware)
	registerUserCenterRouter(v1, authMiddleware)
	registerPostRouter(v1, authMiddleware)
	registerMenuRouter(v1, authMiddleware)
	registerLoginLogRouter(v1, authMiddleware)
	registerOperLogRouter(v1, authMiddleware)
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/getinfo", system.GetInfo)
		v1auth.GET("/menurole", system.GetMenuRole)
		v1auth.PUT("/roledatascope", system.UpdateRoleDataScope)
		v1auth.GET("/roleMenuTreeselect/:roleId", system.GetMenuTreeRoleselect)
		v1auth.GET("/roleDeptTreeselect/:roleId", system.GetDeptTreeRoleselect)

		v1auth.POST("/logout", handler.LogOut)
		v1auth.GET("/menuids", system.GetMenuIDS)

		v1auth.GET("/operloglist", log2.GetOperLogList)
		v1auth.GET("/configKey/:configKey", system.GetConfigByConfigKey)
	}
}

func registerPageRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/deptList", system.GetDeptList)
		v1auth.GET("/deptTree", system.GetDeptTree)
		v1auth.GET("/sysUserList", system.GetSysUserList)
		v1auth.GET("/rolelist", system.GetRoleList)
		v1auth.GET("/configList", system.GetConfigList)
		v1auth.GET("/postlist", system.GetPostList)
		v1auth.GET("/menulist", system.GetMenuList)
		v1auth.GET("/loginloglist", log2.GetLoginLogList)
	}
}

func registerUserCenterRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	user := v1.Group("/user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		user.GET("/profile", system.GetSysUserProfile)
		user.POST("/avatar", system.InsetSysUserAvatar)
		user.PUT("/pwd", system.SysUserUpdatePwd)
	}
}

func registerOperLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	operlog := v1.Group("/operlog").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		operlog.GET("/:operId", log2.GetOperLog)
		operlog.DELETE("/:operId", log2.DeleteOperLog)
	}
}

func registerLoginLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	loginlog := v1.Group("/loginlog").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		loginlog.GET("/:infoId", log2.GetLoginLog)
		loginlog.POST("", log2.InsertLoginLog)
		loginlog.PUT("", log2.UpdateLoginLog)
		loginlog.DELETE("/:infoId", log2.DeleteLoginLog)
	}
}

func registerPostRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	post := v1.Group("/post").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		post.GET("/:postId", system.GetPost)
		post.POST("", system.InsertPost)
		post.PUT("", system.UpdatePost)
		post.DELETE("/:postId", system.DeletePost)
	}
}

func registerMenuRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	menu := v1.Group("/menu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		menu.GET("/:id", system.GetMenu)
		menu.POST("", system.InsertMenu)
		menu.PUT("", system.UpdateMenu)
		menu.DELETE("/:id", system.DeleteMenu)
	}
}

func registerConfigRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	config := v1.Group("/config").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		config.GET("/:configId", system.GetConfig)
		config.POST("", system.InsertConfig)
		config.PUT("", system.UpdateConfig)
		config.DELETE("/:configId", system.DeleteConfig)
	}
}

func registerRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	role := v1.Group("/role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		role.GET("/:roleId", system.GetRole)
		role.POST("", system.InsertRole)
		role.PUT("", system.UpdateRole)
		role.DELETE("/:roleId", system.DeleteRole)
	}
}

func registerSysUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	sysuser := v1.Group("/sysUser").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		sysuser.GET("/:userId", system.GetSysUser)
		sysuser.GET("/", system.GetSysUserInit)
		sysuser.POST("", system.InsertSysUser)
		sysuser.PUT("", system.UpdateSysUser)
		sysuser.DELETE("/:userId", system.DeleteSysUser)
	}
}

func registerDictRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	dicts := v1.Group("/dict").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		dicts.GET("/datalist", dict.GetDictDataList)
		dicts.GET("/typelist", dict.GetDictTypeList)
		dicts.GET("/typeoptionselect", dict.GetDictTypeOptionSelect)

		dicts.GET("/data/:dictCode", dict.GetDictData)
		dicts.POST("/data", dict.InsertDictData)
		dicts.PUT("/data/", dict.UpdateDictData)
		dicts.DELETE("/data/:dictCode", dict.DeleteDictData)

		dicts.GET("/type/:dictId", dict.GetDictType)
		dicts.POST("/type", dict.InsertDictType)
		dicts.PUT("/type", dict.UpdateDictType)
		dicts.DELETE("/type/:dictId", dict.DeleteDictType)
	}
}

func registerDeptRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	dept := v1.Group("/dept").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		dept.GET("/:deptId", system.GetDept)
		dept.POST("", system.InsertDept)
		dept.PUT("", system.UpdateDept)
		dept.DELETE("/:id", system.DeleteDept)
	}
}
func registerSysSettingRouter(v1 *gin.RouterGroup) {
	setting := v1.Group("/setting")
	{
		setting.GET("", system.GetSetting)
		setting.POST("", system.CreateSetting)
		setting.GET("/serverInfo",monitor.ServerInfo)
	}
}

func registerPublicRouter(v1 *gin.RouterGroup) {
	p := v1.Group("/public")
	{
		p.POST("/uploadFile", public.UploadFile)
	}
}
