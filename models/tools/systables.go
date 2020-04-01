package tools

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type SysTables struct {
	//表编码
	TableId int64 `gorm:"column:table_id;primary_key" json:"tableId"`
	//表名称
	TableName string `gorm:"column:table_name" json:"tableName"`
	//表备注
	TableComment string `gorm:"column:table_comment" json:"tableComment"`
	//类名
	ClassName string `gorm:"column:class_name" json:"className"`

	TplCategory string `gorm:"column:tpl_category" json:"tplCategory"`
	//包名
	PackageName string `gorm:"column:package_name" json:"packageName"`
	//模块名
	ModuleName   string `gorm:"column:module_name" json:"moduleName"`
	BusinessName string `gorm:"column:business_name" json:"businessName"`
	//功能名称
	FunctionName string `gorm:"column:function_name" json:"functionName"`
	//功能作者
	FunctionAuthor      string       `gorm:"column:function_author" json:"functionAuthor"`
	PkColumn            string       `gorm:"column:pk_column" json:"pkColumn"`
	PkGoField           string       `gorm:"column:pk_go_field" json:"pkGoField"`
	PkJsonField         string       `gorm:"column:pk_json_field" json:"pkJsonField"`
	Options             string       `gorm:"column:options" json:"options"`
	TreeCode            string       `gorm:"column:tree_code" json:"treeCode"`
	TreeParentCode      string       `gorm:"column:tree_parent_code" json:"treeParentCode"`
	TreeName            string       `gorm:"column:tree_name" json:"treeName"`
	Tree                bool         `gorm:"column:tree" json:"tree"`
	Crud                bool         `gorm:"column:crud" json:"crud"`
	Remark              string       `gorm:"column:remark" json:"remark"`
	IsLogicalDelete     string       `gorm:"column:is_logical_delete" json:"isLogicalDelete"`
	LogicalDelete       bool         `gorm:"column:logical_delete" json:"logicalDelete"`
	LogicalDeleteColumn string       `gorm:"column:logical_delete_column" json:"logicalDeleteColumn"`
	CreateBy            string       `gorm:"column:create_by" json:"createBy"`
	CreateTime          string       `gorm:"column:create_time" json:"createTime"`
	UpdateBy            string       `gorm:"column:update_By" json:"updateBy"`
	UpdateTime          string       `gorm:"column:update_time" json:"updateTime"`
	DataScope           string       `gorm:"-" json:"dataScope"`
	Params              Params       `gorm:"-" json:"params"`
	Columns             []SysColumns `gorm:"-" json:"columns"`
}

type Params struct {
	TreeCode       string `gorm:"-" json:"treeCode"`
	TreeParentCode string `gorm:"-" json:"treeParentCode"`
	TreeName       string `gorm:"-" json:"treeName"`
}

func (e *SysTables) GetPage(pageSize int, pageIndex int) ([]SysTables, int32, error) {
	var doc []SysTables

	table := orm.Eloquent.Select("*").Table("sys_tables")

	if e.TableName != "" {
		table = table.Where("table_name = ?", e.TableName)
	}
	if e.TableComment != "" {
		table = table.Where("table_comment = ?", e.TableComment)
	}

	var count int32

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

func (e *SysTables) Get() (SysTables, error) {
	var doc SysTables
	var err error
	table := orm.Eloquent.Select("*").Table("sys_tables")

	if e.TableName != "" {
		table = table.Where("table_name = ?", e.TableName)
	}
	if e.TableId != 0 {
		table = table.Where("table_id = ?", e.TableId)
	}
	if e.TableComment != "" {
		table = table.Where("table_comment = ?", e.TableComment)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	var col SysColumns
	col.TableId = doc.TableId
	if doc.Columns, err = col.GetList(); err != nil {
		return doc, err
	}

	return doc, nil
}

func (e *SysTables) Create() (SysTables, error) {
	var doc SysTables
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	result := orm.Eloquent.Table("sys_tables").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	for i := 0; i < len(e.Columns); i++ {
		e.Columns[i].TableId = doc.TableId

		e.Columns[i].Create()
	}

	return doc, nil
}

func (e *SysTables) Update() (update SysTables, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_tables").First(&update, e.TableId).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_tables").Model(&update).Updates(&e).Error; err != nil {
		return
	}

	for i := 0; i < len(e.Columns); i++ {
		_, _ = e.Columns[i].Update()
	}
	return
}

func (e *SysTables) Delete() (success bool, err error) {
	if err = orm.Eloquent.Table("sys_tables").Delete(SysTables{}, "table_id = ?", e.TableId).Error; err != nil {
		success = false
		return
	}
	if err = orm.Eloquent.Table("sys_columns").Delete(SysColumns{}, "table_id = ?", e.TableId).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
