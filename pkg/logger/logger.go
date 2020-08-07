package logger

import (
	"github.com/gogf/gf/os/glog"
	"go-admin/global"
	"go-admin/tools"
	"go-admin/tools/config"
)

var Logger *glog.Logger
var JobLogger *glog.Logger
var RequestLogger *glog.Logger

func Setup() {
	Logger = glog.New()
	_ = Logger.SetPath(config.LoggerConfig.Path + "/bus")
	Logger.SetStdoutPrint(config.LoggerConfig.EnabledBUS && config.LoggerConfig.Stdout)
	Logger.SetFile("bus-{Ymd}.log")
	_ = Logger.SetLevelStr(config.LoggerConfig.Level)

	JobLogger = glog.New()
	_ = JobLogger.SetPath(config.LoggerConfig.Path + "/job")
	JobLogger.SetStdoutPrint(false)
	JobLogger.SetFile("db-{Ymd}.log")
	_ = JobLogger.SetLevelStr(config.LoggerConfig.Level)

	RequestLogger = glog.New()
	_ = RequestLogger.SetPath(config.LoggerConfig.Path + "/request")
	RequestLogger.SetStdoutPrint(false)
	RequestLogger.SetFile("access-{Ymd}.log")
	_ = RequestLogger.SetLevelStr(config.LoggerConfig.Level)

	Logger.Info(tools.Green("Logger init success!"))

	global.Logger = Logger.Line()
	global.JobLogger = JobLogger.Line()
	global.RequestLogger = RequestLogger.Line()
}
