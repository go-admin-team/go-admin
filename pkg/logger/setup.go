package logger

import (
	"go-admin/tools"
	"path/filepath"

	"github.com/go-admin-team/go-admin-core/debug/writer"
	"github.com/go-admin-team/go-admin-core/logger"

	"go-admin/common/log"
)

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
