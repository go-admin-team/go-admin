package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"go-admin/pkg/jwtauth"
	jwt "go-admin/pkg/jwtauth"
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

	v1 := r.Group("/api/v1")
	v1.GET("/nilcheckrole", nil)

	// {{无需认证路由自动补充在此处请勿删除}}
}

// 需要认证的路由示例
func examplesCheckRoleRouter(r *gin.Engine, authMiddleware *jwtauth.GinJWTMiddleware) {
	v1 := r.Group("/api/v1")
	v1.GET("/checkrole", nil)

	// {{认证路由自动补充在此处请勿删除}}
}
