package models

import (
	"fmt"
	orm "go-admin/database"
	config2 "go-admin/tools/config"
	"io/ioutil"
	"strings"
)

func InitDb() error {
	filePath := "config/db.sql"
	if config2.DatabaseConfig.Dbtype == "sqlite3" {
		fmt.Println("sqlite3数据库无需初始化！")
		return nil
	}
	sql, err := Ioutil(filePath)
	if err != nil {
		fmt.Println("数据库基础数据初始化脚本读取失败！原因:", err.Error())
		return err
	}
	sqlList := strings.Split(sql, ";")
	for i := 0; i < len(sqlList) - 1; i++ {
		if strings.Contains(sqlList[i], "--") {
			fmt.Println(sqlList[i])
			continue
		}
		sql := strings.Replace(sqlList[i]+";", "\n", "", 0)
		if err = orm.Eloquent.Exec(sql).Error; err != nil {
			if !strings.Contains(err.Error(), "Query was empty") {
				return err
			}
		}
	}
	return nil
}

func Ioutil(name string) (string, error) {
	if contents, err := ioutil.ReadFile(name); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		fmt.Println("Use ioutil.ReadFile to read a file:", result)
		return result, nil
	} else {
		return "", err
	}
}
