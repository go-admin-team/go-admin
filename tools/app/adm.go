package app

import (
	"github.com/go-admin-team/go-admin-core/tools/runtime"
)

const (
	// go-admin Version Info
	Version = "1.2.3"
)

var Runtime runtime.Runtime = runtime.NewConfig()

var (
	Source string
	Driver string
	DBName string
)
