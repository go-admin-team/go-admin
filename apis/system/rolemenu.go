package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/models"
	"go-admin/tools/app"
	"net/http"
)

// @Summary RoleMenu列表数据
// @Description 获取JSON
// @Tags 角色菜单
// @Param RoleId query string false "RoleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/rolemenu [get]
// @Security Bearer
func GetRoleMenu(c *gin.Context) {
	var Rm models.RoleMenu
	err := c.ShouldBind(&Rm)
	result, err := Rm.Get()
	var res app.Response
	if err != nil {
		res.Msg = "抱歉未找到相关信息"
		c.JSON(http.StatusOK, res.ReturnError(200))
		return
	}
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

type RoleMenuPost struct {
	RoleId   string
	RoleMenu []models.RoleMenu
}

func InsertRoleMenu(c *gin.Context) {

	var res app.Response
	res.Msg = "添加成功"
	c.JSON(http.StatusOK, res.ReturnOK())
	return

}

// @Summary 删除用户菜单数据
// @Description 删除数据
// @Tags 角色菜单
// @Param id path string true "id"
// @Param menu_id query string false "menu_id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/rolemenu/{id} [delete]
func DeleteRoleMenu(c *gin.Context) {
	var t models.RoleMenu
	id := c.Param("id")
	menuId := c.Request.FormValue("menu_id")
	fmt.Println(menuId)
	_, err := t.Delete(id, menuId)
	if err != nil {
		var res app.Response
		res.Msg = "删除失败"
		c.JSON(http.StatusOK, res.ReturnError(200))
		return
	}
	var res app.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
	return
}
