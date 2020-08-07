package models

import (
	"errors"
	orm "go-admin/global"
	"go-admin/tools"
	_ "time"
)

type SysDept struct {
	DeptId    int       `json:"deptId" gorm:"primary_key;AUTO_INCREMENT"` //部门编码
	ParentId  int       `json:"parentId" gorm:""`                         //上级部门
	DeptPath  string    `json:"deptPath" gorm:"size:255;"`                //
	DeptName  string    `json:"deptName"  gorm:"size:128;"`               //部门名称
	Sort      int       `json:"sort" gorm:""`                             //排序
	Leader    string    `json:"leader" gorm:"size:128;"`                  //负责人
	Phone     string    `json:"phone" gorm:"size:11;"`                    //手机
	Email     string    `json:"email" gorm:"size:64;"`                    //邮箱
	Status    string    `json:"status" gorm:"size:4;"`                    //状态
	CreateBy  string    `json:"createBy" gorm:"size:64;"`
	UpdateBy  string    `json:"updateBy" gorm:"size:64;"`
	DataScope string    `json:"dataScope" gorm:"-"`
	Params    string    `json:"params" gorm:"-"`
	Children  []SysDept `json:"children" gorm:"-"`
	BaseModel
}

func (SysDept) TableName() string {
	return "sys_dept"
}

type DeptLable struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLable `gorm:"-" json:"children"`
}

func (e *SysDept) Create() (SysDept, error) {
	var doc SysDept
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	deptPath := "/" + tools.IntToString(e.DeptId)
	if int(e.ParentId) != 0 {
		var deptP SysDept
		orm.Eloquent.Table(e.TableName()).Where("dept_id = ?", e.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	var mp = map[string]string{}
	mp["deptPath"] = deptPath
	if err := orm.Eloquent.Table(e.TableName()).Where("dept_id = ?", e.DeptId).Update(mp).Error; err != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	doc.DeptPath = deptPath
	return doc, nil
}

func (e *SysDept) Get() (SysDept, error) {
	var doc SysDept

	table := orm.Eloquent.Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysDept) GetList() ([]SysDept, error) {
	var doc []SysDept

	table := orm.Eloquent.Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Order("sort").Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysDept) GetPage(bl bool) ([]SysDept, error) {
	var doc []SysDept

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.DeptPath != "" {
		table = table.Where("deptPath like %?%", e.DeptPath)
	}
	if bl {
		// 数据权限控制
		dataPermission := new(DataPermission)
		dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
		tableper, err := dataPermission.GetDataScope("sys_dept", table)
		if err != nil {
			return nil, err
		}
		table = tableper
	}

	if err := table.Order("sort").Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *SysDept) SetDept(bl bool) ([]SysDept, error) {
	list, err := e.GetPage(bl)

	m := make([]SysDept, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := Digui(&list, list[i])

		m = append(m, info)
	}
	return m, err
}

func Digui(deptlist *[]SysDept, menu SysDept) SysDept {
	list := *deptlist

	min := make([]SysDept, 0)
	for j := 0; j < len(list); j++ {

		if menu.DeptId != list[j].ParentId {
			continue
		}
		mi := SysDept{}
		mi.DeptId = list[j].DeptId
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.DeptName = list[j].DeptName
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.Children = []SysDept{}
		ms := Digui(deptlist, mi)
		min = append(min, ms)

	}
	menu.Children = min
	return menu
}

func (e *SysDept) Update(id int) (update SysDept, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("dept_id = ?", id).First(&update).Error; err != nil {
		return
	}

	deptPath := "/" + tools.IntToString(e.DeptId)
	if int(e.ParentId) != 0 {
		var deptP SysDept
		orm.Eloquent.Table(e.TableName()).Where("dept_id = ?", e.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	e.DeptPath = deptPath

	if e.DeptPath != "" && e.DeptPath != update.DeptPath {
		return update, errors.New("上级部门不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}

	return
}

func (e *SysDept) Delete(id int) (success bool, err error) {

	user := SysUser{}
	user.DeptId = id
	userlist, err := user.GetList()
	tools.HasError(err, "", 500)
	tools.Assert(len(userlist) <= 0, "当前部门存在用户，不能删除！", 500)

	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		success = false
		return
	}

	if err = tx.Table(e.TableName()).Where("dept_id = ?", id).Delete(&SysDept{}).Error; err != nil {
		success = false
		tx.Rollback()
		return
	}
	if err = tx.Commit().Error; err != nil {
		success = false
		return
	}
	success = true

	return
}

func (dept *SysDept) SetDeptLable() (m []DeptLable, err error) {
	deptlist, err := dept.GetList()

	m = make([]DeptLable, 0)
	for i := 0; i < len(deptlist); i++ {
		if deptlist[i].ParentId != 0 {
			continue
		}
		e := DeptLable{}
		e.Id = deptlist[i].DeptId
		e.Label = deptlist[i].DeptName
		deptsInfo := DiguiDeptLable(&deptlist, e)

		m = append(m, deptsInfo)
	}
	return
}

func DiguiDeptLable(deptlist *[]SysDept, dept DeptLable) DeptLable {
	list := *deptlist

	min := make([]DeptLable, 0)
	for j := 0; j < len(list); j++ {

		if dept.Id != list[j].ParentId {
			continue
		}
		mi := DeptLable{list[j].DeptId, list[j].DeptName, []DeptLable{}}
		ms := DiguiDeptLable(deptlist, mi)
		min = append(min, ms)

	}
	dept.Children = min
	return dept
}
