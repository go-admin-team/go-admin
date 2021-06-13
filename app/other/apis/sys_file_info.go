package apis

import (
	models2 "go-admin/app/other/models"
	service2 "go-admin/app/other/service"
	dto2 "go-admin/app/other/service/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/common/actions"
)

type SysFileInfo struct {
	api.Api
}

func (e SysFileInfo) GetPage(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	search := new(dto2.SysFileInfoSearch)
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

	list := make([]models2.SysFileInfo, 0)
	var count int64
	serviceStudent := service2.SysFileInfo{}
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

func (e SysFileInfo) Get(c *gin.Context) {
	control := new(dto2.SysFileInfoById)
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

	var object models2.SysFileInfo

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysFileInfo := service2.SysFileInfo{}
	serviceSysFileInfo.Log = log
	serviceSysFileInfo.Orm = db
	err = serviceSysFileInfo.GetSysFileInfo(control, p, &object)
	if err != nil {
		log.Errorf("GetSysFileInfo error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.OK(object, "查询成功")
}

func (e SysFileInfo) Insert(c *gin.Context) {
	control := new(dto2.SysFileInfoControl)
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

	serviceSysFileInfo := service2.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.Log = log
	err = serviceSysFileInfo.InsertSysFileInfo(control)
	if err != nil {
		log.Errorf("InsertSysFileInfo error: %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(control.Id, "创建成功")
}

func (e SysFileInfo) Update(c *gin.Context) {
	control := new(dto2.SysFileInfoControl)
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

	serviceSysFileInfo := service2.SysFileInfo{}
	serviceSysFileInfo.Orm = db
	serviceSysFileInfo.Log = log
	err = serviceSysFileInfo.UpdateSysFileInfo(control, p)
	if err != nil {
		log.Errorf("UpdateSysFileInfo error: %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}
	e.OK(control.Id, "更新成功")
}

func (e SysFileInfo) Delete(c *gin.Context) {
	control := new(dto2.SysFileInfoById)
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

	serviceSysFileInfo := service2.SysFileInfo{}
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
