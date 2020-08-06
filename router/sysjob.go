package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/apis/sysjob"
)

// 无需认证的路由代码
func registerSysJobRouter(v1 *gin.RouterGroup) {

	r := v1.Group("/sysjob")
	{
		r.GET("", sysjob.GetSysJobList)
		r.GET("/:jobId", sysjob.GetSysJob)
		r.POST("", sysjob.InsertSysJob)
		r.PUT("", sysjob.UpdateSysJob)
		r.DELETE("/:jobId", sysjob.DeleteSysJob)
	}

	v1.GET("/job/remove/:jobId",sysjob.RemoveJob)
	v1.GET("/job/start/:jobId",sysjob.StartJob)
}
