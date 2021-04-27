package system

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	"go-admin/common/apis"
)

type System struct {
	apis.Api
}

func (e System) GenerateCaptchaHandler(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	id, b64s, err := captcha.DriverDigitFunc()
	if err != nil {
		log.Errorf("DriverDigitFunc error, %s", err.Error())
		e.Error(500, err, "验证码获取失败")
		return
	}
	e.Custom(gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
