package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/scol"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerScolRouter)
}

// 需认证的路由代码
func registerScolRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &scol.Scol{}
	r := v1.Group("/scol").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetScolList)
		r.GET("/:id", api.GetScol)
		r.POST("", api.InsertScol)
		r.PUT("/:id", api.UpdateScol)
		r.DELETE("", api.DeleteScol)
	}
}
