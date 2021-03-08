package app

import (
	"github.com/go-admin-team/go-admin-core/tools/runtime"
)

const (
	// go-admin Version Info
	Version = "1.3.0-rc.0"
)

var Runtime runtime.Runtime = runtime.NewConfig()

var (
	Source string
	Driver string
	DBName string
)
