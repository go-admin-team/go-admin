package log

import "go-admin/logger"

func SetLevel() {

}

// Trace trace级日志输出
func Trace(args ...interface{}) {
	logger.Trace(args...)
}

// Trace trace级日志输出
func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

// Debug debug级日志输出
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf debug级日志输出
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info info级日志输出
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof info级日志输出
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn warn级日志输出
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf warn级日志输出
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error error级日志输出
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf error级日志输出
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal fatal级日志输出
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf fatal级日志输出
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
