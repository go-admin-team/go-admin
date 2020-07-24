package tools

import (
	"github.com/gogf/gf/os/glog"
	"go-admin/global"
	"go-admin/tools/config"
)

var Logger *glog.Logger
var DBLogger *glog.Logger
var AccessLogger *glog.Logger

func InitLogger() {
	Logger = glog.New()
	Logger.SetPath(config.LoggerConfig.Path)
	Logger.Stdout(config.LoggerConfig.Stdout)
	Logger.SetFile("logger-{Ymd}.log")
	_ = Logger.SetLevelStr(config.LoggerConfig.Level)

	Logger.Info("Logger init Success!")

	DBLogger = glog.New()
	_ = DBLogger.SetPath(config.DatabaseConfig.Logger.Path)
	DBLogger.Stdout(config.DatabaseConfig.Logger.Stdout)
	DBLogger.SetFile("db-{Ymd}.log")
	_ = DBLogger.SetLevelStr(config.DatabaseConfig.Logger.Level)
	DBLogger.Info("DBLogger init Success!")

	AccessLogger = glog.New()
	_ = AccessLogger.SetPath(config.ApplicationConfig.Logger.Path)
	AccessLogger.Stdout(config.ApplicationConfig.Logger.Stdout)
	AccessLogger.SetFile("access-{Ymd}.log")
	_ = AccessLogger.SetLevelStr(config.ApplicationConfig.Logger.Level)
	AccessLogger.Info("AccessLogger init Success!")

	global.Logger = Logger.Line()
	global.DBLogger = DBLogger.Line()
	global.AccessLogger = AccessLogger.Line()

}
