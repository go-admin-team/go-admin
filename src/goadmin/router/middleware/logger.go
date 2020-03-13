package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goadmin/config"
	"goadmin/models"
	"goadmin/utils"
	"os"
	"time"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	//写入文件
	src, err := os.OpenFile(config.ApplicationConfig.LogPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof(" %3d  %13v  %15s  %s  %s\r\n",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
		logger.Infoln()
		fmt.Println(
			startTime.Format("\n2006-01-02 15:04:05.9999"),
			"[INFO]",
			reqMethod,
			reqUri,
			statusCode,
			latencyTime,
			reqUri,
			clientIP,
		)
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" {
			SysOperLog := models.SysOperLog{}
			SysOperLog.OperIp = clientIP
			SysOperLog.OperLocation = utils.GetLocation(clientIP)
			SysOperLog.Status = utils.IntToString(statusCode)
			SysOperLog.OperName = utils.GetUserName(c)
			SysOperLog.RequestMethod = c.Request.Method
			SysOperLog.OperUrl = reqUri
			SysOperLog.Method = reqMethod
			b, _ := c.Get("body")
			SysOperLog.OperParam, _ = utils.StructToJsonStr(b)
			SysOperLog.CreateBy = utils.GetCurrntTime()
			SysOperLog.OperTime = utils.GetCurrntTime()
			SysOperLog.LatencyTime = (latencyTime).String()
			SysOperLog.UserAgent = c.Request.UserAgent()
			_, _ = SysOperLog.Create()
		}
	}
}
