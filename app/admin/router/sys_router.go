package router

import (
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"go-admin/app/admin/apis"
	"mime"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/ws"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"go-admin/common/middleware"
	"go-admin/common/middleware/handler"
	_ "go-admin/docs"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	sysBaseRouter(g)
	// 静态文件
	sysStaticFileRouter(g)
	// swagger；注意：生产环境可以注释掉
	if config.ApplicationConfig.Mode != "prod" {
		sysSwaggerRouter(g)
	}
	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)
	return g
}

func sysBaseRouter(r *gin.RouterGroup) {

	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	if config.ApplicationConfig.Mode != "prod" {
		r.GET("/", apis.GoAdmin)
	}
	r.GET("/info", handler.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return
	}
	r.Static("/static", "./static")
	if config.ApplicationConfig.Mode != "prod" {
		r.Static("/form-generator", "./static/form-generator")
	}
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	wss:=r.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		wss.GET("/ws/:id/:channel", ws.WebsocketManager.WsClient)
		wss.GET("/wslogout/:id/:channel", ws.WebsocketManager.UnWsClient)
	}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		// Refresh time can be longer than token timeout
		v1.GET("/refresh_token", authMiddleware.RefreshHandler)
	}
	registerBaseRouter(v1, authMiddleware)
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysMenu{}
	api2 := apis.SysDept{}
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/roleMenuTreeselect/:roleId", api.GetMenuTreeSelect)
		//v1.GET("/menuTreeselect", api.GetMenuTreeSelect)
		v1auth.GET("/roleDeptTreeselect/:roleId", api2.GetDeptTreeRoleSelect)
		v1auth.POST("/logout", handler.LogOut)
	}
}
