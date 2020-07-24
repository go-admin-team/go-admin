package logger

import (
	"github.com/gogf/gf/os/glog"
	"go-admin/tools/config"
)

var Logger *glog.Logger

func InitLog() {
	Logger = glog.New()
	Logger.Path(config.LoggerConfig.Path)
	_ = Logger.SetLevelStr(config.LoggerConfig.Level)
	Logger.Info("Logger init Success!")
}
