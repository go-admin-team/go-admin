package models

import (
	orm "go-admin/common/global"
)

type SysFileInfoOld struct {
	Id        int    `json:"id"`                                 // id
	Type      string `json:"type" gorm:"type:varchar(255);"`     // 文件类型
	Name      string `json:"name" gorm:"type:varchar(255);"`     // 文件名称
	Size      string `json:"size" gorm:"type:int(11);"`          // 文件大小
	PId       int    `json:"pId" gorm:"type:int(11);"`           // 目录id
	Source    string `json:"source" gorm:"type:varchar(255);"`   // 文件源
	Url       string `json:"url" gorm:"type:varchar(255);"`      // 文件路径
	FullUrl   string `json:"fullUrl" gorm:"type:varchar(255);"`  // 文件全路径
	CreateBy  string `json:"createBy" gorm:"type:varchar(128);"` // 创建人
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(128);"` // 编辑人
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
	BaseModel
}

func (SysFileInfoOld) TableName() string {
	return "sys_file_info"
}

// 创建SysFileInfoOld
func (e *SysFileInfoOld) Create() (SysFileInfoOld, error) {
	var doc SysFileInfoOld
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取SysFileInfoOld
func (e *SysFileInfoOld) Get() (SysFileInfoOld, error) {
	var doc SysFileInfoOld
	table := orm.Eloquent.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysFileInfoOld带分页
func (e *SysFileInfoOld) GetPage(pageSize int, pageIndex int) ([]SysFileInfoOld, int, error) {
	var doc []SysFileInfoOld

	table := orm.Eloquent.Table(e.TableName())

	if e.PId != 0 {
		table = table.Where("p_id = ?", e.PId)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	//dataPermission := new(DataPermission)
	//dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	//table, err := dataPermission.GetDataScope(e.TableName(), table)
	//if err != nil {
	//	return nil, 0, err
	//}
	var count int64

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	//table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, int(count), nil
}

// 更新SysFileInfoOld
func (e *SysFileInfoOld) Update(id int) (update SysFileInfoOld, err error) {
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

// 删除SysFileInfoOld
func (e *SysFileInfoOld) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&SysFileInfoOld{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysFileInfoOld) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&SysFileInfoOld{}).Error; err != nil {
		return
	}
	Result = true
	return
}
