package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/sysfile"
)

// 无需认证的路由代码
func registerSysFileInfoRouter(v1 *gin.RouterGroup) {

	v1.GET("/sysfileinfoList", sysfile.GetSysFileInfoList)

	r := v1.Group("/sysfileinfo")
	{
		r.GET("/:id", sysfile.GetSysFileInfo)
		r.POST("", sysfile.InsertSysFileInfo)
		r.PUT("", sysfile.UpdateSysFileInfo)
		r.DELETE("/:id", sysfile.DeleteSysFileInfo)
	}
}

// 无需认证的路由代码
func registerSysFileDirRouter(v1 *gin.RouterGroup) {

	v1.GET("/sysfiledirList", sysfile.GetSysFileDirList)

	r := v1.Group("/sysfiledir")
	{
		r.GET("/:id", sysfile.GetSysFileDir)
		r.POST("", sysfile.InsertSysFileDir)
		r.PUT("", sysfile.UpdateSysFileDir)
		r.DELETE("/:id", sysfile.DeleteSysFileDir)
	}
}
