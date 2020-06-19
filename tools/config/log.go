package config

import "github.com/spf13/viper"

type Log struct {
	Dir string
	Operdb bool
}

func InitLog(cfg *viper.Viper) *Log {
	return &Log{
		Dir: cfg.GetString("dir"),
		Operdb: cfg.GetBool("operdb"),
	}
}

var LogConfig = new(Log)
