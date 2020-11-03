package log

import "github.com/go-admin-team/go-admin-core/logger"

var (
	// Trace trace级日志输出
	Trace = logger.Trace
	// Tracef trace级日志输出
	Tracef = logger.Tracef
	// Debug debug级日志输出
	Debug = logger.Debug
	// Debugf debug级日志输出
	Debugf = logger.Debugf
	// Info info级日志输出
	Info = logger.Info
	// Infof info级日志输出
	Infof = logger.Infof
	// Warn warn级日志输出
	Warn = logger.Warn
	// Warnf warn级日志输出
	Warnf = logger.Warnf
	// Error error级日志输出
	Error = logger.Error
	// Errorf error级日志输出
	Errorf = logger.Errorf
	// Fatal fatal级日志输出
	Fatal = logger.Fatal
	// Fatalf fatal级日志输出
	Fatalf = logger.Fatalf
)
