package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-admin/app/admin/service/dto"
	"go-admin/common"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/api"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"

	"go-admin/common/global"
)

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		// 开始时间
		startTime := time.Now()
		// 处理请求
		//
		// The body is only read when it has a destination. operParam below is
		// the only consumer, and it is written when logger.enableddb is on -
		// off in the shipped configuration, where reading the body was a copy
		// of every request made and discarded.
		var body string
		if config.LoggerConfig.EnabledDB {
			body = readOperParam(c, log)
		}

		c.Next()
		url := c.Request.RequestURI
		if strings.Index(url, "logout") > -1 ||
			strings.Index(url, "login") > -1 {
			return
		}
		// 结束时间
		endTime := time.Now()
		if c.Request.Method == http.MethodOptions {
			return
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
		clientIP := common.GetClientIP(c)
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
		log.WithFields(logData).Info()
		defer func() {
			log.Fields(map[string]interface{}{})
		}()
		if c.Request.Method != "OPTIONS" && config.LoggerConfig.EnabledDB && statusCode != 404 {
			SetDBOperLog(c, clientIP, statusCode, reqUri, reqMethod, latencyTime, body, result, statusBus)
		}
	}
}

// operParamLimit caps what is copied out of a request body for the operation
// log. A file upload is a POST like any other and reaches this middleware
// before any handler, so without a limit the whole upload is held in memory to
// write a log row - a 16MB upload allocated about 67MB. The limit also keeps
// the value inside the column, which is TEXT.
const operParamLimit = 32 << 10

// readOperParam copies the start of the request body for the operation log and
// leaves the request readable by the handler.
//
// The body is not buffered whole: the handler reads the part copied here from
// memory and the rest straight from the connection, so what this holds is
// bounded by operParamLimit however large the request is.
func readOperParam(c *gin.Context, log *logger.Helper) string {
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete:
	default:
		return ""
	}
	if c.Request.Body == nil {
		return ""
	}

	rest := c.Request.Body
	head := make([]byte, operParamLimit)
	n, err := io.ReadFull(rest, head)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		log.Warnf("read body for the operation log: %s", err)
	}
	head = head[:n]

	c.Request.Body = readCloser{
		Reader: io.MultiReader(bytes.NewReader(head), rest),
		Closer: rest,
	}
	return string(head)
}

type readCloser struct {
	io.Reader
	io.Closer
}

// SetDBOperLog 写入操作日志表 fixme 该方法后续即将弃用
func SetDBOperLog(c *gin.Context, clientIP string, statusCode int, reqUri string, reqMethod string, latencyTime time.Duration, body string, result string, status int) {

	log := api.GetRequestLogger(c)
	l := make(map[string]interface{})
	l["_fullPath"] = c.FullPath()
	l["operUrl"] = reqUri
	l["operIp"] = clientIP
	l["operLocation"] = "" // pkg.GetLocation(clientIP, gaConfig.ExtConfig.AMap.Key)
	l["operName"] = user.GetUserName(c)
	l["requestMethod"] = reqMethod
	l["operParam"] = body
	l["operTime"] = time.Now()
	l["jsonResult"] = result
	l["latencyTime"] = latencyTime.String()
	l["statusCode"] = statusCode
	l["userAgent"] = c.Request.UserAgent()
	l["createBy"] = user.GetUserId(c)
	l["updateBy"] = user.GetUserId(c)
	if status == http.StatusOK {
		l["status"] = dto.OperaStatusEnabel
	} else {
		l["status"] = dto.OperaStatusDisable
	}
	q := sdk.Runtime.GetQueuePrefix(c.Request.Host)
	message, err := sdk.Runtime.GetStreamMessage("", global.OperateLog, l)
	if err != nil {
		log.Errorf("GetStreamMessage error, %s", err.Error())
		// 日志报错错误，不中断请求
	} else {
		err = q.Append(message)
		if err != nil {
			log.Errorf("Append message error, %s", err.Error())
		}
	}
}
