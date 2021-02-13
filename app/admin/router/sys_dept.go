package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/system/sys_dept"
	"go-admin/app/admin/middleware"
	jwt "go-admin/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysDeptRouter)
}

// 需认证的路由代码
func registerSysDeptRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := &sys_dept.SysDept{}
	r := v1.Group("/dept").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetSysDeptList)
		r.GET("/:id", api.GetSysDept)
		r.POST("", api.InsertSysDept)
		r.PUT("/:id", api.UpdateSysDept)
		r.DELETE("/:id", api.DeleteSysDept)
	}

	//r1 := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	//{
	//	r1.GET("/deptTree", api.GetDeptTree)
	//}

}
