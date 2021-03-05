package database

import (
	. "log"
	"time"

	logCore "github.com/go-admin-team/go-admin-core/logger"
	toolsDB "github.com/go-admin-team/go-admin-core/tools/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	log "github.com/go-admin-team/go-admin-core/logger"
	"go-admin/common/global"
	mycasbin "go-admin/pkg/casbin"
	"go-admin/tools"
	toolsConfig "go-admin/tools/config"
)

// Setup 配置数据库
func Setup() {
	for k := range toolsConfig.DatabasesConfig {
		setupSimpleDatabase(k, toolsConfig.DatabasesConfig[k])
	}
}

func setupSimpleDatabase(host string, c *toolsConfig.Database) {
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	log.Infof("%s => %s", host, tools.Green(c.Source))
	registers := make([]toolsDB.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = toolsDB.NewResolverConfigure(
			c.Registers[i].Sources,
			c.Registers[i].Replicas,
			c.Registers[i].Policy,
			c.Registers[i].Tables)
	}
	resolverConfig := toolsDB.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifetime, registers)
	db, err := resolverConfig.Init(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			New(logCore.DefaultLogger.Options().Out, "\r\n", LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					logCore.DefaultLogger.Options().Level.LevelForGorm()),
			},
		),
	}, opens[c.Driver])

	if err != nil {
		log.Fatal(tools.Red(c.Driver+" connect error :"), err)
	} else {
		log.Info(tools.Green(c.Driver + " connect success !"))
	}

	e := mycasbin.Setup(db, "sys_")

	//if host == "*" {
	//	global.Eloquent = db
	//}

	global.Runtime.SetDb(host, db)
	global.Runtime.SetCasbin(host, e)
}
