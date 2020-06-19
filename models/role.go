package models

import (
	"github.com/pkg/errors"
	"go-admin/global/orm"
	"go-admin/tools"
)

type SysRole struct {
	RoleId    int    `json:"roleId" gorm:"primary_key;AUTO_INCREMENT"` // 角色编码
	RoleName  string `json:"roleName" gorm:"type:varchar(128);"`       // 角色名称
	Status    string `json:"status" gorm:"type:int(1);"`               //
	RoleKey   string `json:"roleKey" gorm:"type:varchar(128);"`        //角色代码
	RoleSort  int    `json:"roleSort" gorm:"type:int(4);"`             //角色排序
	Flag      string `json:"flag" gorm:"type:varchar(128);"`           //
	CreateBy  string `json:"createBy" gorm:"type:varchar(128);"`       //
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(128);"`       //
	Remark    string `json:"remark" gorm:"type:varchar(255);"`         //备注
	Admin     bool   `json:"admin" gorm:"type:char(1);"`
	DataScope string `json:"dataScope" gorm:"type:varchar(128);"`
	Params    string `json:"params" gorm:"-"`
	MenuIds   []int  `json:"menuIds" gorm:"-"`
	DeptIds   []int  `json:"deptIds" gorm:"-"`
	BaseModel
}

func (SysRole) TableName() string {
	return "sys_role"
}

type MenuIdList struct {
	MenuId int `json:"menuId"`
}

func (e *SysRole) GetPage(pageSize int, pageIndex int) ([]SysRole, int, error) {
	var doc []SysRole

	table := orm.Eloquent.Select("*").Table("sys_role")
	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}
	if e.RoleName != "" {
		table = table.Where("role_name = ?", e.RoleName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.RoleKey != "" {
		table = table.Where("role_key = ?", e.RoleKey)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table,err := dataPermission.GetDataScope("sys_role", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Order("role_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

func (role *SysRole) Get() (SysRole SysRole, err error) {
	table := orm.Eloquent.Table("sys_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if err = table.First(&SysRole).Error; err != nil {
		return
	}

	return
}

func (role *SysRole) GetList() (SysRole []SysRole, err error) {
	table := orm.Eloquent.Table("sys_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if err = table.Order("role_sort").Find(&SysRole).Error; err != nil {
		return
	}

	return
}

// 获取角色对应的菜单ids
func (role *SysRole) GetRoleMeunId() ([]int, error) {
	menuIds := make([]int, 0)
	menuList := make([]MenuIdList, 0)
	if err := orm.Eloquent.Table("sys_role_menu").Select("sys_role_menu.menu_id").Joins("LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id").Where("role_id = ? ", role.RoleId).Where(" sys_role_menu.menu_id not in(select sys_menu.parent_id from sys_role_menu LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id where role_id =? )", role.RoleId).Find(&menuList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

func (role *SysRole) Insert() (id int, err error) {
	i := 0
	orm.Eloquent.Table("sys_role").Where("role_name=? or role_key = ?", role.RoleName, role.RoleKey).Count(&i)
	if i > 0 {
		return 0, errors.New("角色名称或者角色标识已经存在！")
	}
	role.UpdateBy = ""
	result := orm.Eloquent.Table("sys_role").Create(&role)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = role.RoleId
	return
}

type DeptIdList struct {
	DeptId int `json:"DeptId"`
}

func (role *SysRole) GetRoleDeptId() ([]int, error) {
	deptIds := make([]int, 0)
	deptList := make([]DeptIdList, 0)
	if err := orm.Eloquent.Table("sys_role_dept").Select("sys_role_dept.dept_id").Joins("LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id").Where("role_id = ? ", role.RoleId).Where(" sys_role_dept.dept_id not in(select sys_dept.parent_id from sys_role_dept LEFT JOIN sys_dept on sys_dept.dept_id=sys_role_dept.dept_id where role_id =? )", role.RoleId).Find(&deptList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].DeptId)
	}

	return deptIds, nil
}

//修改
func (role *SysRole) Update(id int) (update SysRole, err error) {
	if err = orm.Eloquent.Table("sys_role").First(&update, id).Error; err != nil {
		return
	}

	if role.RoleName != "" && role.RoleName != update.RoleName {
		return update, errors.New("角色名称不允许修改！")
	}

	if role.RoleKey != "" && role.RoleKey != update.RoleKey {
		return update, errors.New("角色标识不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_role").Model(&update).Updates(&role).Error; err != nil {
		return
	}
	return
}

//批量删除
func (e *SysRole) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table("sys_role").Where("role_id in (?)", id).Delete(&SysRole{}).Error; err != nil {
		return
	}
	Result = true
	return
}
