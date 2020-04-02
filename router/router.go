package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	. "go-admin/apis"
	. "go-admin/apis/tools"
	_ "go-admin/docs"
	"go-admin/handler"
	"go-admin/handler/sd"
	_ "go-admin/pkg/jwtauth"
	"go-admin/router/middleware"
	"log"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerToFile())
	r.Use(middleware.CustomError)
	r.Use(middleware.NoCache)
	r.Use(middleware.Options)
	r.Use(middleware.Secure)
	r.Use(middleware.RequestId())
	r.Use(middleware.DemoEvn())

	r.Static("/static", "./static")
	r.GET("/info", Ping)

	// 监控信息
	svcd := r.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
		svcd.GET("/os", sd.OSCheck)
	}

	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()

	if err != nil {
		_ = fmt.Errorf("JWT Error", err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	//r.GET("/dashboard", Dashboard)
	r.GET("/routes", Dashboard)

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/getCaptcha", GenerateCaptchaHandler)
		apiv1.GET("/db/tables/page", GetDBTableList)
		apiv1.GET("/db/columns/page", GetDBColumnList)
		apiv1.GET("/sys/tables/page", GetSysTableList)
		apiv1.POST("/sys/tables/info", InsertSysTable)
		apiv1.PUT("/sys/tables/info", UpdateSysTable)
		apiv1.DELETE("/sys/tables/info/:tableId", DeleteSysTables)
		apiv1.GET("/sys/tables/info/:tableId", GetSysTables)
		apiv1.GET("/gen/preview/:tableId", Preview)
		apiv1.GET("/menuTreeselect", GetMenuTreeelect)
		apiv1.GET("/rolemenu", GetRoleMenu)
		apiv1.POST("/rolemenu", InsertRoleMenu)
		apiv1.DELETE("/rolemenu/:id", DeleteRoleMenu)
		apiv1.GET("/dict/databytype/:dictType", GetDictDataByDictType)


	}

	auth := r.Group("/api/v1")
	auth.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		auth.GET("/deptList", GetDeptList)
		auth.GET("/deptTree", GetDeptTree)
		auth.GET("/dept/:deptId", GetDept)
		auth.POST("/dept", InsertDept)
		auth.PUT("/dept", UpdateDept)
		auth.DELETE("/dept/:id", DeleteDept)

		auth.GET("/dict/datalist", GetDictDataList)
		auth.GET("/dict/data/:dictCode", GetDictData)
		auth.POST("/dict/data", InsertDictData)
		auth.PUT("/dict/data/", UpdateDictData)
		auth.DELETE("/dict/data/:dictCode", DeleteDictData)

		auth.GET("/dict/typelist", GetDictTypeList)
		auth.GET("/dict/type/:dictId", GetDictType)
		auth.POST("/dict/type", InsertDictType)
		auth.PUT("/dict/type", UpdateDictType)
		auth.DELETE("/dict/type/:dictId", DeleteDictType)

		auth.GET("/dict/typeoptionselect", GetDictTypeOptionSelect)

		auth.GET("/sysUserList", GetSysUserList)
		auth.GET("/sysUser/:userId", GetSysUser)
		auth.GET("/sysUser/", GetSysUserInit)
		auth.POST("/sysUser", InsertSysUser)
		auth.PUT("/sysUser", UpdateSysUser)
		auth.DELETE("/sysUser/:userId", DeleteSysUser)

		auth.GET("/rolelist", GetRoleList)
		auth.GET("/role/:roleId", GetRole)
		auth.POST("/role", InsertRole)
		auth.PUT("/role", UpdateRole)
		auth.PUT("/roledatascope", UpdateRoleDataScope)
		auth.DELETE("/role/:roleId", DeleteRole)

		auth.GET("/configList", GetConfigList)
		auth.GET("/config/:configId", GetConfig)
		auth.POST("/config", InsertConfig)
		auth.PUT("/config", UpdateConfig)
		auth.DELETE("/config/:configId", DeleteConfig)

		auth.GET("/roleMenuTreeselect/:roleId", GetMenuTreeRoleselect)
		auth.GET("/roleDeptTreeselect/:roleId", GetDeptTreeRoleselect)

		auth.GET("/getinfo", GetInfo)
		auth.GET("/user/profile", GetSysUserProfile)
		auth.POST("/user/avatar", InsetSysUserAvatar)
		auth.PUT("/user/pwd", SysUserUpdatePwd)

		auth.GET("/postlist", GetPostList)
		auth.GET("/post/:postId", GetPost)
		auth.POST("/post", InsertPost)
		auth.PUT("/post", UpdatePost)
		auth.DELETE("/post/:postId", DeletePost)

		auth.GET("/menulist", GetMenuList)
		auth.GET("/menu/:id", GetMenu)
		auth.POST("/menu", InsertMenu)
		auth.PUT("/menu", UpdateMenu)
		auth.DELETE("/menu/:id", DeleteMenu)
		auth.GET("/menurole", GetMenuRole)

		auth.GET("/menuids", GetMenuIDS)

		auth.GET("/loginloglist", GetLoginLogList)
		auth.GET("/loginlog/:infoId", GetLoginLog)
		auth.POST("/loginlog", InsertLoginLog)
		auth.PUT("/loginlog", UpdateLoginLog)
		auth.DELETE("/loginlog/:infoId", DeleteLoginLog)

		auth.GET("/operloglist", GetOperLogList)
		auth.GET("/operlog/:operId", GetOperLog)
		auth.DELETE("/operlog/:operId", DeleteOperLog)

		auth.GET("/configKey/:configKey", GetConfigByConfigKey)

		auth.POST("/logout", handler.LogOut)
	}

	//r.NoRoute(authMiddleware.MiddlewareFunc(), NoFound)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func Ping(c *gin.Context) {
	log.Println(viper.GetString(`settings.database.type`))
	c.JSON(200, gin.H{
		"message": "ok",
	})
}


func Dashboard(c *gin.Context) {

	var user = make(map[string]interface{})
	user["login_name"] = "admin"
	user["user_id"] = 1
	user["user_name"] = "管理员"
	user["dept_id"] = 1

	var cmenuList = make(map[string]interface{})
	cmenuList["children"] = nil
	cmenuList["parent_id"] = 1
	cmenuList["title"] = "用户管理"
	cmenuList["name"] = "Sysuser"
	cmenuList["icon"] = "user"
	cmenuList["order_num"] = 1
	cmenuList["id"] = 4
	cmenuList["path"] = "sysuser"
	cmenuList["component"] = "sysuser/index"

	var lista = make([]interface{}, 1)
	lista[0] = cmenuList

	var menuList = make(map[string]interface{})
	menuList["children"] = lista
	menuList["parent_id"] = 1
	menuList["name"] = "Upms"
	menuList["title"] = "权限管理"
	menuList["icon"] = "example"
	menuList["order_num"] = 1
	menuList["id"] = 4
	menuList["path"] = "/upms"
	menuList["component"] = "Layout"

	var list = make([]interface{}, 1)
	list[0] = menuList
	var data = make(map[string]interface{})
	data["user"] = user
	data["menuList"] = list

	var r = make(map[string]interface{})
	r["code"] = 200
	r["data"] = data

	c.JSON(200, r)
}
