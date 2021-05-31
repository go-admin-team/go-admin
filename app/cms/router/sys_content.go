package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/cms/models"
	"go-admin/app/cms/service/dto"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysContentRouter)
}

// 需认证的路由代码
func registerSysContentRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/syscontent").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		model := &models.SysContent{}
		r.GET("", actions.PermissionAction(), actions.IndexAction(model, new(dto.SysContentSearch), func() interface{} {
			list := make([]models.SysContent, 0)
			return &list
		}))
		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto.SysContentById), nil))
		r.POST("", actions.CreateAction(new(dto.SysContentControl)))
		r.PUT("/:id", actions.PermissionAction(), actions.UpdateAction(new(dto.SysContentControl)))
		r.DELETE("", actions.PermissionAction(), actions.DeleteAction(new(dto.SysContentById)))
	}
}
