package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/other/apis"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerFileRouter)
}

// 需认证的路由代码
func registerFileRouter(v1 *gin.RouterGroup) {
	var api = apis.File{}
	v1.POST("/public/uploadFile", api.UploadFile)
}
