package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/tools"
	"go-admin/tools/app"
)

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menulist [get]
// @Security Bearer
func GetMenuList(c *gin.Context) {
	var Menu models.Menu
	Menu.MenuName = c.Request.FormValue("menuName")
	Menu.Visible = c.Request.FormValue("visible")
	Menu.Title = c.Request.FormValue("title")
	Menu.DataScope = tools.GetUserIdStr(c)
	result, err := Menu.SetMenu()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menu [get]
// @Security Bearer
func GetMenu(c *gin.Context) {
	var data models.Menu
	id, err := tools.StringToInt(c.Param("id"))
	data.MenuId = id
	result, err := data.GetByMenuId()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

func GetMenuTreeRoleselect(c *gin.Context) {
	var Menu models.Menu
	var SysRole models.SysRole
	id, err := tools.StringToInt(c.Param("roleId"))
	SysRole.RoleId = id
	result, err := Menu.SetMenuLable()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleMeunId()
		tools.HasError(err, "抱歉未找到相关信息", -1)
	}
	app.Custum(c, gin.H{
		"code":        200,
		"menus":       result,
		"checkedKeys": menuIds,
	})
}

// @Summary 获取菜单树
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menuTreeselect [get]
// @Security Bearer
func GetMenuTreeelect(c *gin.Context) {
	var data models.Menu
	result, err := data.SetMenuLable()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

// @Summary 创建菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param menuName formData string true "menuName"
// @Param Path formData string false "Path"
// @Param Action formData string true "Action"
// @Param Permission formData string true "Permission"
// @Param ParentId formData string true "ParentId"
// @Param IsDel formData string true "IsDel"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/menu [post]
// @Security Bearer
func InsertMenu(c *gin.Context) {
	var data models.Menu
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "抱歉未找到相关信息", -1)
	data.CreateBy = tools.GetUserIdStr(c)
	result, err := data.Create()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

// @Summary 修改菜单
// @Description 获取JSON
// @Tags 菜单
// @Accept  application/x-www-form-urlencoded
// @Product application/x-www-form-urlencoded
// @Param id path int true "id"
// @Param data body models.Menu true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/menu/{id} [put]
// @Security Bearer
func UpdateMenu(c *gin.Context) {
	var data models.Menu
	err2 := c.BindWith(&data, binding.JSON)
	data.UpdateBy = tools.GetUserIdStr(c)
	tools.HasError(err2, "修改失败", -1)
	_, err := data.Update(data.MenuId)
	tools.HasError(err, "", 501)
	app.OK(c, "", "修改成功")

}

// @Summary 删除菜单
// @Description 删除数据
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/menu/{id} [delete]
func DeleteMenu(c *gin.Context) {
	var data models.Menu
	id, err := tools.StringToInt(c.Param("id"))
	data.UpdateBy = tools.GetUserIdStr(c)
	_, err = data.Delete(id)
	tools.HasError(err, "删除失败", 500)
	app.OK(c, "", "删除成功")
}

// @Summary 根据角色名称获取菜单列表数据（左菜单使用）
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menurole [get]
// @Security Bearer
func GetMenuRole(c *gin.Context) {
	var Menu models.Menu
	result, err := Menu.SetMenuRole(tools.GetRoleName(c))
	tools.HasError(err, "获取失败", 500)
	app.OK(c, result, "")
}

// @Summary 获取角色对应的菜单id数组
// @Description 获取JSON
// @Tags 菜单
// @Param id path int true "id"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menuids/{id} [get]
// @Security Bearer
func GetMenuIDS(c *gin.Context) {
	var data models.RoleMenu
	data.RoleName = c.GetString("role")
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.GetIDS()
	tools.HasError(err, "获取失败", 500)
	app.OK(c, result, "")
}
