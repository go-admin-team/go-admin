package logger

import (
	"go-admin/tools"
	"path/filepath"

	"github.com/go-admin-team/go-admin-core/debug/writer"
	"github.com/go-admin-team/go-admin-core/logger"

	"go-admin/common/log"
)

//func Setup() (*glog.Logger, *glog.Logger) {
//	var Logger *glog.Logger
//	var JobLogger *glog.Logger
//	var RequestLogger *glog.Logger
//
//	Logger = glog.New()
//	_ = Logger.SetPath(config.LoggerConfig.Path + "/bus")
//	Logger.SetStdoutPrint(config.LoggerConfig.EnabledBUS && config.LoggerConfig.Stdout)
//	Logger.SetFile("bus-{Ymd}.log")
//	_ = Logger.SetLevelStr(config.LoggerConfig.Level)
//
//	JobLogger = glog.New()
//	_ = JobLogger.SetPath(config.LoggerConfig.Path + "/job")
//	JobLogger.SetStdoutPrint(false)
//	JobLogger.SetFile("db-{Ymd}.log")
//	_ = JobLogger.SetLevelStr(config.LoggerConfig.Level)
//
//	RequestLogger = glog.New()
//	_ = RequestLogger.SetPath(config.LoggerConfig.Path + "/request")
//	RequestLogger.SetStdoutPrint(false)
//	RequestLogger.SetFile("access-{Ymd}.log")
//	_ = RequestLogger.SetLevelStr(config.LoggerConfig.Level)
//
//	Logger.Info(tools.Green("Logger init success!"))
//	return Logger, JobLogger
//}

// SetupLogger 日志
func SetupLogger(path string, subPath string) logger.Logger {
	var setLogger logger.Logger
	fullPath := filepath.Join(path, subPath)
	if !tools.PathExist(fullPath) {
		err := tools.PathCreate(fullPath)
		if err != nil {
			log.Fatal("create dir error: %s", err.Error())
		}
	}
	output, err := writer.NewFileWriter(fullPath, "log")
	if err != nil {
		log.Fatal("%s logger setup error: %s", subPath, err.Error())
	}
	setLogger = logger.NewHelper(logger.NewLogger(logger.WithOutput(output)))
	return setLogger
}
