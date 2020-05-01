package config

import "github.com/spf13/viper"

type Log struct {
	Dir string
}

func InitLog(cfg *viper.Viper) *Log {
	return &Log{
		Dir: cfg.GetString("dir"),
	}
}

var LogConfig = new(Log)
