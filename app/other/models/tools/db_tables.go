package tools

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"

	"gorm.io/gorm"

	config2 "github.com/go-admin-team/go-admin-core/v2/sdk/config"
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
	pkg.Assert(config2.DatabaseConfig.Driver == "mysql", "目前只支持mysql数据库", 500)

	var doc []DBTables
	var count int64

	// Tables already registered with the generator are not candidates. Read them
	// through the model on this connection: the subquery used to spell the
	// schema out by hand, so it only resolved when sys_tables happened to live
	// in the schema being generated from, and it counted soft-deleted rows.
	var generated []string
	if err := tx.Model(&SysTables{}).Pluck("table_name", &generated).Error; err != nil {
		return nil, 0, err
	}

	table := tx.Table("information_schema.tables")
	table = table.Where("table_schema= ? ", config2.GenConfig.DBName)
	if len(generated) > 0 {
		// NOT IN (NULL) is unknown for every row, so an empty list has to skip
		// the clause instead of rendering it.
		table = table.Where("TABLE_NAME not in (?)", generated)
	}

	if e.TableName != "" {
		table = table.Where("TABLE_NAME = ?", e.TableName)
	}
	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return doc, int(count), nil
}

func (e *DBTables) Get(tx *gorm.DB) (DBTables, error) {
	pkg.Assert(config2.DatabaseConfig.Driver == "mysql", "目前只支持mysql数据库", 500)

	var doc DBTables
	if e.TableName == "" {
		return doc, errors.New("table name cannot be empty！")
	}
	table := tx.Table("information_schema.tables")
	table = table.Where("table_schema= ? ", config2.GenConfig.DBName)
	table = table.Where("TABLE_NAME = ?", e.TableName)
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}
