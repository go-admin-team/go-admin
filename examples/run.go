// +build examples

package main

import (
	"go-admin/tools/app"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	myCasbin "go-admin/pkg/casbin"
	"gorm.io/driver/mysql"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp/inmg?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	syncEnforce := myCasbin.Setup(db, "sys_")
	app.Runtime.SetDb("*", db)
	app.Runtime.SetCasbin("*", syncEnforce)

	e := gin.Default()
	app.Runtime.SetEngine(e)
	log.Fatal(e.Run(":8000"))
}
