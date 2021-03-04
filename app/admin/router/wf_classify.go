package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/process"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerWfProcessClassifyRouter)
}

// 需认证的路由代码
func registerWfProcessClassifyRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &process.WfProcessClassify{}
	r := v1.Group("/process").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetWfProcessClassifyList)
		r.GET("/:id", api.GetWfProcessClassify)
		r.POST("", api.InsertWfProcessClassify)
		r.PUT("/:id", api.UpdateWfProcessClassify)
		r.DELETE("", api.DeleteWfProcessClassify)
	}
}
