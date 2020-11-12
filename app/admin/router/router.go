package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"go-admin/pkg/jwtauth"
	jwt "go-admin/pkg/jwtauth"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
)

// 路由示例
func InitExamplesRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {

	// 无需认证的路由
	examplesNoCheckRoleRouter(r)
	// 需要认证的路由
	examplesCheckRoleRouter(r, authMiddleware)

	return r
}

// 无需认证的路由示例
func examplesNoCheckRoleRouter(r *gin.Engine) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1")
	// 空接口防止v1定义无使用报错
	v1.GET("/nilcheckrole", nil)

	for _, f := range routerNoCheckRole {
		f(v1)
	}

	// {{无需认证路由自动补充在此处请勿删除}}
	//registerSysFileInfoRouter(v1)
}

// 需要认证的路由示例
func examplesCheckRoleRouter(r *gin.Engine, authMiddleware *jwtauth.GinJWTMiddleware) {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/api/v1")
	// 空接口防止v1定义无使用报错
	v1.GET("/checkrole", nil)

	for _, f := range routerCheckRole {
		f(v1, authMiddleware)
	}

	// {{认证路由自动补充在此处请勿删除}}

}
