package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"net/http"
)

//configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = base64Captcha.DefaultMemStore

func GenerateCaptchaHandler(c *gin.Context) {

	var param configJsonBody
	param.Id = uuid.New().String()
	param.CaptchaType = "string"
	param.DriverString = base64Captcha.NewDriverString(46, 140, 0, 0, 4, "1234567890abcdefghijklmnpqrstuvwxyz", &color.RGBA{45, 45, 45, 45}, []string{"wqy-microhei.ttc"})

	driver := param.DriverString.ConvertFonts()

	cap := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cap.Generate()
	body := map[string]interface{}{"code": 200, "data": b64s, "id": id, "msg": "success"}
	if err != nil {
		body = map[string]interface{}{"code": 0, "msg": err.Error()}
	}
	c.JSON(http.StatusOK, body)
}
