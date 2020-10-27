package config

import (
	"net/http"

	"github.com/go-admin-team/go-admin-core/logger"
)

type Conf interface {
	//多db设置，⚠️SetDbs不允许并发,可以根据自己的业务，例如app分库、host分库
	SetDbs(key string, db *DBConfig)
	GetDbs() map[string]*DBConfig
	GetDbByKey(key string) *DBConfig
	GetSaas() bool
	SetSaas(bool)

	//单库业务实现这两个接口
	SetDb(db *DBConfig)
	GetDb() *DBConfig

	//使用的路由
	SetEngine(engine http.Handler)
	GetEngine() http.Handler

	//使用go-admin定义的logger，参考来源go-micro
	SetLogger(logger logger.Logger)
	GetLogger() logger.Logger
}
