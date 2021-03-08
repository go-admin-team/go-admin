package logger

import (
	"io"
	"os"

	"github.com/go-admin-team/go-admin-core/debug/writer"
	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/plugins/logger/zap"

	log "github.com/go-admin-team/go-admin-core/logger"
	"go-admin/tools"
)

// SetupLogger 日志
func SetupLogger(logType, path, levelStr, outputType string) logger.Logger {
	var setLogger logger.Logger
	if !tools.PathExist(path) {
		err := tools.PathCreate(path)
		if err != nil {
			log.Fatalf("create dir error: %s", err.Error())
		}
	}
	var err error
	var output io.Writer
	switch outputType {
	case "file":
		output, err = writer.NewFileWriter(path, "log")
		if err != nil {
			log.Fatal("logger setup error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}
	var level logger.Level
	level, err = logger.GetLevel(levelStr)
	if err != nil {
		log.Fatalf("get logger level error, %s", err.Error())
	}

	switch logType {
	case "zap":
		setLogger, err = zap.NewLogger(logger.WithLevel(level), logger.WithOutput(output), zap.WithCallerSkip(2))
		if err != nil {
			log.Fatalf("new zap logger error, %s", err.Error())
		}
	//case "logrus":
	//	setLogger = logrus.NewLogger(logger.WithLevel(level), logger.WithOutput(output), logrus.ReportCaller())
	default:
		setLogger = logger.NewLogger(logger.WithLevel(level), logger.WithOutput(output))
	}
	return setLogger
}
