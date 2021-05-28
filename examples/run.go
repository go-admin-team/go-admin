// +build examples

package main

import (
	"github.com/go-admin-team/go-admin-core/sdk"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	myCasbin "github.com/go-admin-team/go-admin-core/sdk/pkg/casbin"
	"gorm.io/driver/mysql"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp/inmg?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	syncEnforce := myCasbin.Setup(db, "sys_")
	sdk.Runtime.SetDb("*", db)
	sdk.Runtime.SetCasbin("*", syncEnforce)

	e := gin.Default()
	sdk.Runtime.SetEngine(e)
	log.Fatal(e.Run(":8000"))
}
