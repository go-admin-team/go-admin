package sys_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
)

type SysFileInfo struct {
	apis.Api
}

func (e SysFileInfo) GetSysFileInfoList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	search := new(dto.SysFileInfoSearch)
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}
	err = c.ShouldBind(search)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysFileInfo, 0)
	var count int64
	serviceStudent := service.SysFileInfo{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysFileInfoPage(search, p, &list, &count)
	if err != nil {
		log.Errorf("GetSysFileInfoPage error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), search.PageIndex, search.PageSize, "查询成功")
}

func (e SysFileInfo) GetSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoById)
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
		log.Warnf("参数验证错误, error:%s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	var object models.SysFileInfo

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Log = log
	serviceSysFileInfo.Orm = db
	err = serviceSysFileInfo.GetSysFileInfo(control, p, &object)
	if err != nil {
		log.Errorf("GetSysFileInfo error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysFileInfo) InsertSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoControl)
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
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	// 设置创建人
	control.CreateBy = user.GetUserId(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.Log = log
	err = serviceSysFileInfo.InsertSysFileInfo(control)
	if err != nil {
		log.Errorf("InsertSysFileInfo error: %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(control.ID, "创建成功")
}

func (e SysFileInfo) UpdateSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoControl)
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	err = c.ShouldBindUri(control)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Warnf("参数验证错误, error:%s", err)
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	// 设置创建人
	control.UpdateBy = user.GetUserId(c)

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.Log = log
	err = serviceSysFileInfo.UpdateSysFileInfo(control, p)
	if err != nil {
		log.Errorf("UpdateSysFileInfo error: %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}
	e.OK(control.ID, "更新成功")
}

func (e SysFileInfo) DeleteSysFileInfo(c *gin.Context) {
	control := new(dto.SysFileInfoById)
	e.Context = c
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	err = c.ShouldBindUri(control)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(422, err, "参数验证失败")
		return
	}
	err = c.ShouldBind(control)
	if err != nil {
		log.Warnf("参数验证错误, error: %s", err)
		e.Error(422, err, "参数验证失败")
		return
	}

	// 设置编辑人
	control.UpdateBy = user.GetUserId(c)

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.Log = log
	err = serviceSysFileInfo.RemoveSysFileInfo(control, p)
	if err != nil {
		log.Errorf("RemoveSysFileInfo error: %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}
	e.OK(control.Id, "删除成功")
}
