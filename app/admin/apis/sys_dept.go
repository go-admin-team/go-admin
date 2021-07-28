package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/admin/models"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysDept struct {
	api.Api
}

// GetPage
// @Summary 分页部门列表数据
// @Description 分页列表
// @Tags 部门
// @Param deptName query string false "deptName"
// @Param deptId query string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept [get]
// @Security Bearer
func (e SysDept) GetPage(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.SysDept, 0)
	list, err = s.SetDeptPage(&req)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

// Get
// @Summary 获取部门数据
// @Description 获取JSON
// @Tags 部门
// @Param deptId path string false "deptId"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept/{deptId} [get]
// @Security Bearer
func (e SysDept) Get(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptGetReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.SysDept

	err = s.Get(&req, &object)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.OK(object, "查询成功")
}

// Insert 添加部门
// @Summary 添加部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysDeptInsertReq true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept [post]
// @Security Bearer
func (e SysDept) Insert(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptInsertReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(req.GetId(), "创建成功")
}

// Update
// @Summary 修改部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysDeptUpdateReq true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept/{deptId} [put]
// @Security Bearer
func (e SysDept) Update(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	err = s.Update(&req)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// Delete
// @Summary 删除部门
// @Description 删除数据
// @Tags 部门
// @Param data body dto.SysDeptDeleteReq true "body"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dept [delete]
// @Security Bearer
func (e SysDept) Delete(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	err = s.Remove(&req)
	if err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Get2Tree 用户管理 左侧部门树
func (e SysDept) Get2Tree(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req,binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	list := make([]dto.DeptLabel, 0)
	list, err = s.SetDeptTree(&req)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "")
}

// GetDeptTreeRoleSelect TODO: 此接口需要调整不应该将list和选中放在一起
func (e SysDept) GetDeptTreeRoleSelect(c *gin.Context) {
	s := service.SysDept{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	id, err := pkg.StringToInt(c.Param("roleId"))
	result, err := s.SetDeptLabel()
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = s.GetWithRoleId(id)
		if err != nil {
			e.Error(500, err, err.Error())
			return
		}
	}
	e.OK(gin.H{
		"depts":       result,
		"checkedKeys": menuIds,
	}, "")
}