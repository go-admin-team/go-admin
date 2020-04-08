package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type SysRole struct {
	// 角色编码
	Id int64 `json:"roleId" gorm:"column:role_id;primary_key"`
	// 角色名称
	Name string `json:"roleName" gorm:"column:role_name"`

	Status string `json:"status" gorm:"column:status"`

	//角色代码
	RoleKey string `json:"roleKey" gorm:"column:role_key"`
	//角色排序
	RoleSort int64 `json:"roleSort" gorm:"column:role_sort"`

	Flag string `json:"flag" gorm:"column:flag"`

	CreateBy   string `gorm:"column:create_by" json:"createBy"`
	CreateTime string `gorm:"column:create_time" json:"createTime"`
	UpdateBy   string `gorm:"column:update_by" json:"updateBy"`
	UpdateTime string `gorm:"column:update_time" json:"updateTime"`

	//备注
	Remark    string `gorm:"column:remark" json:"remark"`
	Params    string `gorm:"column:params" json:"params"`
	DataScope string `gorm:"-" json:"dataScope"`
	IsDel     int `gorm:"column:is_del" json:"isDel"`

	Admin bool `gorm:"column:admin" json:"admin"`

	MenuIds []int64 `gorm:"-" json:"menuIds"`
	DeptIds []int64 `gorm:"-" json:"deptIds"`
}

type MenuIdList struct {
	MenuId int64 `json:"menuId"`
}

func (e *SysRole) GetPage(pageSize int, pageIndex int) ([]SysRole, int32, error) {
	var doc []SysRole

	table := orm.Eloquent.Select("*").Table("sys_role")
	if e.Id != 0 {
		table = table.Where("role_id = ?", e.Id)
	}
	if e.Name != "" {
		table = table.Where("role_name = ?", e.Name)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.RoleKey != "" {
		table = table.Where("role_key = ?", e.RoleKey)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
	table = dataPermission.GetDataScope("sys_role", table)

	var count int32

	if err := table.Where("is_del = 0").Order("role_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("is_del = 0").Count(&count)
	return doc, count, nil
}

func (role *SysRole) Get() (SysRole SysRole, err error) {
	table := orm.Eloquent.Table("sys_role")
	if role.Id != 0 {
		table = table.Where("role_id = ?", role.Id)
	}
	if role.Name != "" {
		table = table.Where("role_name = ?", role.Name)
	}
	if err = table.First(&SysRole).Error; err != nil {
		return
	}

	return
}

func (role *SysRole) GetList() (SysRole []SysRole, err error) {
	table := orm.Eloquent.Table("sys_role").Where("is_del = 0")
	if role.Id != 0 {
		table = table.Where("role_id = ?", role.Id)
	}
	if role.Name != "" {
		table = table.Where("role_name = ?", role.Name)
	}
	if err = table.Order("role_sort").Find(&SysRole).Error; err != nil {
		return
	}

	return
}

func (role *SysRole) GetRoleMeunId() ([]int64, error) {
	menuIds := make([]int64, 0)
	menuList := make([]MenuIdList, 0)
	if err := orm.Eloquent.Table("sys_role_menu").Select("sys_role_menu.menu_id").Joins("LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id").Where("role_id = ? ", role.Id).Where(" sys_role_menu.menu_id not in(select sys_menu.parent_id from sys_role_menu LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id where role_id =? )", role.Id).Find(&menuList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

func (role *SysRole) Insert() (id int64, err error) {
	role.CreateTime = utils.GetCurrntTime()
	role.UpdateBy = ""
	role.UpdateTime = utils.GetCurrntTime()
	role.IsDel = 0
	result := orm.Eloquent.Table("sys_role").Create(&role)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = role.Id
	return
}

type DeptIdList struct {
	DeptId int64 `json:"DeptId"`
}

func (role *SysRole) GetRoleDeptId() ([]int64, error) {
	deptIds := make([]int64, 0)
	deptList := make([]DeptIdList, 0)
	if err := orm.Eloquent.Table("sys_role_dept").Select("sys_role_dept.dept_id").Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").Where("role_id = ? ", role.Id).Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id where role_id =? )", role.Id).Order("role_sort").Find(&deptList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].DeptId)
	}

	return deptIds, nil
}

//修改
func (role *SysRole) Update(id int64) (update SysRole, err error) {
	role.UpdateTime = utils.GetCurrntTime()
	if err = orm.Eloquent.Table("sys_role").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_role").Model(&update).Updates(&role).Error; err != nil {
		return
	}
	return
}

//批量删除
func (e *SysRole) BatchDelete(id []int64) (Result bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_role").Where("is_del=0 and role_id in (?)", id).Update(mp).Error; err != nil {
		return
	}
	Result = true
	return
}
