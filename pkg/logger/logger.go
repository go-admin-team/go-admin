package logger

import (
	"github.com/gogf/gf/os/glog"
	"go-admin/global"
	"go-admin/tools"
	"go-admin/tools/config"
)

var Logger *glog.Logger
var DBLogger *glog.Logger
var RequestLogger *glog.Logger

func Setup() {
	Logger = glog.New()
	Logger.SetPath(config.LoggerConfig.Path + "/bus")
	Logger.SetStdoutPrint(config.LoggerConfig.EnabledBUS)
	Logger.SetFile("logger-{Ymd}.log")
	_ = Logger.SetLevelStr(config.LoggerConfig.Level)

	DBLogger = glog.New()
	_ = DBLogger.SetPath(config.LoggerConfig.Path + "/db")
	DBLogger.SetStdoutPrint(false)
	DBLogger.SetFile("db-{Ymd}.log")
	_ = DBLogger.SetLevelStr(config.LoggerConfig.Level)

	RequestLogger = glog.New()
	_ = RequestLogger.SetPath(config.LoggerConfig.Path + "/request")
	RequestLogger.SetStdoutPrint(false)
	RequestLogger.SetFile("access-{Ymd}.log")
	_ = RequestLogger.SetLevelStr(config.LoggerConfig.Level)


	Logger.Info(tools.Green("Logger init success!"))

	global.Logger = Logger.Line()
	global.DBLogger = DBLogger.Line()
	global.RequestLogger = RequestLogger.Line()
}
