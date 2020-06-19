package tools

import (
	"go-admin/global/orm"
	"go-admin/models"
)

type SysColumns struct {
	ColumnId      int  `gorm:"primary_key;auto_increment" json:"columnId"`
	TableId       int  `gorm:"type:int(11);" json:"tableId"`
	ColumnName    string `gorm:"type:varchar(128);" json:"columnName"`
	ColumnComment string `gorm:"column:column_comment;type:varchar(128);" json:"columnComment"`
	ColumnType    string `gorm:"column:column_type;type:varchar(128);" json:"columnType"`
	GoType        string `gorm:"column:go_type;type:varchar(128);" json:"goType"`
	GoField       string `gorm:"column:go_field;type:varchar(128);" json:"goField"`
	JsonField     string `gorm:"column:json_field;type:varchar(128);" json:"jsonField"`
	IsPk          string `gorm:"column:is_pk;type:char(4);" json:"isPk"`
	IsIncrement   string `gorm:"column:is_increment;type:char(4);" json:"isIncrement"`
	IsRequired    string `gorm:"column:is_required;type:char(4);" json:"isRequired"`
	IsInsert      string `gorm:"column:is_insert;type:char(4);" json:"isInsert"`
	IsEdit        string `gorm:"column:is_edit;type:char(4);" json:"isEdit"`
	IsList        string `gorm:"column:is_list;type:char(4);" json:"isList"`
	IsQuery       string `gorm:"column:is_query;type:char(4);" json:"isQuery"`
	QueryType     string `gorm:"column:query_type;type:varchar(128);" json:"queryType"`
	HtmlType      string `gorm:"column:html_type;type:varchar(128);" json:"htmlType"`
	DictType      string `gorm:"column:dict_type;type:varchar(128);" json:"dictType"`
	Sort          int    `gorm:"column:sort;type:int(4);" json:"sort"`
	List          string `gorm:"column:list;type:char(1);" json:"list"`
	Pk            bool   `gorm:"column:pk;type:char(1);" json:"pk"`
	Required      bool   `gorm:"column:required;type:char(1);" json:"required"`
	SuperColumn   bool   `gorm:"column:super_column;type:char(1);" json:"superColumn"`
	UsableColumn  bool   `gorm:"column:usable_column;type:char(1);" json:"usableColumn"`
	Increment     bool   `gorm:"column:increment;type:char(1);" json:"increment"`
	Insert        bool   `gorm:"column:insert;type:char(1);" json:"insert"`
	Edit          bool   `gorm:"column:edit;type:char(1);" json:"edit"`
	Query         bool   `gorm:"column:query;type:char(1);" json:"query"`
	Remark        string `gorm:"column:remark;type:varchar(255);" json:"remark"`
	CreateBy      string `gorm:"column:create_by;type:varchar(128);" json:"createBy"`
	UpdateBy      string `gorm:"column:update_By;type:varchar(128);" json:"updateBy"`

	models.BaseModel
}

func (SysColumns) TableName() string {
	return "sys_columns"
}

func (e *SysColumns) GetList() ([]SysColumns, error) {
	var doc []SysColumns

	table := orm.Eloquent.Select("*").Table("sys_columns")

	table = table.Where("table_id = ?", e.TableId)

	if err := table.Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *SysColumns) Create() (SysColumns, error) {
	var doc SysColumns
	result := orm.Eloquent.Table("sys_columns").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *SysColumns) Update() (update SysColumns, err error) {
	if err = orm.Eloquent.Table("sys_columns").First(&update, e.ColumnId).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_columns").Model(&update).Updates(&e).Error; err != nil {
		return
	}

	return
}
