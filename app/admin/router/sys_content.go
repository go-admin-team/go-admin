package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/sys_content"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)


func init()  {
	routerCheckRole = append(routerCheckRole, registerSysContentRouter)
}

// 需认证的路由代码
func registerSysContentRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/syscontent").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("/:id", sys_content.GetSysContent)
		r.POST("", sys_content.InsertSysContent)
		r.PUT("", sys_content.UpdateSysContent)
		r.DELETE("/:id", sys_content.DeleteSysContent)
	}

	l := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		l.GET("/syscontentList", sys_content.GetSysContentList)
	}

}
