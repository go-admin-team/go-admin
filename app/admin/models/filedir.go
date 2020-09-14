package models

import (
	orm "go-admin/common/global"
	"go-admin/tools"
)

type SysFileDir struct {
	Id        int          `json:"id"`
	Label     string       `json:"label" gorm:"type:varchar(255);"`    // 名称
	PId       int          `json:"pId" gorm:"type:int(11);"`           // 父id
	Sort      int          `json:"sort" gorm:""`                       //排序
	Path      string       `json:"path" gorm:"size:255;"`              //
	CreateBy  string       `json:"createBy" gorm:"type:varchar(128);"` // 创建人
	UpdateBy  string       `json:"updateBy" gorm:"type:varchar(128);"` // 编辑人
	Children  []SysFileDir `json:"children" gorm:"-"`
	DataScope string       `json:"dataScope" gorm:"-"`
	Params    string       `json:"params"  gorm:"-"`
	BaseModel
}

func (SysFileDir) TableName() string {
	return "sys_file_dir"
}

// 创建SysFileDir
func (e *SysFileDir) Create() (SysFileDir, error) {
	var doc SysFileDir
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}

	path := "/" + tools.IntToString(e.Id)
	if int(e.PId) != 0 {
		var deptP SysFileDir
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

// 获取SysFileDir
func (e *SysFileDir) Get() (SysFileDir, error) {
	var doc SysFileDir
	table := orm.Eloquent.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysFileDir带分页
func (e *SysFileDir) GetPage() ([]SysFileDir, int, error) {
	var doc []SysFileDir

	table := orm.Eloquent.Table(e.TableName())

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	//dataPermission := new(DataPermission)
	//dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	//table, err := dataPermission.GetDataScope(e.TableName(), table)
	//if err != nil {
	//	return nil, 0, err
	//}
	var count int64

	if err := table.Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, int(count), nil
}

// 更新SysFileDir
func (e *SysFileDir) Update(id int) (update SysFileDir, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	path := "/" + tools.IntToString(e.Id)
	if int(e.Id) != 0 {
		var deptP SysFileDir
		orm.Eloquent.Table(e.TableName()).Where("id = ?", e.Id).First(&deptP)
		path = deptP.Path + path
	} else {
		path = "/0" + path
	}
	e.Path = path

	//if e.Path != "" && e.Path != update.Path {
	//	return update, errors.New("上级不允许修改！")
	//}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除SysFileDir
func (e *SysFileDir) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&SysFileDir{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysFileDir) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&SysFileDir{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *SysFileDir) SetSysFileDir() ([]SysFileDir, error) {
	list, _, err := e.GetPage()

	m := make([]SysFileDir, 0)
	for i := 0; i < len(list); i++ {
		if list[i].PId != 0 {
			continue
		}
		info := SysFileDirDigui(&list, list[i])

		m = append(m, info)
	}
	return m, err
}

func SysFileDirDigui(deptlist *[]SysFileDir, menu SysFileDir) SysFileDir {
	list := *deptlist

	min := make([]SysFileDir, 0)
	for j := 0; j < len(list); j++ {

		if menu.Id != list[j].PId {
			continue
		}
		mi := SysFileDir{}
		mi.Id = list[j].Id
		mi.PId = list[j].PId
		mi.Label = list[j].Label
		mi.Sort = list[j].Sort
		mi.CreatedAt = list[j].CreatedAt
		mi.UpdatedAt = list[j].UpdatedAt
		mi.Children = []SysFileDir{}
		ms := SysFileDirDigui(deptlist, mi)
		min = append(min, ms)

	}
	menu.Children = min
	return menu
}
