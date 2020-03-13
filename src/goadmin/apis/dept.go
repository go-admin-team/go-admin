package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"goadmin/models"
	"goadmin/pkg"
	"goadmin/utils"
	"net/http"
	"strconv"
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
	Dept.Deptname = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.Deptid, _ = utils.StringToInt64(c.Request.FormValue("deptId"))
	Dept.DataScope = utils.GetUserIdStr(c)
	result, err := Dept.SetDept(true)
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)

	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
}

func GetDeptTree(c *gin.Context) {
	var Dept models.Dept
	Dept.Deptname = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.Deptid, _ = utils.StringToInt64(c.Request.FormValue("deptId"))
	result, err := Dept.SetDept(false)
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)

	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	Dept.Deptid, _ = strconv.ParseInt(c.Param("deptId"), 10, 64)
	Dept.DataScope = utils.GetUserIdStr(c)
	result, err := Dept.Get()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)

	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	pkg.AssertErr(err, "", 500)
	data.CreateBy = utils.GetUserIdStr(c)
	result, err := data.Create()
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	pkg.AssertErr(err, "", -1)
	data.UpdateBy = utils.GetUserIdStr(c)
	result, err := data.Update(data.Deptid)
	pkg.AssertErr(err, "", -1)
	var res models.Response
	res.Data = result
	c.JSON(http.StatusOK, res.ReturnOK())
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	_, err = data.Delete(id)
	pkg.AssertErr(err, "删除失败", 500)

	var res models.Response
	res.Msg = "删除成功"
	c.JSON(http.StatusOK, res.ReturnOK())
}

func GetDeptTreeRoleselect(c *gin.Context) {
	var Dept models.Dept
	var SysRole models.SysRole
	id, err := strconv.ParseInt(c.Param("roleId"), 10, 64)
	SysRole.Id = id
	result, err := Dept.SetDeptLable()
	pkg.AssertErr(err, "抱歉未找到相关信息", -1)
	menuIds := make([]int64, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleDeptId()
		pkg.AssertErr(err, "抱歉未找到相关信息", -1)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"depts":       result,
		"checkedKeys": menuIds,
	})
}
