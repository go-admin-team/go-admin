package models

import (
	"fmt"
	"go-admin/common/global"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"strings"
)

func InitDb(db *gorm.DB) (err error) {
	filePath := "config/db.sql"
	if global.Driver == "postgres" {
		filePath = "config/pg.sql"
		err = ExecSql(db, filePath)
	} else if global.Driver == "mysql" {
		filePath = "config/db-begin-mysql.sql"
		err = ExecSql(db, filePath)
		filePath = "config/db.sql"
		err = ExecSql(db, filePath)
		filePath = "config/db-end-mysql.sql"
		err = ExecSql(db, filePath)
	} else {
		err = ExecSql(db, filePath)
	}
	return err
}

func ExecSql(db *gorm.DB, filePath string) error {
	sql, err := Ioutil(filePath)
	if err != nil {
		fmt.Println("数据库基础数据初始化脚本读取失败！原因:", err.Error())
		return err
	}
	sqlList := strings.Split(sql, ";")
	for i := 0; i < len(sqlList)-1; i++ {
		if strings.Contains(sqlList[i], "--") {
			fmt.Println(sqlList[i])
			continue
		}
		sql := strings.Replace(sqlList[i]+";", "\n", "", -1)
		sql = strings.TrimSpace(sql)
		if err = db.Exec(sql).Error; err != nil {
			log.Printf("error sql: %s", sql)
			if !strings.Contains(err.Error(), "Query was empty") {
				return err
			}
		}
	}
	return nil
}

func Ioutil(filePath string) (string, error) {
	if contents, err := ioutil.ReadFile(filePath); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		fmt.Println("Use ioutil.ReadFile to read a file:", result)
		return result, nil
	} else {
		return "", err
	}
}
