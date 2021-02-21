package sys_role

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
	"go-admin/common/global"
	"go-admin/common/log"
	"go-admin/tools"
	"net/http"
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
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role [get]
// @Security Bearer
func (e *SysRole) GetSysRoleList(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	d := new(dto.SysRoleSearch)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	list := make([]system.SysRole, 0)
	var count int64
	serviceStudent := service.SysRole{}
	serviceStudent.MsgID = msgID
	serviceStudent.Orm = db
	err = serviceStudent.GetSysRolePage(d, &list, &count)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role/{id} [get]
// @Security Bearer
func (e *SysRole) GetSysRole(c *gin.Context) {
	control := new(dto.SysRoleById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object system.SysRole

	serviceSysOperlog := service.SysRole{}
	serviceSysOperlog.MsgID = msgID
	serviceSysOperlog.Orm = db
	err = serviceSysOperlog.GetSysRole(control, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
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
// @Security Bearer
func (e *SysRole) InsertSysRole(c *gin.Context) {
	control := new(dto.SysRoleControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//新增操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.CreateBy = tools.GetUserId(c)

	s := service.SysRole{}
	s.Orm = db
	s.MsgID = msgID
	err = s.InsertSysRole(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "")
		return
	}
	e.OK(c, object.GetId(), "创建成功")
}

// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body models.SysRole true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role/{id} [put]
// @Security Bearer
func (e *SysRole) UpdateSysRole(c *gin.Context) {
	control := new(dto.SysRoleControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//更新操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.UpdateBy = tools.GetUserId(c)

	s := service.SysRole{}
	s.Orm = db
	s.MsgID = msgID
	err = s.UpdateSysRole(object)
	if err != nil {
		log.Error(err)
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "")
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param roleId path int true "roleId"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role/{roleId} [delete]
// @Security Bearer
func (e *SysRole) DeleteSysRole(c *gin.Context) {
	control := new(dto.SysRoleById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("MsgID[%s] Bind error: %s", msgID, err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	s := service.SysRole{}
	s.Orm = db
	s.MsgID = msgID
	err = s.RemoveSysRole(control)
	if err != nil {
		log.Error(err)
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "")
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}

func (e *SysRole) UpdateRoleDataScope(c *gin.Context) {
	control := new(dto.RoleDataScopeReq)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//更新操作
	err = c.Bind(control)
	if err != nil {
		log.Errorf("msgID[%s] request bind error, %s", msgID, err.Error())
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	data := &system.SysRole{
		RoleId:    control.RoleId,
		DataScope: control.DataScope,
		DeptIds:   control.DeptIds,
	}
	data.UpdateBy = tools.GetUserId(c)
	s := &service.SysRole{}
	s.Orm = db
	s.MsgID = msgID
	err = s.UpdateDataScope(data)
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "")
		return
	}
	e.OK(c, nil, "操作成功")
}
