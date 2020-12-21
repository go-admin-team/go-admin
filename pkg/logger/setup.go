package logger

import (
	"github.com/go-admin-team/go-admin-core/debug/writer"
	"github.com/go-admin-team/go-admin-core/logger"

	"go-admin/common/log"
	"go-admin/tools"
)

// SetupLogger 日志
func SetupLogger(path string, prefix string) logger.Logger {
	var setLogger logger.Logger
	prefix = "[" + prefix + "] "
	if !tools.PathExist(path) {
		err := tools.PathCreate(path)
		if err != nil {
			log.Fatal("create dir error: %s", err.Error())
		}
	}
	output, err := writer.NewFileWriter(path, prefix, "log")
	if err != nil {
		log.Fatal("%s logger setup error: %s", prefix, err.Error())
	}
	setLogger = logger.NewHelper(logger.NewLogger(logger.WithOutput(output)))
	return setLogger
}
