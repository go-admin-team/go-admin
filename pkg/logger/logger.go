package logger

import (
	"github.com/gin-gonic/gin"
	"go-admin/tools"
	"os"

	"github.com/go-admin-team/go-admin-core/logger"
)

// Logger 通用log个性化实现
type Logger struct {
	logger.Logger
}

// Info info级日志输出
func (l *Logger) Info(args ...interface{}) {
	l.Log(logger.InfoLevel, args...)
}

// Infof info级日志输出
func (l *Logger) Infof(template string, args ...interface{}) {
	l.Logf(logger.InfoLevel, template, args...)
}

// Trace trace级日志输出
func (l *Logger) Trace(args ...interface{}) {
	l.Log(logger.TraceLevel, args...)
}

// Tracef trace级日志输出
func (l *Logger) Tracef(template string, args ...interface{}) {
	l.Logf(logger.TraceLevel, template, args...)
}

// Debug debug级日志输出
func (l *Logger) Debug(args ...interface{}) {
	l.Log(logger.DebugLevel, args...)
}

// Debugf debug级日志输出
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Logf(logger.DebugLevel, template, args...)
}

// Warn warn级日志输出
func (l *Logger) Warn(args ...interface{}) {
	l.Log(logger.WarnLevel, args...)
}

// Warnf warn级日志输出
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Logf(logger.WarnLevel, template, args...)
}

// Error error级日志输出
func (l *Logger) Error(args ...interface{}) {
	l.Log(logger.ErrorLevel, args...)
}

// Errorf error级日志输出
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Logf(logger.ErrorLevel, template, args...)
}

// Fatal fatal级日志输出
func (l *Logger) Fatal(args ...interface{}) {
	l.Log(logger.FatalLevel, args...)
	os.Exit(1)
}

// Fatalf fatal级日志输出
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Logf(logger.FatalLevel, template, args...)
	os.Exit(1)
}

// GetRequestLogger 获取上下文提供的日志
func GetRequestLogger(c *gin.Context) Logger {
	if c == nil {
		return Logger{Logger: logger.DefaultLogger}
	}
	l, ok := c.Get(tools.LoggerKey)
	if !ok {
		return Logger{Logger: logger.DefaultLogger}
	}
	requestLogger, ok := l.(Logger)
	if !ok {
		return Logger{Logger: logger.DefaultLogger}
	}
	return requestLogger
}
