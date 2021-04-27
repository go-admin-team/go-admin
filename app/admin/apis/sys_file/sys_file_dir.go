package sys_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
)

type SysFileDir struct {
	apis.Api
}

func (e SysFileDir) GetSysFileDirList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	search := new(dto.SysFileDirSearch)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	err = c.ShouldBind(search)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}

	var list *[]models.SysFileDirL
	serviceStudent := service.SysFileDir{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	list, err = serviceStudent.SetSysFileDir(search)
	if err != nil {
		log.Errorf("SetSysFileDir error, %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(list, "查询成功")
}

func (e SysFileDir) GetSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirById)
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Warnf("ShouldBindUri error: %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
	}

	var object models.SysFileDir

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Log = log
	serviceSysFileDir.Orm = db
	err = serviceSysFileDir.GetSysFileDir(control, &object)
	if err != nil {
		log.Errorf("GetSysFileDir error, %s", err)
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysFileDir) InsertSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirControl)
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Warnf("ShouldBindUri error: %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Warnf("ShouldBind error: %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	// 设置创建人
	control.CreateBy = user.GetUserId(c)

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Orm = db
	serviceSysFileDir.Log = log
	err = serviceSysFileDir.InsertSysFileDir(control)
	if err != nil {
		log.Errorf("InsertSysFileDir error, %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(control.ID, "创建成功")
}

func (e SysFileDir) UpdateSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirControl)
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	err = c.ShouldBindUri(control)
	if err != nil {
		log.Warnf("ShouldBindUri error: %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Warnf("ShouldBind error: %#v", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
	}
	// 设置创建人
	control.UpdateBy = user.GetUserId(c)

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Orm = db
	serviceSysFileDir.Log = log
	err = serviceSysFileDir.UpdateSysFileDir(control, p)
	if err != nil {
		log.Errorf("UpdateSysFileDir error, %s", err)
		e.Error(http.StatusInternalServerError, err, "更新失败")
		return
	}
	e.OK(control.ID, "更新成功")
}

func (e SysFileDir) DeleteSysFileDir(c *gin.Context) {
	control := new(dto.SysFileDirById)
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	msgID := pkg.GenerateMsgIDFromContext(c)
	//删除操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Warnf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Warnf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
	}

	// 设置编辑人
	control.UpdateBy = user.GetUserId(c)

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileDir := service.SysFileDir{}
	serviceSysFileDir.Orm = db
	serviceSysFileDir.MsgID = msgID
	err = serviceSysFileDir.RemoveSysFileDir(control, p)
	if err != nil {
		log.Errorf("RemoveSysFileDir error, %s", err)
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(control.Id, "删除成功")
}
