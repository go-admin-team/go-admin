package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/apis/sysfiledir"
	"go-admin/app/admin/apis/sysfileinfo"
)

// 无需认证的路由代码
func registerSysFileInfoRouter(v1 *gin.RouterGroup) {

	v1.GET("/sysfileinfoList", sysfileinfo.GetSysFileInfoList)

	r := v1.Group("/sysfileinfo")
	{
		r.GET("/:id", sysfileinfo.GetSysFileInfo)
		r.POST("", sysfileinfo.InsertSysFileInfo)
		r.PUT("", sysfileinfo.UpdateSysFileInfo)
		r.DELETE("/:id", sysfileinfo.DeleteSysFileInfo)
	}
}

// 无需认证的路由代码
func registerSysFileDirRouter(v1 *gin.RouterGroup) {

	v1.GET("/sysfiledirList", sysfiledir.GetSysFileDirList)

	r := v1.Group("/sysfiledir")
	{
		r.GET("/:id", sysfiledir.GetSysFileDir)
		r.POST("", sysfiledir.InsertSysFileDir)
		r.PUT("", sysfiledir.UpdateSysFileDir)
		r.DELETE("/:id", sysfiledir.DeleteSysFileDir)
	}
}
