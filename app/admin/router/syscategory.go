package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/syscategory"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

// 需认证的路由代码
func registerSysCategoryRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/syscategory").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("/:id", syscategory.GetSysCategory)
		r.POST("", syscategory.InsertSysCategory)
		r.PUT("", syscategory.UpdateSysCategory)
		r.DELETE("/:id", syscategory.DeleteSysCategory)
	}

	l := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		l.GET("/syscategoryList", syscategory.GetSysCategoryList)
	}

}
