package migrate

import (
	"fmt"
	"go-admin/database"
	orm "go-admin/database"
	"go-admin/models"
	"go-admin/models/gorm"
	tools2 "go-admin/models/tools"
	"go-admin/tools"
	config2 "go-admin/tools/config"
	"log"

	"github.com/spf13/cobra"
)

var (
	config   string
	mode     string
	StartCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate from struct to database",
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
	usage := `start migrate`
	fmt.Println(usage)
	//1. 读取配置
	config2.ConfigSetup(config)
	//2. 设置日志
	tools.InitLogger()
	//3. 初始化数据库链接
	database.Setup()
	//4. 数据库迁移
	gorm.AutoMigrate(orm.Eloquent)
	log.Println("数据库结构初始化成功！")
	//5. 数据初始化完成
	if config2.ApplicationConfig.IsInit {
		if err := models.InitDb(); err != nil {
			log.Fatal("数据库基础数据初始化失败！")
		} else {
			config2.SetApplicationIsInit()
		}
		log.Println("数据库基础数据初始化成功！")
	}

	usage = `finish`
	fmt.Println(usage)
}

func migrateModel() error {
	if config2.DatabaseConfig.Dbtype == "mysql" {
		orm.Eloquent = orm.Eloquent.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	orm.Eloquent.SingularTable(true)
	return orm.Eloquent.AutoMigrate(
		new(models.CasbinRule),
		new(tools2.SysTables),
		new(tools2.SysColumns),
		new(models.Dept),
		new(models.Menu),
		new(models.LoginLog),
		new(models.SysOperLog),
		new(models.RoleMenu),
		new(models.SysRoleDept),
		new(models.SysUser),
		new(models.SysRole),
		new(models.Post),
		new(models.DictData),
		new(models.SysConfig),
		new(models.DictType),
	).Error
}
