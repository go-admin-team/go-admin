package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/common/apis"
)

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		if c.Request.Method == http.MethodOptions {
			return
		}
		log := apis.GetRequestLogger(c)
		// 开始时间
		startTime := time.Now()
		// 处理请求

		bd, bl := c.Get("body")
		var body = ""
		if bl {
			body = bd.(string)
		}

		rt, bl := c.Get("result")
		var result = ""
		if bl {
			result = rt.(string)
		}

		st, bl := c.Get("status")
		var statusBus = 0
		if bl {
			statusBus = st.(int)
		}

		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 日志格式
		logData := map[string]interface{}{
			"statusCode":  statusCode,
			"latencyTime": latencyTime,
			"clientIP":    clientIP,
			"method":      reqMethod,
			"uri":         reqUri,
		}
		log.Info(logData)
		//l := logger.Logger{Logger: log.Fields(logData)}
		//l.Info(logData)
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" && config.LoggerConfig.EnabledDB {
			SetDBOperLog(c, clientIP, statusCode, reqUri, reqMethod, latencyTime, body, result, statusBus)
		}
	}
}

// SetDBOperLog 写入操作日志表 fixme 该方法后续即将弃用
func SetDBOperLog(c *gin.Context, clientIP string, statusCode int, reqUri string, reqMethod string, latencyTime time.Duration, body string, result string, status int) {
	log := apis.GetRequestLogger(c)
	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db connection error, %s", err.Error())
		app.Error(c, http.StatusInternalServerError, err, "数据库连接获取失败")
		return
	}

	menu := models.Menu{}
	menu.Path = reqUri
	menu.Action = reqMethod
	menuList, _ := menu.Get(db)
	sysOperaLog := system.SysOperaLog{}
	sysOperaLog.OperIp = clientIP
	sysOperaLog.OperLocation = pkg.GetLocation(clientIP)
	sysOperaLog.Status = pkg.IntToString(statusCode)
	sysOperaLog.OperName = user.GetUserName(c)
	sysOperaLog.RequestMethod = c.Request.Method
	sysOperaLog.OperUrl = reqUri
	sysOperaLog.OperParam = body
	if reqUri == "/login" {
		sysOperaLog.BusinessType = "10"
		sysOperaLog.Title = "用户登录"
		sysOperaLog.OperName = "-"
	} else if strings.Contains(reqUri, "/api/v1/logout") {
		sysOperaLog.BusinessType = "11"
	} else if strings.Contains(reqUri, "/api/v1/getCaptcha") {
		sysOperaLog.BusinessType = "12"
		sysOperaLog.Title = "验证码"
	} else {
		if reqMethod == "POST" {
			sysOperaLog.BusinessType = "1"
		} else if reqMethod == "PUT" {
			sysOperaLog.BusinessType = "2"
		} else if reqMethod == "DELETE" {
			sysOperaLog.BusinessType = "3"
		}
	}
	sysOperaLog.Method = reqMethod
	if len(menuList) > 0 {
		sysOperaLog.Title = menuList[0].Title
	}
	sysOperaLog.CreateBy = user.GetUserId(c)
	sysOperaLog.OperTime = pkg.GetCurrentTime()
	sysOperaLog.LatencyTime = fmt.Sprintf("%v", latencyTime)

	sysOperaLog.JsonResult = result
	sysOperaLog.UserAgent = c.Request.UserAgent()
	if status == 200 {
		sysOperaLog.Status = "2"
	} else {
		sysOperaLog.Status = "1"
	}
	serviceOperaLog := service.SysOperaLog{}
	serviceOperaLog.Orm = db
	serviceOperaLog.Log = log
	_ = serviceOperaLog.InsertSysOperaLog(&sysOperaLog)
}
