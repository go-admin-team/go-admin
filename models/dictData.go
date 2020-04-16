package models

import (
	orm "go-admin/database"
	"go-admin/pkg/utils"
)

type DictData struct {
	DictCode  int    `gorm:"primary_key;AUTO_INCREMENT" json:"dictCode" example:"1"` //字典编码
	DictSort  int    `gorm:"type:int(4);" json:"dictSort"`                           //显示顺序
	DictLabel string `gorm:"type:varchar(128);" json:"dictLabel"`                    //数据标签
	DictValue string `gorm:"type:varchar(255);" json:"dictValue"`                    //数据键值
	DictType  string `gorm:"type:varchar(64);" json:"dictType"`                      //字典类型
	CssClass  string `gorm:"type:varchar(128);" json:"cssClass"`                     //
	ListClass string `gorm:"type:varchar(128);" json:"listClass"`                    //
	IsDefault string `gorm:"type:varchar(8);" json:"isDefault"`                      //
	Status    string `gorm:"type:int(1);" json:"status"`                             //状态
	Default   string `gorm:"type:varchar(8);" json:"default"`                        //
	CreateBy  string `gorm:"type:varchar(64);" json:"createBy"`                      //
	UpdateBy  string `gorm:"type:varchar(64);" json:"updateBy"`                      //
	Remark    string `gorm:"type:varchar(255);" json:"remark"`                       //备注
	Params    string `gorm:"-" json:"params"`
	DataScope string `gorm:"-" json:"dataScope"`
	BaseModel
}

func (DictData) TableName() string {
	return "sys_dict_data"
}

func (e *DictData) Create() (DictData, error) {
	var doc DictData
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *DictData) GetByCode() (DictData, error) {
	var doc DictData

	table := orm.Eloquent.Table(e.TableName())
	if e.DictCode != 0 {
		table = table.Where("dict_code = ?", e.DictCode)
	}
	if e.DictLabel != "" {
		table = table.Where("dict_label = ?", e.DictLabel)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictData) Get() ([]DictData, error) {
	var doc []DictData

	table := orm.Eloquent.Table(e.TableName())
	if e.DictCode != 0 {
		table = table.Where("dict_code = ?", e.DictCode)
	}
	if e.DictLabel != "" {
		table = table.Where("dict_label = ?", e.DictLabel)
	}
	if e.DictType != "" {
		table = table.Where("dict_type = ?", e.DictType)
	}

	if err := table.Order("dict_sort").Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *DictData) GetPage(pageSize int, pageIndex int) ([]DictData, int, error) {
	var doc []DictData

	table := orm.Eloquent.Select("*").Table(e.TableName())
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
	dataPermission.UserId, _ = utils.StringToInt(e.DataScope)
	table = dataPermission.GetDataScope("sys_dict_data", table)

	var count int

	if err := table.Order("dict_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

func (e *DictData) Update(id int) (update DictData, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("dict_code = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *DictData) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("dict_code = ?", id).Delete(&DictData{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
