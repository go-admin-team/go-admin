package tools

import (
	common "go-admin/common/models"

	"gorm.io/gorm"
)

type SysColumns struct {
	ColumnId           int          `gorm:"primaryKey;autoIncrement" json:"columnId"`
	TableId            int          `gorm:"" json:"tableId"`
	ColumnName         string       `gorm:"size:128;" json:"columnName"`
	ColumnComment      string       `gorm:"column:column_comment;size:128;" json:"columnComment"`
	ColumnType         string       `gorm:"column:column_type;size:128;" json:"columnType"`
	GoType             string       `gorm:"column:go_type;size:128;" json:"goType"`
	GoField            string       `gorm:"column:go_field;size:128;" json:"goField"`
	JsonField          string       `gorm:"column:json_field;size:128;" json:"jsonField"`
	IsPk               string       `gorm:"column:is_pk;size:4;" json:"isPk"`
	IsIncrement        string       `gorm:"column:is_increment;size:4;" json:"isIncrement"`
	IsRequired         string       `gorm:"column:is_required;size:4;" json:"isRequired"`
	IsInsert           string       `gorm:"column:is_insert;size:4;" json:"isInsert"`
	IsEdit             string       `gorm:"column:is_edit;size:4;" json:"isEdit"`
	IsList             string       `gorm:"column:is_list;size:4;" json:"isList"`
	IsQuery            string       `gorm:"column:is_query;size:4;" json:"isQuery"`
	QueryType          string       `gorm:"column:query_type;size:128;" json:"queryType"`
	HtmlType           string       `gorm:"column:html_type;size:128;" json:"htmlType"`
	DictType           string       `gorm:"column:dict_type;size:128;" json:"dictType"`
	Sort               int          `gorm:"column:sort;" json:"sort"`
	List               string       `gorm:"column:list;size:1;" json:"list"`
	Pk                 bool         `gorm:"column:pk;size:1;" json:"pk"`
	Required           bool         `gorm:"column:required;size:1;" json:"required"`
	SuperColumn        bool         `gorm:"column:super_column;size:1;" json:"superColumn"`
	UsableColumn       bool         `gorm:"column:usable_column;size:1;" json:"usableColumn"`
	Increment          bool         `gorm:"column:increment;size:1;" json:"increment"`
	Insert             bool         `gorm:"column:insert;size:1;" json:"insert"`
	Edit               bool         `gorm:"column:edit;size:1;" json:"edit"`
	Query              bool         `gorm:"column:query;size:1;" json:"query"`
	Remark             string       `gorm:"column:remark;size:255;" json:"remark"`
	FkTableName        string       `gorm:"" json:"fkTableName"`
	FkTableNameClass   string       `gorm:"" json:"fkTableNameClass"`
	FkTableNamePackage string       `gorm:"" json:"fkTableNamePackage"`
	FkCol              []SysColumns `gorm:"-" json:"fkCol"`
	FkLabelId          string       `gorm:"" json:"fkLabelId"`
	FkLabelName        string       `gorm:"size:255;" json:"fkLabelName"`
	CreateBy           int          `gorm:"column:create_by;size:20;" json:"createBy"`
	UpdateBy           int          `gorm:"column:update_By;size:20;" json:"updateBy"`

	// ColWidth and DefaultValue back PRD 010 F1/F2 (代码生成器前端模板迁移 Vue 3).
	// Both use a sentinel default (0 / "") rather than NULL - see
	// docs-prd/010-代码生成器前端模板迁移Vue3/数据库变更.md §1.1: a non-pointer
	// int/string field can never read NULL back out, and NULL would give
	// "unconfigured" two representations instead of one. Callers test
	// ColWidth == 0 / DefaultValue == "" to detect "not configured".
	//
	// ColWidth deliberately has no gorm size tag: this codebase's "size:N"
	// convention on numeric fields maps to a narrow SQL integer type (see
	// column_width_test.go), and col_width needs to hold values up to 800.
	ColWidth     int    `gorm:"column:col_width;not null;default:0;comment:table column width in px, 0 = not configured" json:"colWidth"`
	DefaultValue string `gorm:"column:default_value;size:255;not null;default:'';comment:form field default value, empty = not configured" json:"defaultValue"`

	common.ModelTime
}

func (*SysColumns) TableName() string {
	return "sys_columns"
}

func (e *SysColumns) GetList(tx *gorm.DB, exclude bool) ([]SysColumns, error) {
	var doc []SysColumns
	table := tx.Table("sys_columns")
	table = table.Where("table_id = ? ", e.TableId)
	if exclude {
		notIn := make([]string, 0, 6)
		notIn = append(notIn, "id")
		notIn = append(notIn, "create_by")
		notIn = append(notIn, "update_by")
		notIn = append(notIn, "created_at")
		notIn = append(notIn, "updated_at")
		notIn = append(notIn, "deleted_at")
		table = table.Where(" column_name not in(?)", notIn)
	}

	if err := table.Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *SysColumns) Create(tx *gorm.DB) (SysColumns, error) {
	var doc SysColumns
	e.CreateBy = 0
	result := tx.Table("sys_columns").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *SysColumns) Update(tx *gorm.DB) (update SysColumns, err error) {
	if err = tx.Table("sys_columns").First(&update, e.ColumnId).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	e.UpdateBy = 0
	if err = tx.Table("sys_columns").Model(&update).Updates(&e).Error; err != nil {
		return
	}

	// Updates(&e) above skips zero-value fields (GORM's struct-form Updates
	// always does), but ColWidth/DefaultValue's own "unconfigured" sentinel
	// is 0/"" (see the field comments on SysColumns) - so clearing either one
	// back to its sentinel is indistinguishable, to a struct-form Updates,
	// from "the caller didn't touch this field" and silently does not get
	// written. A map-form Updates does not skip zero values, so it is used
	// here for just these two columns rather than widening this to
	// Select("*") (which would also start writing every other zero-valued
	// field on this struct - Sort, the Pk/Required/... bools - and that is a
	// pre-existing gap in this method affecting fields outside PRD 010's
	// scope, not fixed here).
	if err = tx.Table("sys_columns").Model(&update).Updates(map[string]interface{}{
		"col_width":     e.ColWidth,
		"default_value": e.DefaultValue,
	}).Error; err != nil {
		return
	}

	return
}
