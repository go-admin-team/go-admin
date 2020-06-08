package system

import (
	"github.com/gin-gonic/gin"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/captcha"
)

func GenerateCaptchaHandler(c *gin.Context) {
	id, b64s, err := captcha.DriverDigitFunc()
	tools.HasError(err, "验证码获取失败", 500)
	app.Custum(c, gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
