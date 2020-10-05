package migrate

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-admin/cmd/migrate/migration"
	_ "go-admin/cmd/migrate/migration/version"
	"go-admin/common/database"
	"go-admin/common/global"
	"go-admin/common/models"
	"go-admin/pkg/logger"
	"go-admin/tools/config"
)

var (
	configYml string
	exec      bool
	StartCmd  = &cobra.Command{
		Use:     "migrate",
		Short:   "Initialize the database",
		Example: "go-admin migrate -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	//StartCmd.PersistentFlags().BoolVarP(&exec, "exec", "e", false, "exec script")
}

func run() {
	usage := `start init`
	fmt.Println(usage)
	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	logger.Setup()
	//3. 初始化数据库链接
	database.Setup(config.DatabaseConfig.Driver)
	//4. 数据库迁移
	fmt.Println("数据库迁移开始")
	_ = migrateModel()
	//fmt.Println("数据库结构初始化成功！")
	//5. 数据初始化完成
	//if err := models.InitDb(); err != nil {
	//	global.Logger.Fatalf("数据库基础数据初始化失败！error: %v ", err)
	//}
	usage = `数据库基础数据初始化成功`
	fmt.Println(usage)
}

func migrateModel() error {
	if config.DatabaseConfig.Driver == "mysql" {
		global.Eloquent.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	err := global.Eloquent.Debug().AutoMigrate(&models.Migration{})
	if err != nil {
		return err
	}
	migration.Migrate.SetDb(global.Eloquent.Debug())
	migration.Migrate.Migrate()
	return err
}
