package config

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/config"
	"github.com/go-admin-team/go-admin-core/config/source"
	"log"
)

var (
	ExtendConfig interface{}
	_watch       config.Watcher
	_cfg         *Settings
)

// Settings 兼容原先的配置结构
type Settings struct {
	Settings Config `yaml:"settings"`
}

// Config 配置集合
type Config struct {
	Application *Application          `yaml:"application"`
	Ssl         *Ssl                  `yaml:"ssl"`
	Logger      *Logger               `yaml:"logger"`
	Jwt         *Jwt                  `yaml:"jwt"`
	Database    *Database             `yaml:"database"`
	Databases   *map[string]*Database `yaml:"databases"`
	Gen         *Gen                  `yaml:"gen"`
	Extend      interface{}           `yaml:"extend"`
}

// 多db改造
func (e *Config) multiDatabase() {
	if len(*e.Databases) == 0 {
		*e.Databases = map[string]*Database{
			"*": e.Database,
		}

	}
}

// Setup 载入配置文件
func Setup(f func(opts ...source.Option) source.Source, options ...source.Option) {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("New config object fail: %s", err.Error()))
	}
	err = c.Load(
		f(options...),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("Load config source fail: %s", err.Error()))
	}

	_cfg = &Settings{Config{
		Application: ApplicationConfig,
		Ssl:         SslConfig,
		Logger:      LoggerConfig,
		Jwt:         JwtConfig,
		Database:    DatabaseConfig,
		Databases:   &DatabasesConfig,
		Gen:         GenConfig,
		Extend:      ExtendConfig,
	}}
	err = c.Scan(_cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Scan config fail: %s", err.Error()))
	}
	_cfg.Settings.multiDatabase()

	_watch, err = c.Watch()
	if err != nil {
		log.Fatal(fmt.Sprintf("Watch config source fail: %s", err.Error()))
	}
}

// Watch 配置监听, 重载时报错，不影响运行 fixme 数据连接 redis连接还没支持动态配置
func Watch() {
	for {
		v, err := _watch.Next()
		if err != nil {
			if err.Error() != "watch stopped" {
				log.Println(fmt.Sprintf("Next config fail: %s", err.Error()))
			}
			break
		}
		err = v.Scan(_cfg)
		if err != nil {
			log.Println(fmt.Sprintf("Scan config fail: %s", err.Error()))
			break
		}
		_cfg.Settings.multiDatabase()
	}
}

// Stop 停止监听
func Stop() {
	err := _watch.Stop()
	if err != nil {
		log.Println(err)
	}
}
