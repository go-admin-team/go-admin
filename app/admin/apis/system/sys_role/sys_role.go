package sys_role

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
	"go-admin/common/global"
)

type SysRole struct {
	apis.Api
}

// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role [get]
// @Security Bearer
func (e SysRole) GetSysRoleList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysRoleSearch)
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

	list := make([]system.SysRole, 0)
	var count int64
	s := service.SysRole{}
	s.Log = log
	s.Orm = db
	err = s.GetSysRolePage(d, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role/{id} [get]
// @Security Bearer
func (e SysRole) GetSysRole(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysRoleById)
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
	var object system.SysRole

	s := service.SysRole{}
	s.Log = log
	s.Orm = db
	err = s.GetSysRole(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysRoleControl true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
// @Security Bearer
func (e SysRole) InsertSysRole(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysRoleControl)
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
	object.CreateBy = user.GetUserId(c)
	if object.Status == "" {
		object.Status = "2"
	}

	s := service.SysRole{}
	s.Orm = db
	s.Log = log
	err = s.InsertSysRole(object)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(object.GetId(), "创建成功")
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysRoleControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role/{id} [put]
// @Security Bearer
func (e SysRole) UpdateSysRole(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysRoleControl)
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
	object.UpdateBy = user.GetUserId(c)

	s := service.SysRole{}
	s.Orm = db
	s.Log = log
	err = s.UpdateSysRole(object)
	if err != nil {
		log.Error(err)
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(object.GetId(), "更新成功")
}

// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param data body dto.SysRoleById true "body"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role [delete]
// @Security Bearer
func (e SysRole) DeleteSysRole(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysRoleById)
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

	s := service.SysRole{}
	s.Orm = db
	s.Log = log
	err = s.RemoveSysRole(control)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(control.GetId(), "删除成功")
}

func (e SysRole) UpdateRoleDataScope(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.RoleDataScopeReq)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//更新操作
	err = c.Bind(control)
	if err != nil {
		log.Errorf("request bind error, %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	data := &system.SysRole{
		RoleId:    control.RoleId,
		DataScope: control.DataScope,
		DeptIds:   control.DeptIds,
	}
	data.UpdateBy = user.GetUserId(c)
	s := &service.SysRole{}
	s.Orm = db
	s.Log = log
	err = s.UpdateDataScope(data)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(nil, "操作成功")
}
