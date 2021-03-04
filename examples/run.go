// +build examples

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/common/global"
	myCasbin "go-admin/pkg/casbin"
	"gorm.io/driver/mysql"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp/inmg?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	syncEnforce := myCasbin.Setup(db, "sys_")
	global.Cfg.SetDb("*", db)
	global.Cfg.SetCasbin("*", syncEnforce)

	e := gin.Default()
	global.Cfg.SetEngine(e)
	log.Fatal(e.Run(":8000"))
}
