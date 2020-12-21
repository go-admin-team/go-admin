package sys_file

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/service"
	"net/http"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"
)

type SysFileDir struct {
	apis.Api
}

func (e *SysFileDir) GetSysFileDirList(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	search := new(dto.SysFileDirSearch)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	err = c.ShouldBind(search)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
	}

	var list *[]models.SysFileDirL
	serviceStudent := service.SysFileDir{}
	serviceStudent.MsgID = msgID
	serviceStudent.Orm = db
	list, err = serviceStudent.SetSysFileDir(search)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, list, "查询成功")
}

func (e *SysFileDir) GetSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//查看详情
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
	}

	var object models.SysFileDir

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.MsgID = msgID
	serviceSysFileDir.Orm = db
	err = serviceSysFileDir.GetSysFileDir(control, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *SysFileDir) InsertSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//新增操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}
	// 设置创建人
	control.CreateBy = tools.GetUserId(c)

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Orm = db
	serviceSysFileDir.MsgID = msgID
	err = serviceSysFileDir.InsertSysFileDir(control)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, control.ID, "创建成功")
}

func (e *SysFileDir) UpdateSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}
	// 设置创建人
	control.UpdateBy = tools.GetUserId(c)

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Orm = db
	serviceSysFileDir.MsgID = msgID
	err = serviceSysFileDir.UpdateSysFileDir(control, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.ID, "更新成功")
}

func (e *SysFileDir) DeleteSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//删除操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}

	// 设置编辑人
	control.UpdateBy = tools.GetUserId(c)

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Orm = db
	serviceSysFileDir.MsgID = msgID
	err = serviceSysFileDir.RemoveSysFileDir(control, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.Id, "删除成功")
}
