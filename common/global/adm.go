package global

import (
	"go-admin/common/config"
)

const (
	// go-admin Version Info
	Version = "1.2.3"
)

var Cfg config.Conf = config.NewConfig()

var (
	Source string
	Driver string
	DBName string
)
