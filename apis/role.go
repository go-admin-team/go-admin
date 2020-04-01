package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/utils"
	"net/http"
	"strconv"
)

// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/rolelist [get]
// @Security
func GetRoleList(c *gin.Context) {
	var data models.SysRole
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	data.RoleKey = c.Request.FormValue("roleKey")
	data.Name = c.Request.FormValue("roleName")
	data.Status = c.Request.FormValue("status")
	data.DataScope = utils.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	pkg.AssertErr(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
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
	Role.Id, _ = strconv.ParseInt(c.Param("roleId"), 10, 64)
	result, err := Role.Get()

	menuIds := make([]int64, 0)

	menuIds, err = Role.GetRoleMeunId()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)

	result.MenuIds = menuIds
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "抱歉未找到相关信息",
		})
		return
	}
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.Config true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
func InsertRole(c *gin.Context) {
	var data models.SysRole
	data.CreateBy = utils.GetUserIdStr(c)
	err := c.BindWith(&data, binding.JSON)
	pkg.AssertErr(err, "", 500)
	id, err := data.Insert()
	data.Id = id
	pkg.AssertErr(err, "", -1)
	var t models.RoleMenu
	_, err = t.Insert(id, data.MenuIds)
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = data
	res.Msg = "添加成功"
	c.JSON(http.StatusOK, res.ReturnOK())
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
	data.UpdateBy = utils.GetUserIdStr(c)
	err := c.Bind(&data)
	pkg.AssertErr(err, "数据解析失败", -1)
	result, err := data.Update(data.Id)

	var t models.RoleMenu
	_, err = t.DeleteRoleMenu(data.Id)
	pkg.AssertErr(err, "添加失败1", -1)

	_, err2 := t.Insert(data.Id, data.MenuIds)
	pkg.AssertErr(err2, "添加失败2", -1)

	var res models.Response
	res.Data = result
	res.Msg = "修改成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}

func UpdateRoleDataScope(c *gin.Context) {
	var data models.SysRole
	data.UpdateBy = utils.GetUserIdStr(c)
	err := c.Bind(&data)
	pkg.AssertErr(err, "数据解析失败", -1)
	result, err := data.Update(data.Id)

	var t models.SysRoleDept
	_, err = t.DeleteRoleDept(data.Id)
	pkg.AssertErr(err, "添加失败1", -1)
	if data.DataScope == "2" {
		_, err2 := t.Insert(data.Id, data.DeptIds)
		pkg.AssertErr(err2, "添加失败2", -1)
	}
	var res models.Response
	res.Data = result
	res.Msg = "修改成功"
	c.JSON(http.StatusOK, res.ReturnOK())
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
	Role.UpdateBy = utils.GetUserIdStr(c)

	IDS := utils.IdsStrToIdsInt64Group("roleId", c)
	_, err := Role.BatchDelete(IDS)
	pkg.AssertErr(err, "删除失败1", -1)
	var t models.RoleMenu
	_, err = t.BatchDeleteRoleMenu(IDS)
	pkg.AssertErr(err, "删除失败1", -1)

	var res models.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}
