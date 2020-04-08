package models

import (
	"errors"
	orm "go-admin/database"
	"go-admin/pkg"
	"go-admin/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// User
type User struct {
	// key
	IdentityKey string
	// 用户名
	UserName  string
	FirstName string
	LastName  string
	// 角色
	Role string
}

type UserName struct {
	Username string `gorm:"column:username" json:"username"`
}

type PassWord struct {
	// 密码
	Password string `gorm:"column:password" json:"password"`
}

type LoginM struct {
	UserName
	PassWord
}

type SysUserId struct {
	// 编码
	Id int64 `gorm:"column:user_id;primary_key"  json:"userId"`
}

type SysUserB struct {
	// 昵称
	NickName string `gorm:"column:nick_name" json:"nickName"`
	// 手机号
	Phone string `gorm:"column:phone" json:"phone"`
	// 角色编码
	RoleId int64 `gorm:"column:role_id" json:"roleId"`
	//盐
	Salt string `gorm:"column:salt" json:"salt"`
	//头像
	Avatar string `gorm:"column:avatar" json:"avatar"`
	//性别
	Sex string `gorm:"column:sex" json:"sex"`
	//邮箱
	Email string `gorm:"column:email" json:"email"`
	// 创建时间
	CreateTime string `gorm:"column:create_time" json:"createTime"`
	// 修改时间
	UpdateTime string `gorm:"column:update_time" json:"updateTime"`

	//部门编码
	DeptId int64 `gorm:"column:dept_id" json:"deptId"`

	//职位编码
	PostId int64 `gorm:"column:post_id" json:"postId"`

	CreateBy string `gorm:"column:create_by" json:"createBy"`
	UpdateBy string `gorm:"column:update_by" json:"updateBy"`

	//备注
	Remark    string `gorm:"column:remark" json:"remark"`
	Params    string `gorm:"column:params" json:"params"`
	Status    string `gorm:"column:status" json:"status"`
	DataScope string `gorm:"-" json:"dataScope"`
	IsDel     int `gorm:"column:is_del" json:"isDel"`
}

type SysUser struct {
	SysUserId
	SysUserB
	LoginM
}

type SysUserPwd struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type SysUserPage struct {
	SysUserId
	SysUserB
	LoginM
	DeptName string `gorm:"-" json:"deptName"`
}

type SysUserView struct {
	SysUserId
	SysUserB
	LoginM
	RoleName string `gorm:"column:role_name"  json:"role_name"`
}

// 获取用户数据
func (e *SysUser) Get() (SysUserView SysUserView, err error) {

	table := orm.Eloquent.Table("sys_user").Select([]string{"sys_user.*", "sys_role.role_name"})
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if e.Id != 0 {
		table = table.Where("user_id = ?", e.Id)
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
	return
}

func (e *SysUser) GetPage(pageSize int, pageIndex int) ([]SysUserPage, int32, error) {
	var doc []SysUserPage

	table := orm.Eloquent.Select("sys_user.*,sys_dept.dept_name").Table("sys_user")
	table = table.Joins("left join sys_dept on sys_dept.dept_id = sys_user.dept_id")

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.DeptId != 0 {
		table = table.Where("sys_user.dept_id in (select dept_id from sys_dept where dept_path like ? )", "%"+utils.Int64ToString(e.DeptId)+"%")
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = utils.StringToInt64(e.DataScope)
	table = dataPermission.GetDataScope("sys_user", table)

	var count int32

	if err := table.Where("sys_user.is_del = 0").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("sys_user.is_del = 0").Count(&count)
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
func (e SysUser) Insert() (id int64, err error) {
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	e.IsDel = "0"
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int
	orm.Eloquent.Table("sys_user").Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = orm.Eloquent.Table("sys_user").Create(&e).Error; err != nil {
		return
	}
	id = e.Id
	return
}

//修改
func (e *SysUser) Update(id int64) (update SysUser, err error) {
	e.UpdateTime = utils.GetCurrntTime()
	if err = e.Encrypt(); err != nil {
		return
	}

	if err = orm.Eloquent.Table("sys_user").First(&update, id).Error; err != nil {
		return
	}
	if e.RoleId == 0 {
		e.RoleId = update.RoleId
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_user").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysUser) BatchDelete(id []int64) (Result bool, err error) {
	var mp = map[string]string{}
	mp["is_del"] = "1"
	mp["update_time"] = utils.GetCurrntTime()
	mp["update_by"] = e.UpdateBy
	if err = orm.Eloquent.Table("sys_user").Where("is_del=0 and user_id in (?)", id).Update(mp).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *SysUser) SetPwd(pwd SysUserPwd) (Result bool, err error) {
	user, _ := e.Get()
	_, err = pkg.CompareHashAndPassword(user.Password, pwd.OldPassword)
	if err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			pkg.AssertErr(err, "密码错误(代码202)", 500)
		}
		log.Print(err)
		return
	}
	e.Password = pwd.NewPassword
	_, err = e.Update(e.Id)
	pkg.AssertErr(err, "更新密码失败(代码202)", 500)
	return
}
