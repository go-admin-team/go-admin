package tools

import (
	"errors"
	orm "go-admin/database"
)

type DBColumns struct {
	//select COLUMN_NAME,COLUMN_DEFAULT,IS_NULLABLE,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,CHARACTER_SET_NAME,COLUMN_TYPE,COLUMN_KEY,EXTRA,COLUMN_COMMENT
	//from information_schema.`COLUMNS` where table_schema='test1db' and TABLE_NAME='sys_config'

	TableSchema string `gorm:"column:TABLE_SCHEMA" json:"tableSchema"`
	TableName string `gorm:"column:TABLE_NAME" json:"tableName"`
	Engine string `gorm:"column:COLUMN_NAME" json:"engine"`
	TableRows string `gorm:"column:COLUMN_DEFAULT" json:"tableRows"`
	TableCollation string `gorm:"column:IS_NULLABLE" json:"tableCollation"`
	CreateTime string `gorm:"column:DATA_TYPE" json:"createTime"`
	UpdateTime string `gorm:"column:CHARACTER_MAXIMUM_LENGTH" json:"updateTime"`
	TableComment string `gorm:"column:CHARACTER_SET_NAME" json:"tableComment"`
	ColumnType string `gorm:"column:COLUMN_TYPE" json:"columnType"`
	ColumnKey string `gorm:"column:COLUMN_KEY" json:"columnKey"`
	Extra string `gorm:"column:EXTRA" json:"extra"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT" json:"columnComment"`
}

func (e *DBColumns) GetPage(pageSize int, pageIndex int) ([]DBColumns, int32, error) {
	var doc []DBColumns

	table := orm.Eloquent.Select("*").Table("information_schema.`COLUMNS`")
	table=table.Where("table_schema='test1db'")

	if e.TableName != "" {
		return nil,0,errors.New("table name cannot be empty！")
	}

	table = table.Where("TABLE_NAME = ?", e.TableName)

	var count int32

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}
