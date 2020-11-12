package models

import (
	orm "go-admin/common/global"
	"go-admin/tools"
)

type SysFileDirOld struct {
	Id        int          `json:"id"`
	Label     string       `json:"label" gorm:"type:varchar(255);"`    // 名称
	PId       int          `json:"pId" gorm:"type:int(11);"`           // 父id
	Sort      int          `json:"sort" gorm:""`                       //排序
	Path      string       `json:"path" gorm:"size:255;"`              //
	CreateBy  string       `json:"createBy" gorm:"type:varchar(128);"` // 创建人
	UpdateBy  string       `json:"updateBy" gorm:"type:varchar(128);"` // 编辑人
	Children  []SysFileDirOld `json:"children" gorm:"-"`
	DataScope string       `json:"dataScope" gorm:"-"`
	Params    string       `json:"params"  gorm:"-"`
	BaseModel
}

func (SysFileDirOld) TableName() string {
	return "sys_file_dir"
}

// 创建SysFileDirOld
func (e *SysFileDirOld) Create() (SysFileDirOld, error) {
	var doc SysFileDirOld
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}

	path := "/" + tools.IntToString(e.Id)
	if int(e.PId) != 0 {
		var deptP SysFileDirOld
		orm.Eloquent.Table(e.TableName()).Where("id = ?", e.PId).First(&deptP)
		path = deptP.Path + path
	} else {
		path = "/0" + path
	}
	var mp = map[string]string{}
	mp["path"] = path
	if err := orm.Eloquent.Table(e.TableName()).Where("id = ?", e.Id).Updates(mp).Error; err != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	doc.Path = path

	return doc, nil
}

// 获取SysFileDirOld
func (e *SysFileDirOld) Get() (SysFileDirOld, error) {
	var doc SysFileDirOld
	table := orm.Eloquent.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysFileDirOld带分页
func (e *SysFileDirOld) GetPage() ([]SysFileDirOld, int, error) {
	var doc []SysFileDirOld

	table := orm.Eloquent.Table(e.TableName())
	var count int64

	if err := table.Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, int(count), nil
}

// 更新SysFileDirOld
func (e *SysFileDirOld) Update(id int) (update SysFileDirOld, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	path := "/" + tools.IntToString(e.Id)
	if int(e.Id) != 0 {
		var deptP SysFileDirOld
		orm.Eloquent.Table(e.TableName()).Where("id = ?", e.Id).First(&deptP)
		path = deptP.Path + path
	} else {
		path = "/0" + path
	}
	e.Path = path

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除SysFileDirOld
func (e *SysFileDirOld) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&SysFileDirOld{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysFileDirOld) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&SysFileDirOld{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *SysFileDirOld) SetSysFileDirOld() ([]SysFileDirOld, error) {
	list, _, err := e.GetPage()

	m := make([]SysFileDirOld, 0)
	for i := 0; i < len(list); i++ {
		if list[i].PId != 0 {
			continue
		}
		info := SysFileDirOldDigui(&list, list[i])

		m = append(m, info)
	}
	return m, err
}

func SysFileDirOldDigui(deptlist *[]SysFileDirOld, menu SysFileDirOld) SysFileDirOld {
	list := *deptlist

	min := make([]SysFileDirOld, 0)
	for j := 0; j < len(list); j++ {

		if menu.Id != list[j].PId {
			continue
		}
		mi := SysFileDirOld{}
		mi.Id = list[j].Id
		mi.PId = list[j].PId
		mi.Label = list[j].Label
		mi.Sort = list[j].Sort
		mi.CreatedAt = list[j].CreatedAt
		mi.UpdatedAt = list[j].UpdatedAt
		mi.Children = []SysFileDirOld{}
		ms := SysFileDirOldDigui(deptlist, mi)
		min = append(min, ms)

	}
	menu.Children = min
	return menu
}
