package service

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go-admin/global/orm"
	. "go-admin/models"
	"go-admin/tools"
	"golang.org/x/crypto/bcrypt"
	"strings"
)


// 获取用户数据
func (e *SysUser) Get() (SysUserView SysUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&SysUserView).Error; err != nil {
		return
	}
	SysUserView.Password = ""
	return
}

func (e *SysUser) GetPage(pageSize int, pageIndex int) ([]SysUserPage, int, error) {
	var doc []SysUserPage
	table := orm.Eloquent.Select("sys_user.*,sys_dept.dept_name").Table(e.TableName())
	table = table.Joins("left join sys_dept on sys_dept.dept_id = sys_user.dept_id")

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	if e.Status != "" {
		table = table.Where("sys_user.status = ?", e.Status)
	}

	if e.Phone != "" {
		table = table.Where("sys_user.phone = ?", e.Phone)
	}

	if e.DeptId != 0 {
		table = table.Where("sys_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+tools.IntToString(e.DeptId)+"%")
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("sys_user", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("sys_user.deleted_at IS NULL").Count(&count)
	return doc, count, nil
}

//加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

//添加
func (e SysUser) Insert() (id int, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int
	orm.Eloquent.Table(e.TableName()).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
		return
	}
	id = e.UserId
	return
}

//修改
func (e *SysUser) Update(id int) (update SysUser, err error) {
	if e.Password!="" {
		if err = e.Encrypt(); err != nil {
			return
		}
	}
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if e.RoleId == 0 {
		e.RoleId = update.RoleId
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("user_id in (?)", id).Delete(&SysUser{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *SysUser) SetPwd(pwd SysUserPwd) (Result bool, err error) {
	user, err := e.Get()
	if err != nil {
		tools.HasError(err, "获取用户数据失败(代码202)", 500)
	}
	_, err = tools.CompareHashAndPassword(user.Password, pwd.OldPassword)
	if err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			tools.HasError(err, "密码错误(代码202)", 500)
		}
		log.Print(err)
		return
	}
	e.Password = pwd.NewPassword
	_, err = e.Update(e.UserId)
	tools.HasError(err, "更新密码失败(代码202)", 500)
	return
}
