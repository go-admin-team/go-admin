package models

import (
	orm "go-admin/common/global"
	"go-admin/tools"
)

type SysCategory struct {
	Id        int    `json:"id" gorm:"type:int(11);primary_key;AUTO_INCREMENT"` // 分类Id
	Name      string `json:"name" gorm:"type:varchar(255);"`                    // 名称
	Img       string `json:"img" gorm:"type:varchar(255);"`                     // 图片
	Sort      string `json:"sort" gorm:"type:int(4);"`                          // 排序
	Status    string `json:"status" gorm:"type:int(1);"`                        // 状态
	Remark    string `json:"remark" gorm:"type:varchar(255);"`                  // 备注
	CreateBy  string `json:"createBy" gorm:"type:varchar(64);"`                 // 创建者
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(64);"`                 // 更新者
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
	BaseModel
}

func (SysCategory) TableName() string {
	return "sys_category"
}

// 创建SysCategory
func (e *SysCategory) Create() (SysCategory, error) {
	var doc SysCategory
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取SysCategory
func (e *SysCategory) Get() (SysCategory, error) {
	var doc SysCategory
	table := orm.Eloquent.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysCategory带分页
func (e *SysCategory) GetPage(pageSize int, pageIndex int) ([]SysCategory, int, error) {
	var doc []SysCategory

	table := orm.Eloquent.Table(e.TableName())

	if e.Name != "" {
		table = table.Where("name = ?", e.Name)
	}

	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int64

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	//table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, int(count), nil
}

// 更新SysCategory
func (e *SysCategory) Update(id int) (update SysCategory, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除SysCategory
func (e *SysCategory) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&SysCategory{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysCategory) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&SysCategory{}).Error; err != nil {
		return
	}
	Result = true
	return
}
