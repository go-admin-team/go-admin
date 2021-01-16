package sys_role

import (
	"github.com/gin-gonic/gin"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"

	"net/http"
)

type SysRole struct {
	apis.Api
}

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
	object.CreateBy= tools.GetUserId(c)

	serviceSysRole := service.SysRole{}
	serviceSysRole.Orm = db
	serviceSysRole.MsgID = msgID
	err = serviceSysRole.InsertSysRole(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

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

	serviceSysRole := service.SysRole{}
	serviceSysRole.Orm = db
	serviceSysRole.MsgID = msgID
	err = serviceSysRole.UpdateSysRole(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

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

	serviceSysRole := service.SysRole{}
	serviceSysRole.Orm = db
	serviceSysRole.MsgID = msgID
	err = serviceSysRole.RemoveSysRole(control)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}

