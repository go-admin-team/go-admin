package migrate

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-admin/database"
	"go-admin/global/orm"
	"go-admin/models"
	"go-admin/models/gorm"
	"go-admin/tools"
	config2 "go-admin/tools/config"

	"github.com/spf13/cobra"
)

var (
	config   string
	mode     string
	StartCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize the database",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func run() {
	usage := `start init`
	fmt.Println(usage)
	//1. 读取配置
	config2.ConfigSetup(config)
	//2. 设置日志
	tools.InitLogger()
	//3. 初始化数据库链接
	database.Setup()
	//4. 数据库迁移
	_ = migrateModel()
	log.Println("数据库结构初始化成功！")
	//5. 数据初始化完成
	if err := models.InitDb(); err != nil {
		log.Fatal("数据库基础数据初始化失败！")
	}

	usage = `数据库基础数据初始化成功`
	fmt.Println(usage)
}

func migrateModel() error {
	if config2.DatabaseConfig.Dbtype == "mysql" {
		orm.Eloquent = orm.Eloquent.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	return gorm.AutoMigrate(orm.Eloquent)
}
