package database

import (
	"database/sql"
	. "log"
	"time"

	goAdminLogger "github.com/go-admin-team/go-admin-core/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/config"
	"go-admin/common/global"
	"go-admin/common/log"
	"go-admin/tools"
	toolsConfig "go-admin/tools/config"
)

// Mysql mysql配置结构体
type Mysql struct {
}

// Setup 配置步骤
func (e *Mysql) Setup() {
	global.Source = e.GetConnect()
	log.Info(tools.Green(global.Source))
	db, err := sql.Open("mysql", global.Source)
	if err != nil {
		log.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	}
	global.Cfg.SetDb(&config.DBConfig{
		Driver: "mysql",
		DB:     db,
	})
	global.Eloquent, err = e.Open(db, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(tools.Red(e.GetDriver()+" connect error :"), err)
	} else {
		log.Info(tools.Green(e.GetDriver() + " connect success !"))
	}

	if global.Eloquent.Error != nil {
		log.Fatal(tools.Red(" database error :"), global.Eloquent.Error)
	}

	if toolsConfig.LoggerConfig.EnabledDB {
		global.Eloquent.Logger = logger.New(
			New(goAdminLogger.DefaultLogger.Options().Out, "\r\n", LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					goAdminLogger.DefaultLogger.Options().Level.LevelForGorm()),
			})
	}
}

// Open 打开数据库连接
func (e *Mysql) Open(db *sql.DB, cfg *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{Conn: db}), cfg)
}

// GetConnect 获取数据库连接
func (e *Mysql) GetConnect() string {
	return toolsConfig.DatabaseConfig.Source
}

// GetDriver 获取连接
func (e *Mysql) GetDriver() string {
	return toolsConfig.DatabaseConfig.Driver
}
