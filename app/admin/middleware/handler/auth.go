package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/mssola/user_agent"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/common/global"
	"go-admin/common/log"
	jwt "go-admin/pkg/jwtauth"
	"go-admin/tools"
	"go-admin/tools/config"
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
		global.RequestLogger.Error(e)
	}
	return nil, jwt.ErrFailedAuthentication
}

// LoginLogToDB Write log to database
func LoginLogToDB(c *gin.Context, status string, msg string, username string) {
	if config.LoggerConfig.EnabledDB {
		var loginLog system.SysLoginLog
		msgID := tools.GenerateMsgIDFromContext(c)
		db, err := tools.GetOrm(c)
		if err != nil {
			log.Errorf("msgID[%s] 获取Orm失败, error:%s", msgID, err)
		}
		ua := user_agent.New(c.Request.UserAgent())
		loginLog.Ipaddr = c.ClientIP()
		loginLog.Username = username
		location := tools.GetLocation(c.ClientIP())
		loginLog.LoginLocation = location
		loginLog.LoginTime = tools.GetCurrentTime()
		loginLog.Status = status
		loginLog.Remark = c.Request.UserAgent()
		browserName, browserVersion := ua.Browser()
		loginLog.Browser = browserName + " " + browserVersion
		loginLog.Os = ua.OS()
		loginLog.Msg = msg
		loginLog.Platform = ua.Platform()
		serviceLoginLog := service.SysLoginLog{}
		serviceLoginLog.Orm = db
		_ = serviceLoginLog.InsertSysLoginLog(loginLog.Generate())
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
	var loginLog system.SysLoginLog
	ua := user_agent.New(c.Request.UserAgent())
	loginLog.Ipaddr = c.ClientIP()
	location := tools.GetLocation(c.ClientIP())
	loginLog.LoginLocation = location
	loginLog.LoginTime = tools.GetCurrentTime()
	loginLog.Status = "0"
	loginLog.Remark = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	loginLog.Browser = browserName + " " + browserVersion
	loginLog.Os = ua.OS()
	loginLog.Platform = ua.Platform()
	loginLog.Username = tools.GetUserName(c)
	loginLog.Msg = "退出成功"
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] 获取Orm失败, error:%s", msgID, err)
	}
	serviceLoginLog := service.SysLoginLog{}
	serviceLoginLog.Orm = db
	_ = serviceLoginLog.InsertSysLoginLog(loginLog.Generate())

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
