package actions

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/jwtauth/user"
	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/response"
	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"gorm.io/gorm"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permission() below returns the query untouched when data permission
		// is off, so the lookup that feeds it has nothing to feed. It used to
		// run anyway: a sys_user join on every list, detail, update and delete,
		// with the result discarded.
		if !config.ApplicationConfig.EnableDP {
			c.Set(PermissionKey, new(DataPermission))
			c.Next()
			return
		}

		userId := user.GetUserIdStr(c)
		if userId == "" {
			c.Set(PermissionKey, new(DataPermission))
			c.Next()
			return
		}

		// The token already carries what the scope is decided by. Reading it
		// there costs nothing, and goes no more stale than rolekey does - which
		// Casbin has always read from the token.
		if p, ok := permissionFromClaims(c); ok {
			c.Set(PermissionKey, p)
			c.Next()
			return
		}

		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			// Same fix as the newDataPermission branch below: without Abort,
			// gin's "return means continue" semantics send the request on to
			// the business handler with PermissionKey never set. The caller
			// then reads a zero-value DataPermission, which used to fall
			// into Permission()'s fail-open default - a database hiccup
			// silently turning into "see everything". PRD 006 F14/H1.
			response.Error(c, 500, err, "权限范围鉴定错误")
			c.Abort()
			return
		}
		msgID := pkg.GenerateMsgIDFromContext(c)
		p, err := newDataPermission(db, userId)
		if err != nil {
			log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
			response.Error(c, 500, err, "权限范围鉴定错误")
			c.Abort()
			return
		}
		c.Set(PermissionKey, p)
		c.Next()
	}
}

// permissionFromClaims builds the scope from the token, reporting false when
// the token predates deptid being carried. Such a token still exists until it
// expires, and it has to keep working.
func permissionFromClaims(c *gin.Context) (*DataPermission, bool) {
	claims := user.ExtractClaims(c)
	if claims["deptid"] == nil || claims["datascope"] == nil {
		return nil, false
	}
	scope, ok := claims["datascope"].(string)
	if !ok {
		return nil, false
	}
	return &DataPermission{
		DataScope: scope,
		UserId:    user.GetUserId(c),
		DeptId:    user.GetDeptId(c),
		RoleId:    user.GetRoleId(c),
	}, true
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*DataPermission, error) {
	var err error
	p := &DataPermission{}

	err = tx.Table("sys_user").
		Select("sys_user.user_id", "sys_role.role_id", "sys_user.dept_id", "sys_role.data_scope").
		Joins("left join sys_role on sys_role.role_id = sys_user.role_id").
		Where("sys_user.user_id = ?", userId).
		Scan(p).Error
	if err != nil {
		err = errors.New("获取用户数据出错 msg:" + err.Error())
		return nil, err
	}
	return p, nil
}

func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
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
			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%/"+pkg.IntToString(p.DeptId)+"/%")
		case "5":
			return db.Where(tableName+".create_by = ?", p.UserId)
		default:
			return db
		}
	}
}

func getPermissionFromContext(c *gin.Context) *DataPermission {
	p := new(DataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch pm.(type) {
		case *DataPermission:
			p = pm.(*DataPermission)
		}
	}
	return p
}

// GetPermissionFromContext 提供非action写法数据范围约束
func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return getPermissionFromContext(c)
}
