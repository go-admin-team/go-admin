package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-admin/models"
	"go-admin/pkg"
	"go-admin/pkg/app"
	"go-admin/pkg/app/msg"
	"go-admin/pkg/utils"
)

// @Summary 分页部门列表数据
// @Description 分页列表
// @Tags 部门
// @Param name query string false "name"
// @Param id query string false "id"
// @Param position query string false "position"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/deptList [get]
// @Security
func GetDeptList(c *gin.Context) {
	var Dept models.Dept
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId, _ = utils.StringToInt(c.Request.FormValue("deptId"))
	Dept.DataScope = utils.GetUserIdStr(c)
	result, err := Dept.SetDept(true)
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c,result,"")
}

func GetDeptTree(c *gin.Context) {
	var Dept models.Dept
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status= c.Request.FormValue("status")
	Dept.DeptId, _ = utils.StringToInt(c.Request.FormValue("deptId"))
	result, err := Dept.SetDept(false)
	pkg.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c,result,"")
}

// @Summary 部门列表数据
// @Description 获取JSON
// @Tags 部门
// @Param deptId path string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} models.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept/{deptId} [get]
// @Security
func GetDept(c *gin.Context) {
	var Dept models.Dept
	Dept.DeptId, _ = utils.StringToInt(c.Param("deptId"))
	Dept.DataScope = utils.GetUserIdStr(c)
	result, err := Dept.Get()
	pkg.HasError(err, msg.NotFound, 404)
	app.OK(c,result,msg.GetSuccess)
}

// @Summary 添加部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param data body models.Dept true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept [post]
// @Security Bearer
func InsertDept(c *gin.Context) {
	var data models.Dept
	err := c.BindWith(&data, binding.JSON)
	pkg.HasError(err, "", 500)
	data.CreateBy = utils.GetUserIdStr(c)
	result, err := data.Create()
	pkg.HasError(err, "", -1)
	app.OK(c,result,msg.CreatedSuccess)
}

// @Summary 修改部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body models.Dept true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept [put]
// @Security Bearer
func UpdateDept(c *gin.Context) {
	var data models.Dept
	err := c.BindWith(&data, binding.JSON)
	pkg.HasError(err, "", -1)
	data.UpdateBy = utils.GetUserIdStr(c)
	result, err := data.Update(data.DeptId)
	pkg.HasError(err, "", -1)
	app.OK(c,result,msg.UpdatedSuccess)
}

// @Summary 删除部门
// @Description 删除数据
// @Tags 部门
// @Param id path int true "id"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dept/{id} [delete]
func DeleteDept(c *gin.Context) {
	var data models.Dept
	id, err := utils.StringToInt(c.Param("id"))
	_, err = data.Delete(id)
	pkg.HasError(err, "删除失败", 500)
	app.OK(c,"",msg.DeletedSuccess)
}

func GetDeptTreeRoleselect(c *gin.Context) {
	var Dept models.Dept
	var SysRole models.SysRole
	id, err := utils.StringToInt(c.Param("roleId"))
	SysRole.RoleId = id
	result, err := Dept.SetDeptLable()
	pkg.HasError(err, msg.NotFound, -1)
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleDeptId()
		pkg.HasError(err, "抱歉未找到相关信息", -1)
	}
	app.Custum(c,gin.H{
		"code":        200,
		"depts":       result,
		"checkedKeys": menuIds,
	})
}
