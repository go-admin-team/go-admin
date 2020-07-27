package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/tools"
	"go-admin/tools/app"
)

// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/rolelist [get]
// @Security Bearer
func GetRoleList(c *gin.Context) {
	var data models.SysRole
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.RoleKey = c.Request.FormValue("roleKey")
	data.RoleName = c.Request.FormValue("roleName")
	data.Status = c.Request.FormValue("status")
	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role [get]
// @Security Bearer
func GetRole(c *gin.Context) {
	var Role models.SysRole
	Role.RoleId, _ = tools.StringToInt(c.Param("roleId"))
	result, err := Role.Get()
	menuIds := make([]int, 0)
	menuIds, err = Role.GetRoleMeunId()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	result.MenuIds = menuIds
	app.OK(c, result, "")

}

// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
func InsertRole(c *gin.Context) {
	var data models.SysRole
	data.CreateBy = tools.GetUserIdStr(c)
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "", 500)
	id, err := data.Insert()
	data.RoleId = id
	tools.HasError(err, "", -1)
	var t models.RoleMenu
	_, err = t.Insert(id, data.MenuIds)
	tools.HasError(err, "", -1)
	app.OK(c, data, "添加成功")
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role [put]
func UpdateRole(c *gin.Context) {
	var data models.SysRole
	data.UpdateBy = tools.GetUserIdStr(c)
	err := c.Bind(&data)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.RoleId)
	tools.HasError(err, "", -1)
	var t models.RoleMenu
	_, err = t.DeleteRoleMenu(data.RoleId)
	tools.HasError(err, "添加失败1", -1)
	_, err2 := t.Insert(data.RoleId, data.MenuIds)
	tools.HasError(err2, "添加失败2", -1)

	app.OK(c, result, "修改成功")
}

func UpdateRoleDataScope(c *gin.Context) {
	var data models.SysRole
	data.UpdateBy = tools.GetUserIdStr(c)
	err := c.Bind(&data)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.RoleId)

	var t models.SysRoleDept
	_, err = t.DeleteRoleDept(data.RoleId)
	tools.HasError(err, "添加失败1", -1)
	if data.DataScope == "2" {
		_, err2 := t.Insert(data.RoleId, data.DeptIds)
		tools.HasError(err2, "添加失败2", -1)
	}
	app.OK(c, result, "修改成功")
}

// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param roleId path int true "roleId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role/{roleId} [delete]
func DeleteRole(c *gin.Context) {
	var Role models.SysRole
	Role.UpdateBy = tools.GetUserIdStr(c)

	IDS := tools.IdsStrToIdsIntGroup("roleId", c)
	_, err := Role.BatchDelete(IDS)
	tools.HasError(err, "删除失败1", -1)
	var t models.RoleMenu
	_, err = t.BatchDeleteRoleMenu(IDS)
	tools.HasError(err, "删除失败1", -1)
	app.OK(c, "", "删除成功")
}
