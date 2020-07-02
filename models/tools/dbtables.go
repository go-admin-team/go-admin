package tools

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-admin/global/orm"
	config2 "go-admin/tools/config"
)

type DBTables struct {
	TableName      string `gorm:"column:TABLE_NAME" json:"tableName"`
	Engine         string `gorm:"column:ENGINE" json:"engine"`
	TableRows      string `gorm:"column:TABLE_ROWS" json:"tableRows"`
	TableCollation string `gorm:"column:TABLE_COLLATION" json:"tableCollation"`
	CreateTime     string `gorm:"column:CREATE_TIME" json:"createTime"`
	UpdateTime     string `gorm:"column:UPDATE_TIME" json:"updateTime"`
	TableComment   string `gorm:"column:TABLE_COMMENT" json:"tableComment"`
}

func (e *DBTables) GetPage(pageSize int, pageIndex int) ([]DBTables, int, error) {
	var doc []DBTables
	table := new(gorm.DB)
	var count int

	if config2.DatabaseConfig.DbType == "mysql" {
		table = orm.Eloquent.Select("*").Table("information_schema.tables")
		table = table.Where("TABLE_NAME not in (select table_name from " + config2.DatabaseConfig.Mysql.DBName + ".sys_tables) ")
		table = table.Where("table_schema= ? ", config2.DatabaseConfig.Mysql.DBName)

		if e.TableName != "" {
			table = table.Where("TABLE_NAME = ?", e.TableName)
		}
		if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
			return nil, 0, err
		}
	} else if config2.DatabaseConfig.DbType == "sqlite3" {
		table = orm.Eloquent.Select("name as TABLE_NAME,name as ENGINE,name as TABLE_ROWS,name as TABLE_COLLATION,name as CREATE_TIME,name as UPDATE_TIME,name as TABLE_COMMENT ").Table("sqlite_master")
		table = table.Where("type = 'table' ")

		if e.TableName != "" {
			table = table.Where("name = ?", e.TableName)
		}
		if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
			return nil, 0, err
		}
	}

	table.Count(&count)
	return doc, count, nil
}

func (e *DBTables) Get() (DBTables, error) {
	var doc DBTables
	table := new(gorm.DB)
	if config2.DatabaseConfig.DbType == "mysql" {
		table = orm.Eloquent.Select("*").Table("information_schema.tables")
		table = table.Where("table_schema= ? ", config2.DatabaseConfig.Mysql.DBName)
		if e.TableName == "" {
			return doc, errors.New("table name cannot be empty！")
		}
		table = table.Where("TABLE_NAME = ?", e.TableName)
	} else if config2.DatabaseConfig.DbType == "sqlite3" {
		table = orm.Eloquent.Select("name as TABLE_NAME,name as ENGINE,name as TABLE_ROWS,name as TABLE_COLLATION,name as CREATE_TIME,name as UPDATE_TIME,name as TABLE_COMMENT ").Table("sqlite_master")
		table = table.Where("type = 'table' ")

		if e.TableName != "" {
			table = table.Where("name = ?", e.TableName)
		}
	}
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}
