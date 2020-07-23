package models

import (
	"errors"
	orm "go-admin/global"
	"go-admin/tools"
)

type DictType struct {
	DictId    int    `gorm:"primary_key;AUTO_INCREMENT" json:"dictId"`
	DictName  string `gorm:"size:128;" json:"dictName"` //字典名称
	DictType  string `gorm:"size:128;" json:"dictType"` //字典类型
	Status    string `gorm:"size:4;" json:"status"`     //状态
	DataScope string `gorm:"-" json:"dataScope"`        //
	Params    string `gorm:"-" json:"params"`           //
	CreateBy  string `gorm:"size:11;" json:"createBy"`  //创建者
	UpdateBy  string `gorm:"size:11;" json:"updateBy"`  //更新者
	Remark    string `gorm:"size:255;" json:"remark"`   //备注
	BaseModel
}

func (DictType) TableName() string {
	return "sys_dict_type"
}

func (e *DictType) Create() (DictType, error) {
	var doc DictType

	i := 0
	orm.Eloquent.Table(e.TableName()).Where("dict_name=? or dict_type = ?", e.DictName, e.DictType).Count(&i)
	if i > 0 {
		return doc, errors.New("字典名称或者字典类型已经存在！")
	}

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *DictType) Get() (DictType, error) {
	var doc DictType

	table := orm.Eloquent.Table(e.TableName())
	if e.DictId != 0 {
		table = table.Where("dict_id = ?", e.DictId)
	}
	if e.DictName != "" {
		table = table.Where("dict_name = ?", e.DictName)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictType) GetList() ([]DictType, error) {
	var doc []DictType

	table := orm.Eloquent.Table(e.TableName())
	if e.DictId != 0 {
		table = table.Where("dict_id = ?", e.DictId)
	}
	if e.DictName != "" {
		table = table.Where("dict_name = ?", e.DictName)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictType) GetPage(pageSize int, pageIndex int) ([]DictType, int, error) {
	var doc []DictType

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if e.DictId != 0 {
		table = table.Where("dict_id = ?", e.DictId)
	}
	if e.DictName != "" {
		table = table.Where("dict_name = ?", e.DictName)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("sys_dict_type", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

func (e *DictType) Update(id int) (update DictType, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if e.DictName != "" && e.DictName != update.DictName {
		return update, errors.New("名称不允许修改！")
	}

	if e.DictType != "" && e.DictType != update.DictType {
		return update, errors.New("类型不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *DictType) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("dict_id = ?", id).Delete(&DictData{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (e *DictType) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("dict_id in (?)", id).Delete(&DictType{}).Error; err != nil {
		return
	}
	Result = true
	return
}
