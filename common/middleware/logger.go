package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/common/global"
)

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求

		c.Next()
		// 结束时间
		endTime := time.Now()
		if c.Request.Method == http.MethodOptions {
			return
		}
		log := api.GetRequestLogger(c)

		bd, bl := c.Get("body")
		var body = ""
		if bl {
			body = bd.(string)
		}

		rt, bl := c.Get("result")
		var result = ""
		if bl {
			rb, err := json.Marshal(rt)
			if err != nil {
				log.Warnf("json Marshal result error, %s", err.Error())
			} else {
				result = string(rb)
			}
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
	log := api.GetRequestLogger(c)
	l := make(map[string]interface{})
	l["_fullPath"] = c.FullPath()
	l["operUrl"] = reqUri
	l["method"] = reqMethod
	l["operIp"] = clientIP
	l["operLocation"] = pkg.GetLocation(clientIP)
	l["operName"] = user.GetUserName(c)
	l["requestMethod"] = c.Request.Method
	l["operParam"] = body
	l["operTime"] = time.Now()
	if reqUri == "/login" {
		l["businessType"] = "10"
		l["title"] = "用户登录"
		l["operName"] = "-"
	} else if strings.Contains(reqUri, "/api/v1/logout") {
		l["businessType"] = "11"
		l["title"] = "退出登录"
	} else if strings.Contains(reqUri, "/api/v1/getCaptcha") {
		l["businessType"] = "12"
		l["title"] = "验证码"
	} else {
		if reqMethod == "POST" {
			l["businessType"] = "1"
		} else if reqMethod == "PUT" {
			l["businessType"] = "2"
		} else if reqMethod == "DELETE" {
			l["businessType"] = "3"
		}
	}
	if status == http.StatusOK {
		l["status"] = "2"
	} else {
		l["status"] = "1"
	}
	q := sdk.Runtime.GetCachePrefix(c.Request.Host)
	message, err := sdk.Runtime.GetStreamMessage("", global.OperateLog, l)
	if err != nil {
		log.Errorf("GetStreamMessage error, %s", err.Error())
		//日志报错错误，不中断请求
	} else {
		err = q.Append(message)
		if err != nil {
			log.Errorf("Append message error, %s", err.Error())
		}
	}
}
