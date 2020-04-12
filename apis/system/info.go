package system

import (
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/app"
	"go-admin/pkg/utils"
)

func GetInfo(c *gin.Context) {

	var roles = make([]string, 1)
	roles[0] = utils.GetRoleName(c)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"
	RoleMenu := models.RoleMenu{}
	RoleMenu.RoleId = utils.GetRoleId(c)

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	if utils.GetRoleName(c) == "admin" || utils.GetRoleName(c) == "系统管理员" {
		mp["permissions"] = permissions
	} else {
		list, _ := RoleMenu.GetPermis()
		mp["permissions"] = list
	}

	sysuser := models.SysUser{}
	sysuser.UserId = utils.GetUserId(c)
	user, err := sysuser.Get()
	pkg.HasError(err, "", 500)

	mp["introduction"] = " am a super administrator"

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}
	mp["name"] = user.NickName

	app.OK(c,mp,"")
}
