package logger

import (
	"github.com/go-admin-team/go-admin-core/logger"
	"os"
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
	l.Log(logger.InfoLevel, args...)
}

// Tracef trace级日志输出
func (l *Logger) Tracef(template string, args ...interface{}) {
	l.Logf(logger.InfoLevel, template, args...)
}

// Debug debug级日志输出
func (l *Logger) Debug(args ...interface{}) {
	l.Log(logger.InfoLevel, args...)
}

// Debugf debug级日志输出
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Logf(logger.InfoLevel, template, args...)
}

// Warn warn级日志输出
func (l *Logger) Warn(args ...interface{}) {
	l.Log(logger.InfoLevel, args...)
}

// Warnf warn级日志输出
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Logf(logger.InfoLevel, template, args...)
}

// Error error级日志输出
func (l *Logger) Error(args ...interface{}) {
	l.Log(logger.InfoLevel, args...)
}

// Errorf error级日志输出
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Logf(logger.InfoLevel, template, args...)
}

// Fatal fatal级日志输出
func (l *Logger) Fatal(args ...interface{}) {
	l.Log(logger.InfoLevel, args...)
	os.Exit(1)
}

// Fatalf fatal级日志输出
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Logf(logger.InfoLevel, template, args...)
	os.Exit(1)
}
