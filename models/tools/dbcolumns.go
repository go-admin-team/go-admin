package tools

import (
	"errors"
	"go-admin/global/orm"
	config2 "go-admin/tools/config"
	"unsafe"
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

type SqliteColumns struct {
	TableSchema            string `gorm:"column:-" json:"tableSchema"`
	TableName              string `gorm:"column:-" json:"tableName"`
	ColumnName             string `gorm:"column:name" json:"columnName"`
	ColumnDefault          string `gorm:"column:dflt_val" json:"columnDefault"`
	IsNullable             string `gorm:"column:notnull" json:"isNullable"`
	DataType               string `gorm:"column:-" json:"dataType"`
	CharacterMaximumLength string `gorm:"column:-" json:"characterMaximumLength"`
	CharacterSetName       string `gorm:"column:-" json:"characterSetName"`
	ColumnType             string `gorm:"column:type" json:"columnType"`
	ColumnKey              string `gorm:"column:-" json:"columnKey"`
	Extra                  string `gorm:"column:-" json:"extra"`
	ColumnComment          string `gorm:"column:-" json:"columnComment"`
}

func (e *DBColumns) GetPage(pageSize int, pageIndex int) ([]DBColumns, int, error) {
	var doc []DBColumns

	table := orm.Eloquent
	if config2.DatabaseConfig.Dbtype == "sqlite3" {
		if e.TableName == "" {
			return nil, 0, errors.New("table name cannot be empty！")
		}
		table = table.Raw("PRAGMA table_info(" + e.TableName + ");")
		var doc1 []SqliteColumns
		if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc1).Error; err != nil {
			return nil, 0, err
		}
		doc = *(*[]DBColumns)(unsafe.Pointer(&doc1))
	} else {
		table = table.Select("*").Table("information_schema.`COLUMNS`")
		table = table.Where("table_schema= ? ", config2.DatabaseConfig.Name)
		if e.TableName == "" {
			return nil, 0, errors.New("table name cannot be empty！")
		}
		table = table.Where("TABLE_NAME = ?", e.TableName)
		if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
			return nil, 0, err
		}
	}

	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

func (e *DBColumns) GetList() ([]DBColumns, error) {
	var doc []DBColumns

	table := orm.Eloquent
	if config2.DatabaseConfig.Dbtype == "sqlite3" {
		if e.TableName == "" {
			return nil, errors.New("table name cannot be empty！")
		}
		table = table.Raw("PRAGMA table_info(" + e.TableName + ");")
		var doc1 []SqliteColumns
		if err := table.Find(&doc1).Error; err != nil {
			return doc, err
		}
		doc = *(*[]DBColumns)(unsafe.Pointer(&doc1))
	} else {
		table = table.Select("*").Table("information_schema.columns")
		table = table.Where("table_schema= ? ", config2.DatabaseConfig.Name)
		if e.TableName == "" {
			return nil, errors.New("table name cannot be empty！")
		}
		table = table.Where("TABLE_NAME = ?", e.TableName)
		if err := table.Find(&doc).Error; err != nil {
			return doc, err
		}
	}
	return doc, nil
}
