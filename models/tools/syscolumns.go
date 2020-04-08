package tools

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type SysColumns struct {
	ColumnId      int64  `gorm:"column:column_id;primary_key" json:"columnId"`
	TableId       int64  `gorm:"column:table_id" json:"tableId"`
	ColumnName    string `gorm:"column:column_name" json:"columnName"`
	ColumnComment string `gorm:"column:column_comment" json:"columnComment"`
	ColumnType    string `gorm:"column:column_type" json:"columnType"`
	GoType        string `gorm:"column:go_type" json:"goType"`
	GoField       string `gorm:"column:go_field" json:"goField"`
	JsonField     string `gorm:"column:json_field" json:"jsonField"`
	IsPk          string `gorm:"column:is_pk" json:"isPk"`
	IsIncrement   string `gorm:"column:is_increment" json:"isIncrement"`
	IsRequired    string `gorm:"column:is_required" json:"isRequired"`
	IsInsert      string `gorm:"column:is_insert" json:"isInsert"`
	IsEdit        string `gorm:"column:is_edit" json:"isEdit"`
	IsList        string `gorm:"column:is_list" json:"isList"`
	IsQuery       string `gorm:"column:is_query" json:"isQuery"`
	QueryType     string `gorm:"column:query_type" json:"queryType"`
	HtmlType      string `gorm:"column:html_type" json:"htmlType"`
	DictType      string `gorm:"column:dict_type" json:"dictType"`
	Sort          int `gorm:"column:sort" json:"sort"`
	List          string `gorm:"column:list" json:"list"`
	Pk            bool   `gorm:"column:pk" json:"pk"`
	Required      bool   `gorm:"column:required" json:"required"`
	SuperColumn   bool   `gorm:"column:super_column" json:"superColumn"`
	UsableColumn  bool   `gorm:"column:usable_column" json:"usableColumn"`
	Increment     bool   `gorm:"column:increment" json:"increment"`
	Insert        bool   `gorm:"column:insert" json:"insert"`
	Edit          bool   `gorm:"column:edit" json:"edit"`
	Query         bool   `gorm:"column:query" json:"query"`
	Remark        string `gorm:"column:remark" json:"remark"`
	CreateBy      string `gorm:"column:create_by" json:"createBy"`
	CreateTime    string `gorm:"column:create_time" json:"createTime"`
	UpdateBy      string `gorm:"column:update_By" json:"updateBy"`
	UpdateTime    string `gorm:"column:update_time" json:"updateTime"`
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
	e.CreateTime = utils.GetCurrntTime()
	result := orm.Eloquent.Table("sys_columns").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *SysColumns) Update() (update SysColumns, err error) {
	e.UpdateTime = utils.GetCurrntTime()
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
