package log

import "go-admin/logger"

// Trace trace级日志输出
func Trace(args ...interface{}) {
	logger.DefaultLogger.Log(logger.TraceLevel, args...)
}

// Trace trace级日志输出
func Tracef(format string, args ...interface{}) {
	logger.DefaultLogger.Logf(logger.TraceLevel, format, args...)
}

// Debug debug级日志输出
func Debug(args ...interface{}) {
	logger.DefaultLogger.Log(logger.TraceLevel, args...)
}

// Debugf debug级日志输出
func Debugf(format string, args ...interface{}) {
	logger.DefaultLogger.Logf(logger.TraceLevel, format, args...)
}

// Info info级日志输出
func Info(args ...interface{}) {
	logger.DefaultLogger.Log(logger.TraceLevel, args...)
}

// Infof info级日志输出
func Infof(format string, args ...interface{}) {
	logger.DefaultLogger.Logf(logger.TraceLevel, format, args...)
}

// Warn warn级日志输出
func Warn(args ...interface{}) {
	logger.DefaultLogger.Log(logger.TraceLevel, args...)
}

// Warnf warn级日志输出
func Warnf(format string, args ...interface{}) {
	logger.DefaultLogger.Logf(logger.TraceLevel, format, args...)
}

// Error error级日志输出
func Error(args ...interface{}) {
	logger.DefaultLogger.Log(logger.TraceLevel, args...)
}

// Errorf error级日志输出
func Errorf(format string, args ...interface{}) {
	logger.DefaultLogger.Logf(logger.TraceLevel, format, args...)
}

// Fatal fatal级日志输出
func Fatal(args ...interface{}) {
	logger.DefaultLogger.Log(logger.TraceLevel, args...)
}

// Fatalf fatal级日志输出
func Fatalf(format string, args ...interface{}) {
	logger.DefaultLogger.Logf(logger.TraceLevel, format, args...)
}
