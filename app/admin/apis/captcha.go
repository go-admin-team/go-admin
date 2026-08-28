package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/captcha"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
)

type System struct {
	api.Api
}

// GenerateCaptchaHandler 获取验证码
// @Summary 获取验证码
// @Description 获取验证码
// @Tags 登陆
// @Success 200 {object} response.Response{data=string,id=string,msg=string} "{"code": 200, "data": [...]}"
// @Router /api/v1/captcha [get]
func (e System) GenerateCaptchaHandler(c *gin.Context) {
	if err := e.MakeContext(c).Errors; err != nil {
		e.Error(500, err, "服务初始化失败！")
		return
	}
	// The answer is deliberately discarded rather than logged. It used to be
	// written at info level, which put a currently valid captcha answer in the
	// application log - anyone able to read the log could bypass the check the
	// captcha exists to enforce.
	id, b64s, _, err := captcha.DriverDigitFunc()
	if err != nil {
		e.Logger.Errorf("DriverDigitFunc error, %s", err.Error())
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
