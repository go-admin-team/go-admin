package sys_file

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/service"
	"go-admin/tools/app"
	"net/http"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"
)

type SysFileInfo struct {
	apis.Api
}

func (e *SysFileInfo) GetSysFileInfoList(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	search := new(dto.SysFileInfoSearch)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}
	err = c.ShouldBind(search)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysFileInfo, 0)
	var count int64
	serviceStudent := service.SysFileInfo{}
	serviceStudent.MsgID = msgID
	serviceStudent.Orm = db
	err = serviceStudent.GetSysFileInfoPage(search, p, &list, &count)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), search.PageIndex, search.PageSize, "查询成功")
}

func (e *SysFileInfo) GetSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//查看详情
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}

	var object models.SysFileInfo

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.MsgID = msgID
	serviceSysFileInfo.Orm = db
	err = serviceSysFileInfo.GetSysFileInfo(control, p, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *SysFileInfo) InsertSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//新增操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	// 设置创建人
	control.CreateBy = tools.GetUserId(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.MsgID = msgID
	err = serviceSysFileInfo.InsertSysFileInfo(control)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, control.ID, "创建成功")
}

func (e *SysFileInfo) UpdateSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoControl)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	// 设置创建人
	control.UpdateBy = tools.GetUserId(c)

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.MsgID = msgID
	err = serviceSysFileInfo.UpdateSysFileInfo(control, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.ID, "更新成功")
}

func (e *SysFileInfo) DeleteSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoById)
	db, err := tools.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	msgID := tools.GenerateMsgIDFromContext(c)
	//删除操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}

	// 设置编辑人
	control.UpdateBy = tools.GetUserId(c)

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.MsgID = msgID
	err = serviceSysFileInfo.RemoveSysFileInfo(control, p)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.Id, "删除成功")
}
