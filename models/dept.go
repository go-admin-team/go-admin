package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type Dept struct {
	//部门编码
	Deptid int64 `gorm:"column:dept_id;primary_key" json:"deptId" example:"1" extensions:"x-description=标示"`

	//上级部门
	ParentId int64 `gorm:"column:parent_id" json:"parent_id"`

	DeptPath string `gorm:"column:dept_path" json:"deptath"`

	//部门名称
	Deptname string `gorm:"column:dept_name" json:"deptname"`

	//排序
	Sort int64 `gorm:"column:sort" json:"sort"`

	//负责人
	Leader string `gorm:"column:leader" json:"leader"`

	//手机
	Phone string `gorm:"column:phone" json:"phone"`

	//邮箱
	Email string `gorm:"column:email" json:"email"`

	//状态
	Status string `gorm:"column:status" json:"status"`

	CreateBy   string `gorm:"column:create_by" json:"createBy"`
	CreateTime string `gorm:"column:create_time" json:"createTime"`
	UpdateBy   string `gorm:"column:update_by" json:"updateBy"`
	UpdateTime string `gorm:"column:update_time" json:"updateTime"`
	DataScope  string `gorm:"-" json:"dataScope"`
	Params     string `gorm:"-" json:"params"`
	IsDel      int    `gorm:"column:is_del" json:"isDel"`
	Children   []Dept `gorm:"-"json:"children"`
}

type DeptLable struct {
	Id       int64       `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLable `gorm:"-" json:"children"`
}

func (e *Dept) Create() (Dept, error) {
	var doc Dept
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	e.IsDel = 0
	result := orm.Eloquent.Table("sys_dept").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	deptPath := "/" + utils.Int64ToString(e.Deptid)
	if int(e.ParentId) != 0 {
		var deptP Dept
		orm.Eloquent.Table("sys_dept").Where("dept_id = ?", e.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	var mp = map[string]string{}
	mp["deptPath"] = deptPath
	if err := orm.Eloquent.Table("sys_dept").Where("dept_id = ?", e.Deptid).Update(mp).Error; err != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	doc.DeptPath = deptPath
	return doc, nil
}

func (e *Dept) Get() (Dept, error) {
	var doc Dept

	table := orm.Eloquent.Table("sys_dept")
	if e.Deptid != 0 {
		table = table.Where("dept_id = ?", e.Deptid)
	}
	if e.Deptname != "" {
		table = table.Where("dept_name = ?", e.Deptname)
	}

	if err := table.Where("is_del = 0").First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Dept) GetList() ([]Dept, error) {
	var doc []Dept

	table := orm.Eloquent.Table("sys_dept")
	if e.Deptid != 0 {
		table = table.Where("dept_id = ?", e.Deptid)
	}
	if e.Deptname != "" {
		table = table.Where("dept_name = ?", e.Deptname)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.Where("is_del = 0").Order("sort").Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *Dept) GetPage(bl bool) ([]Dept, error) {
	var doc []Dept

	table := orm.Eloquent.Select("*").Table("sys_dept")
	if e.Deptid != 0 {
		table = table.Where("dept_id = ?", e.Deptid)
	}
	if e.Deptname != "" {
		table = table.Where("dept_name = ?", e.Deptname)
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
		dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
		table = dataPermission.GetDataScope("sys_dept", table)
	}

	if err := table.Where("is_del = 0").Order("sort").Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *Dept) SetDept(bl bool) ([]Dept, error) {
	list, err := e.GetPage(bl)

	m := make([]Dept, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := Digui(&list, list[i])

		m = append(m, info)
	}
	return m, err
}

func Digui(deptlist *[]Dept, menu Dept) Dept {
	list := *deptlist

	min := make([]Dept, 0)
	for j := 0; j < len(list); j++ {

		if menu.Deptid != list[j].ParentId {
			continue
		}
		mi := Dept{}
		mi.Deptid = list[j].Deptid
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.Deptname = list[j].Deptname
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.CreateTime = list[j].CreateTime
		mi.UpdateTime = list[j].UpdateTime
		mi.IsDel = list[j].IsDel
		mi.Children = []Dept{}
		ms := Digui(deptlist, mi)
		min = append(min, ms)

	}
	menu.Children = min
	return menu
}

func (e *Dept) Update(id int64) (update Dept, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_dept").Where("dept_id = ?", id).First(&update).Error; err != nil {
		return
	}

	deptPath := "/" + utils.Int64ToString(e.Deptid)
	if int(e.ParentId) != 0 {
		var deptP Dept
		orm.Eloquent.Table("sys_dept").Where("dept_id = ?", e.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	e.DeptPath = deptPath

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_dept").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *Dept) Delete(id int64) (success bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_dept").Where("dept_id = ?", id).Update(mp).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (dept *Dept) SetDeptLable() (m []DeptLable, err error) {
	deptlist, err := dept.GetList()

	m = make([]DeptLable, 0)
	for i := 0; i < len(deptlist); i++ {
		if deptlist[i].ParentId != 0 {
			continue
		}
		e := DeptLable{}
		e.Id = deptlist[i].Deptid
		e.Label = deptlist[i].Deptname
		deptsInfo := DiguiDeptLable(&deptlist, e)

		m = append(m, deptsInfo)
	}
	return
}

func DiguiDeptLable(deptlist *[]Dept, dept DeptLable) DeptLable {
	list := *deptlist

	min := make([]DeptLable, 0)
	for j := 0; j < len(list); j++ {

		if dept.Id != list[j].ParentId {
			continue
		}
		mi := DeptLable{list[j].Deptid, list[j].Deptname, []DeptLable{}}
		ms := DiguiDeptLable(deptlist, mi)
		min = append(min, ms)

	}
	dept.Children = min
	return dept
}
