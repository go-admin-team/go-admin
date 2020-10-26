package mycasbin

import (
	"sync/atomic"

	"github.com/go-admin-team/go-admin-core/logger"
)

// Logger is the implementation for a Logger using golang log.
type Logger struct {
	enable int32
}

func (l *Logger) EnableLog(enable bool) {
	i := 0
	if enable {
		i = 1
	}
	atomic.StoreInt32(&(l.enable), int32(i))
}

func (l *Logger) IsEnabled() bool {
	return atomic.LoadInt32(&(l.enable)) != 0
}

func (l *Logger) Print(v ...interface{}) {
	if l.IsEnabled() {
		logger.DefaultLogger.Log(logger.InfoLevel, v...)
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.IsEnabled() {
		logger.DefaultLogger.Logf(logger.InfoLevel, format, v...)
	}
}
