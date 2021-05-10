package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysApiRouter)
}

// 需认证的路由代码
func registerSysApiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r := v1.Group("/sys-api").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		model := &models.SysApi{}
		r.GET("", actions.PermissionAction(), actions.IndexAction(model, new(dto.SysApiSearch), func() interface{} {
			list := make([]models.SysApi, 0)
			return &list
		}))
		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto.SysApiById), nil))
		r.POST("", actions.CreateAction(new(dto.SysApiControl)))
		r.PUT("/:id", actions.PermissionAction(), actions.UpdateAction(new(dto.SysApiControl)))
		r.DELETE("", actions.PermissionAction(), actions.DeleteAction(new(dto.SysApiById)))
	}
}
