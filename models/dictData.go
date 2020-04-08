package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type DictData struct {
	//字典编码
	DictCode int64 `gorm:"column:dict_code;primary_key" json:"dictCode" example:"1"`

	//显示顺序
	DictSort int `gorm:"column:dict_sort" json:"dictSort" example:"1"`

	//数据标签
	DictLabel string `gorm:"column:dict_label" json:"dictLabel"`

	//数据键值
	DictValue string `gorm:"column:dict_value" json:"dictValue"`

	//字典类型
	DictType  string `gorm:"column:dict_type" json:"dictType"`
	CssClass  string `gorm:"column:css_class" json:"cssClass"`
	ListClass string `gorm:"column:list_class" json:"listClass"`
	IsDefault string `gorm:"column:is_default" json:"isDefault"`

	//状态
	Status     string `gorm:"column:status" json:"status"`
	Default    string `gorm:"column:default" json:"default"`
	CreateBy   string `gorm:"column:create_by" json:"createBy"`
	CreateTime string `gorm:"column:create_time" json:"createTime"`
	UpdateBy   string `gorm:"column:update_by" json:"updateBy"`
	UpdateTime string `gorm:"column:update_time" json:"updateTime"`

	//备注
	Remark    string `gorm:"column:remark" json:"remark"`
	Params    string `gorm:"column:params" json:"params"`
	DataScope string `gorm:"column:data_scope" json:"dataScope"`
	IsDel     int `gorm:"column:is_del" json:"isDel"`
}

func (e *DictData) Create() (DictData, error) {
	var doc DictData
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	result := orm.Eloquent.Table("sys_dict_data").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *DictData) GetByCode() (DictData, error) {
	var doc DictData

	table := orm.Eloquent.Table("sys_dict_data")
	if e.DictCode != 0 {
		table = table.Where("dict_code = ?", e.DictCode)
	}
	if e.DictLabel != "" {
		table = table.Where("dict_label = ?", e.DictLabel)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.Where("is_del = 0").First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictData) Get() ([]DictData, error) {
	var doc []DictData

	table := orm.Eloquent.Table("sys_dict_data")
	if e.DictCode != 0 {
		table = table.Where("dict_code = ?", e.DictCode)
	}
	if e.DictLabel != "" {
		table = table.Where("dict_label = ?", e.DictLabel)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.Where("is_del = 0").Order("dict_sort").Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictData) GetPage(pageSize int, pageIndex int) ([]DictData, int32, error) {
	var doc []DictData

	table := orm.Eloquent.Select("*").Table("sys_dict_data")
	if e.DictCode != 0 {
		table = table.Where("dict_code = ?", e.DictCode)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}
	if e.DictLabel != "" {
		table = table.Where("dict_label = ?", e.DictLabel)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
	table = dataPermission.GetDataScope("sys_dict_data", table)

	var count int32

	if err := table.Where("is_del = 0").Order("dict_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("is_del = 0").Count(&count)
	return doc, count, nil
}

func (e *DictData) Update(id int64) (update DictData, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_dict_data").Where("dict_code = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_dict_data").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *DictData) Delete(id int64) (success bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_dict_data").Where("dict_code = ?", id).Update(mp).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
