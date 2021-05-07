package sys_login_log

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
)

type SysLoginLog struct {
	apis.Api
}

func (e SysLoginLog) GetSysLoginLogList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysLoginLogSearch)
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

	list := make([]system.SysLoginLog, 0)
	var count int64
	serviceStudent := service.SysLoginLog{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysLoginLogPage(d, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e SysLoginLog) GetSysLoginLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysLoginLogById)
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
	var object system.SysLoginLog

	serviceSysLoginLog := service.SysLoginLog{}
	serviceSysLoginLog.Log = log
	serviceSysLoginLog.Orm = db
	err = serviceSysLoginLog.GetSysLoginLog(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysLoginLog) InsertSysLoginLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysLoginLogControl)
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

	serviceSysLoginLog := service.SysLoginLog{}
	serviceSysLoginLog.Orm = db
	serviceSysLoginLog.Log = log
	err = serviceSysLoginLog.InsertSysLoginLog(object)
	if err != nil {
		log.Errorf("InsertSysLoginLog error, %s", err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(object.GetId(), "创建成功")
}

func (e SysLoginLog) UpdateSysLoginLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysLoginLogControl)
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

	serviceSysLoginLog := service.SysLoginLog{}
	serviceSysLoginLog.Orm = db
	serviceSysLoginLog.Log = log
	err = serviceSysLoginLog.UpdateSysLoginLog(object)
	if err != nil {
		log.Errorf("UpdateSysLoginLog error, %s", err)
		e.Error(http.StatusInternalServerError, err, "更新失败")
		return
	}
	e.OK(object.GetId(), "更新成功")
}

func (e SysLoginLog) DeleteSysLoginLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysLoginLogById)
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
	object, err := control.GenerateM()
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "模型生成失败")
		return
	}

	// 设置编辑人
	object.SetUpdateBy(user.GetUserId(c))

	serviceSysLoginLog := service.SysLoginLog{}
	serviceSysLoginLog.Orm = db
	serviceSysLoginLog.Log = log
	err = serviceSysLoginLog.RemoveSysLoginLog(control, object)
	if err != nil {
		log.Errorf("RemoveSysLoginLog error, %s", err)
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(object.GetId(), "删除成功")
}
