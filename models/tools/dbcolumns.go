package tools

import (
	"errors"
	"github.com/jinzhu/gorm"
	orm "go-admin/global"
	"go-admin/tools"
	config2 "go-admin/tools/config"
)

type DBColumns struct {
	TableSchema            string `gorm:"column:TABLE_SCHEMA" json:"tableSchema"`
	TableName              string `gorm:"column:TABLE_NAME" json:"tableName"`
	ColumnName             string `gorm:"column:COLUMN_NAME" json:"columnName"`
	ColumnDefault          string `gorm:"column:COLUMN_DEFAULT" json:"columnDefault"`
	IsNullable             string `gorm:"column:IS_NULLABLE" json:"isNullable"`
	DataType               string `gorm:"column:DATA_TYPE" json:"dataType"`
	CharacterMaximumLength string `gorm:"column:CHARACTER_MAXIMUM_LENGTH" json:"characterMaximumLength"`
	CharacterSetName       string `gorm:"column:CHARACTER_SET_NAME" json:"characterSetName"`
	ColumnType             string `gorm:"column:COLUMN_TYPE" json:"columnType"`
	ColumnKey              string `gorm:"column:COLUMN_KEY" json:"columnKey"`
	Extra                  string `gorm:"column:EXTRA" json:"extra"`
	ColumnComment          string `gorm:"column:COLUMN_COMMENT" json:"columnComment"`
}

func (e *DBColumns) GetPage(pageSize int, pageIndex int) ([]DBColumns, int, error) {
	var doc []DBColumns
	var count int
	table := new(gorm.DB)

	if config2.DatabaseConfig.Driver == "mysql" {
		table = orm.Eloquent.Select("*").Table("information_schema.`COLUMNS`")
		table = table.Where("table_schema= ? ", config2.GenConfig.DBName)

		if e.TableName != "" {
			return nil, 0, errors.New("table name cannot be empty！")
		}

		table = table.Where("TABLE_NAME = ?", e.TableName)
	}

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil

}

func (e *DBColumns) GetList() ([]DBColumns, error) {
	var doc []DBColumns
	table := new(gorm.DB)

	if e.TableName == "" {
		return nil, errors.New("table name cannot be empty！")
	}

	if config2.DatabaseConfig.Driver == "mysql" {
		table = orm.Eloquent.Select("*").Table("information_schema.columns")
		table = table.Where("table_schema= ? ", config2.GenConfig.DBName)

		table = table.Where("TABLE_NAME = ?", e.TableName)
	} else {
		tools.Assert(true, "目前只支持mysql数据库", 500)
	}
	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}
