package apis

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"go-admin/app/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysDept struct {
	api.Api
}

// GetSysDeptList
// @Summary 分页部门列表数据
// @Description 分页列表
// @Tags 部门
// @Param name query string false "name"
// @Param id query string false "id"
// @Param position query string false "position"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept [get]
// @Security Bearer
func (e SysDept) GetSysDeptList(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}
	list := make([]models.SysDept, 0)

	list, err = s.SetDeptPage(&req)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(list, "查询成功")
}

// GetSysDept
// @Summary 部门列表数据
// @Description 获取JSON
// @Tags 部门
// @Param deptId path string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept/{deptId} [get]
// @Security Bearer
func (e SysDept) GetSysDept(c *gin.Context) {
	s :=service.SysDept{}
	req := dto.SysDeptById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}
	var object models.SysDept

	err = s.GetSysDept(&req, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

// InsertSysDept 添加部门
// @Summary 添加部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysDeptControl true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept [post]
// @Security Bearer
func (e SysDept) InsertSysDept(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.InsertSysDept(&req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// UpdateSysDept
// @Summary 修改部门
// @Description 获取JSON
// @Tags 部门
// @Accept  application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysDeptControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/dept/{deptId} [put]
// @Security Bearer
func (e SysDept) UpdateSysDept(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	req.SetUpdateBy(user.GetUserId(c))

	err = s.UpdateSysDept(&req)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// DeleteSysDept
// @Summary 删除部门
// @Description 删除数据
// @Tags 部门
// @Param data body dto.SysDeptById true "body"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dept [delete]
func (e SysDept) DeleteSysDept(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	err = s.RemoveSysDept(&req)
	if err != nil {
		e.Logger.Errorf("RemoveSysDept error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// GetDeptTree 用户管理 左侧部门树
func (e SysDept) GetDeptTree(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	list := make([]dto.DeptLabel, 0)

	list, err = s.SetDeptTree(&req)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
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
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	id, err := pkg.StringToInt(c.Param("roleId"))
	result, err := s.SetDeptLabel()
	if err != nil {
		e.Logger.Errorf("SetDeptLabel error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "")
	}
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = s.GetRoleDeptId(id)
		if err != nil {
			e.Logger.Errorf("抱歉未找到相关信息, %s", err.Error())
			e.Error(http.StatusInternalServerError, err, "")
		}
	}
	e.OK(gin.H{
		"depts":       result,
		"checkedKeys": menuIds,
	}, "")
}
