package tools

import (
	"errors"
	orm "go-admin/database"
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

	table := orm.Eloquent.Select("*").Table("information_schema.tables")
	table = table.Where("TABLE_NAME not in (select table_name from "+config2.DatabaseConfig.Name+".sys_tables) ")
	table = table.Where("table_schema= ? ", config2.DatabaseConfig.Name)

	if e.TableName != "" {
		table = table.Where("TABLE_NAME = ?", e.TableName)
	}

	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

func (e *DBTables) Get() (DBTables, error) {
	var doc DBTables

	table := orm.Eloquent.Select("*").Table("information_schema.tables")
	table = table.Where("table_schema= ? ", config2.DatabaseConfig.Name)
	if e.TableName == "" {
		return doc, errors.New("table name cannot be empty！")
	}
	table = table.Where("TABLE_NAME = ?", e.TableName)

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}
