package actions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/tools/config"

	"gorm.io/gorm"

	"go-admin/tools"
)

type dataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		switch idb.(type) {
		case *gorm.DB:
			db := idb.(*gorm.DB)
			var p = new(dataPermission)
			if userId := tools.GetUserIdStr(c); userId != "" {
				p, err = newDataPermission(db, userId)
				tools.HasError(err, "权限范围鉴定错误", 500)
			}
			c.Set(PermissionKey, p)
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		c.Next()
	}
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*dataPermission, error) {
	var err error
	p := &dataPermission{}
	sysUser := new(models.SysUser)
	sysRole := new(models.SysRole)

	err = sysUser.GetByUserId(tx, userId)
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}
	p.UserId = sysUser.UserId
	p.RoleId = sysUser.RoleId
	p.DeptId = sysUser.DeptId
	err = sysRole.GetById(tx, sysUser.RoleId)
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}
	p.DataScope = sysRole.DataScope
	return p, nil
}

func Permission(tableName string, p *dataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !config.ApplicationConfig.EnableDP {
			return db
		}
		switch p.DataScope {
		case "2":
			return db.Where(tableName+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", p.RoleId)
		case "3":
			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where dept_id = ? )", p.DeptId)
		case "4":
			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+tools.IntToString(p.DeptId)+"%")
		case "5":
			return db.Where(tableName+".create_by = ?", p.UserId)
		default:
			return db
		}
	}
}

func getPermissionFromContext(c *gin.Context) *dataPermission {
	p := new(dataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch pm.(type) {
		case *dataPermission:
			p = pm.(*dataPermission)
		}
	}
	return p
}
