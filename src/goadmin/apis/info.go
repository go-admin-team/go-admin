package apis

import (
	"github.com/gin-gonic/gin"
	"goadmin/models"
	"goadmin/utils"
	"net/http"
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
	mp["introduction"] = " am a super administrator"
	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	mp["name"] = utils.GetUserName(c)

	var res models.Response
	res.Data = mp
	c.JSON(http.StatusOK, res.ReturnOK())
}
