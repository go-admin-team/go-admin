package tools

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"gorm.io/gorm"

	config2 "github.com/go-admin-team/go-admin-core/sdk/config"
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

func (e *DBTables) GetPage(tx *gorm.DB, pageSize int, pageIndex int) ([]DBTables, int, error) {
	var doc []DBTables
	table := new(gorm.DB)
	var count int64

	if config2.DatabaseConfig.Driver == "mysql" {
		table = tx.Table("information_schema.tables")
		table = table.Where("TABLE_NAME not in (select table_name from `" + config2.GenConfig.DBName + "`.sys_tables) ")
		table = table.Where("table_schema= ? ", config2.GenConfig.DBName)

		if e.TableName != "" {
			table = table.Where("TABLE_NAME = ?", e.TableName)
		}
		if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
			return nil, 0, err
		}
	} else {
		pkg.Assert(true, "目前只支持mysql数据库", 500)
	}

	//table.Count(&count)
	return doc, int(count), nil
}

func (e *DBTables) Get(tx *gorm.DB) (DBTables, error) {
	var doc DBTables
	if config2.DatabaseConfig.Driver == "mysql" {
		table := tx.Table("information_schema.tables")
		table = table.Where("table_schema= ? ", config2.GenConfig.DBName)
		if e.TableName == "" {
			return doc, errors.New("table name cannot be empty！")
		}
		table = table.Where("TABLE_NAME = ?", e.TableName)
		if err := table.First(&doc).Error; err != nil {
			return doc, err
		}
	} else {
		pkg.Assert(true, "目前只支持mysql数据库", 500)
	}
	return doc, nil
}
