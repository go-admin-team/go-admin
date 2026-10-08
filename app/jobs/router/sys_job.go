package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/v2/jwtauth"
	"go-admin/app/jobs/apis"
	models2 "go-admin/app/jobs/models"
	dto2 "go-admin/app/jobs/service/dto"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysJobRouter)
}

// 需认证的路由代码
func registerSysJobRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/sysjob").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), actions.Index[models2.SysJob, dto2.SysJobSearch]())
		// The detail answers SysJobItem, whose entryId the edit form reads;
		// the model itself marshals it as entry_id.
		r.GET("/:id", actions.PermissionAction(), actions.ViewAs[models2.SysJob, dto2.SysJobById, dto2.SysJobItem]())
		r.POST("", actions.Create[models2.SysJob, dto2.SysJobControl]())
		r.PUT("", actions.PermissionAction(), actions.Update[models2.SysJob, dto2.SysJobControl]())
		r.DELETE("", actions.PermissionAction(), actions.Delete[models2.SysJob, dto2.SysJobById]())
	}
	sysJob := apis.SysJob{}

	r2 := v1.Group("/job").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r2.GET("/remove/:id", sysJob.RemoveJobForService)
		r2.GET("/start/:id", sysJob.StartJobForService)
	}
}
