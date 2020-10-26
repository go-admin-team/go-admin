package config

import (
	"database/sql"
	"net/http"

	"github.com/go-admin-team/go-admin-core/logger"
)

type Config struct {
	saas   bool
	dbs    map[string]*DBConfig
	db     *DBConfig
	engine http.Handler
}

type DBConfig struct {
	Driver string
	DB     *sql.DB
}

// SetDbs 设置对应key的db
func (c *Config) SetDbs(key string, db *DBConfig) {
	c.dbs[key] = db
}

// GetDbs 获取所有map里的db数据
func (c *Config) GetDbs() map[string]*DBConfig {
	return c.dbs
}

// GetDbByKey 根据key获取db
func (c *Config) GetDbByKey(key string) *DBConfig {
	return c.dbs[key]
}

// SetDb 设置单个db
func (c *Config) SetDb(db *DBConfig) {
	c.db = db
}

// GetDb 获取单个db
func (c *Config) GetDb() *DBConfig {
	return c.db
}

// SetEngine 设置路由引擎
func (c *Config) SetEngine(engine http.Handler) {
	c.engine = engine
}

// GetEngine 获取路由引擎
func (c *Config) GetEngine() http.Handler {
	return c.engine
}

// SetLogger 设置日志组件
func (c *Config) SetLogger(l logger.Logger) {
	logger.DefaultLogger = l
}

// GetLogger 获取日志组件
func (c *Config) GetLogger() logger.Logger {
	return logger.DefaultLogger
}

// SetSaas 设置是否是saas应用
func (c *Config) SetSaas(saas bool) {
	c.saas = saas
}

// GetSaas 获取是否是saas应用
func (c *Config) GetSaas() bool {
	return c.saas
}

func DefaultConfig() *Config {
	return &Config{}
}
