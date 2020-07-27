package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/mssola/user_agent"
	"go-admin/global"
	"go-admin/models"
	jwt "go-admin/pkg/jwtauth"
	"go-admin/tools"
	"go-admin/tools/config"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(models.SysUser)
		r, _ := v["role"].(models.SysRole)
		return jwt.MapClaims{
			jwt.IdentityKey:  u.UserId,
			jwt.RoleIdKey:    r.RoleId,
			jwt.RoleKey:      r.RoleKey,
			jwt.NiceKey:      u.Username,
			jwt.DataScopeKey: r.DataScope,
			jwt.RoleNameKey:  r.RoleName,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["nice"],
		"RoleKey":     claims["rolekey"],
		"UserId":      claims["identity"],
		"RoleIds":     claims["roleid"],
		"DataScope":   claims["datascope"],
	}
}

// @Summary 登陆
// @Description 获取token
// @Description LoginHandler can be used by clients to get a jwt token.
// @Description Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// @Description Reply will be of the form {"token": "TOKEN"}.
// @Description dev mode：It should be noted that all fields cannot be empty, and a value of 0 can be passed in addition to the account password
// @Description 注意：开发模式：需要注意全部字段不能为空，账号密码外可以传入0值
// @Accept  application/json
// @Product application/json
// @Param account body models.Login  true "account"
// @Success 200 {string} string "{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": ".eyJleHAiOjE1NjUxNTMxNDgsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2NTE0OTU0OH0.-zvzHvbg0A" }"
// @Router /login [post]
func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals models.Login
	var status = "0"
	var msg = "登录成功"
	var username = ""

	if err := c.ShouldBind(&loginVals); err != nil {
		username = loginVals.Username
		msg = "数据解析失败"
		status = "1"
		LoginLogToDB(c, status, msg, username)

		return nil, jwt.ErrMissingLoginValues
	}
	if config.ApplicationConfig.Mode != "dev" {
		if !store.Verify(loginVals.UUID, loginVals.Code, true) {
			username = loginVals.Username
			msg = "验证码错误"
			status = "1"
			LoginLogToDB(c, status, msg, username)

			return nil, jwt.ErrInvalidVerificationode
		}
	}
	user, role, e := loginVals.GetUser()
	if e == nil {
		username = loginVals.Username
		LoginLogToDB(c, status, msg, username)

		return map[string]interface{}{"user": user, "role": role}, nil
	} else {
		msg = "登录失败"
		status = "1"
		LoginLogToDB(c, status, msg, username)
		global.RequestLogger.Println(e.Error())
	}
	return nil, jwt.ErrFailedAuthentication
}

// Write log to database
func LoginLogToDB(c *gin.Context, status string, msg string, username string) {
	if config.LoggerConfig.EnabledDB {
		var loginlog models.LoginLog
		ua := user_agent.New(c.Request.UserAgent())
		loginlog.Ipaddr = c.ClientIP()
		loginlog.Username = username
		location := tools.GetLocation(c.ClientIP())
		loginlog.LoginLocation = location
		loginlog.LoginTime = tools.GetCurrentTime()
		loginlog.Status = status
		loginlog.Remark = c.Request.UserAgent()
		browserName, browserVersion := ua.Browser()
		loginlog.Browser = browserName + " " + browserVersion
		loginlog.Os = ua.OS()
		loginlog.Msg = msg
		loginlog.Platform = ua.Platform()
		_, _ = loginlog.Create()
	}
}

// @Summary 退出登录
// @Description 获取token
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Success 200 {string} string "{"code": 200, "msg": "成功退出系统" }"
// @Router /logout [post]
// @Security Bearer
func LogOut(c *gin.Context) {
	var loginlog models.LoginLog
	ua := user_agent.New(c.Request.UserAgent())
	loginlog.Ipaddr = c.ClientIP()
	location := tools.GetLocation(c.ClientIP())
	loginlog.LoginLocation = location
	loginlog.LoginTime = tools.GetCurrentTime()
	loginlog.Status = "0"
	loginlog.Remark = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	loginlog.Browser = browserName + " " + browserVersion
	loginlog.Os = ua.OS()
	loginlog.Platform = ua.Platform()
	loginlog.Username = tools.GetUserName(c)
	loginlog.Msg = "退出成功"
	loginlog.Create()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})

}

func Authorizator(data interface{}, c *gin.Context) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(models.SysUser)
		r, _ := v["role"].(models.SysRole)
		c.Set("role", r.RoleName)
		c.Set("roleIds", r.RoleId)
		c.Set("userId", u.UserId)
		c.Set("userName", u.UserName)
		c.Set("dataScope", r.DataScope)

		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
