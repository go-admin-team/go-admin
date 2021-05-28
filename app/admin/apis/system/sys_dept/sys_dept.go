package sys_dept

import (
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
)

type SysDept struct {
	apis.Api
}

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
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysDeptSearch)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	list := make([]system.SysDept, 0)
	serviceStudent := service.SysDept{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	list, err = serviceStudent.SetDeptPage(d)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(list, "查询成功")
}

// @Summary 部门列表数据
// @Description 获取JSON
// @Tags 部门
// @Param deptId path string false "deptId"
// @Param position query string false "position"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept/{deptId} [get]
// @Security Bearer
func (e SysDept) GetSysDept(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysDeptById)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object system.SysDept

	serviceSysOperlog := service.SysDept{}
	serviceSysOperlog.Log = log
	serviceSysOperlog.Orm = db
	err = serviceSysOperlog.GetSysDept(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

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
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysDeptControl)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	err = control.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.SetCreateBy(user.GetUserId(c))

	serviceSysDept := service.SysDept{}
	serviceSysDept.Orm = db
	serviceSysDept.Log = log
	err = serviceSysDept.InsertSysDept(object)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(object.GetId(), "创建成功")
}

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
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysDeptControl)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//更新操作
	err = control.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	serviceSysDept := service.SysDept{}
	serviceSysDept.Orm = db
	serviceSysDept.Log = log
	err = serviceSysDept.UpdateSysDept(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(object.GetId(), "更新成功")
}

// @Summary 删除部门
// @Description 删除数据
// @Tags 部门
// @Param data body dto.SysDeptById true "body"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/dept [delete]
func (e SysDept) DeleteSysDept(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysDeptById)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("Bind error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysDept := service.SysDept{}
	serviceSysDept.Orm = db
	serviceSysDept.Log = log
	err = serviceSysDept.RemoveSysDept(control)
	if err != nil {
		log.Errorf("RemoveSysDept error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(control.GetId(), "删除成功")
}

// GetDeptTree 用户管理 左侧部门树
func (e SysDept) GetDeptTree(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysDeptSearch)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	list := make([]dto.DeptLabel, 0)
	serviceStudent := service.SysDept{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	list, err = serviceStudent.SetDeptTree(d)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	//var Dept models.SysDepts
	//Dept.DeptName = c.Request.FormValue("deptName")
	//Dept.Status = c.Request.FormValue("status")
	//Dept.DeptId, _ = tools.StringToInt(c.Request.FormValue("deptId"))
	//result, err := Dept.SetDept(false)
	//tools.HasError(err, "抱歉未找到相关信息", -1)
	e.OK(list, "")
}

func (e SysDept) GetDeptTreeRoleSelect(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	s := service.SysDept{}
	s.Orm = db
	s.Log = log
	id, err := pkg.StringToInt(c.Param("roleId"))
	result, err := s.SetDeptLabel()
	if err != nil {
		log.Errorf("SetDeptLabel error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "")
	}
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = s.GetRoleDeptId(id)
		if err != nil {
			log.Errorf("抱歉未找到相关信息, %s", err.Error())
			e.Error(http.StatusInternalServerError, err, "")
		}
	}
	e.OK(gin.H{
		"depts":       result,
		"checkedKeys": menuIds,
	}, "")
}
