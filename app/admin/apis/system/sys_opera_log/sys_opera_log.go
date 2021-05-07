package sys_opera_log

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
)

type SysOperaLog struct {
	apis.Api
}

func (e SysOperaLog) GetSysOperaLogList(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	d := new(dto.SysOperaLogSearch)
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

	list := make([]system.SysOperaLog, 0)
	var count int64
	serviceStudent := service.SysOperaLog{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysOperaLogPage(d, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e SysOperaLog) GetSysOperaLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysOperaLogById)
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
	var object system.SysOperaLog

	serviceSysOperlog := service.SysOperaLog{}
	serviceSysOperlog.Log = log
	serviceSysOperlog.Orm = db
	err = serviceSysOperlog.GetSysOperaLog(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysOperaLog) InsertSysOperaLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysOperaLogControl)
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

	serviceSysOperaLog := service.SysOperaLog{}
	serviceSysOperaLog.Orm = db
	serviceSysOperaLog.Log = log
	err = serviceSysOperaLog.InsertSysOperaLog(object)
	if err != nil {
		log.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(object.GetId(), "创建成功")
}

func (e SysOperaLog) UpdateSysOperaLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysOperaLogControl)
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

	serviceSysOperaLog := service.SysOperaLog{}
	serviceSysOperaLog.Orm = db
	serviceSysOperaLog.Log = log
	err = serviceSysOperaLog.UpdateSysOperaLog(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(object.GetId(), "更新成功")
}

func (e SysOperaLog) DeleteSysOperaLog(c *gin.Context) {
	e.Context = c
	log := e.GetLogger()
	control := new(dto.SysOperaLogById)
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

	serviceSysOperaLog := service.SysOperaLog{}
	serviceSysOperaLog.Orm = db
	serviceSysOperaLog.Log = log
	err = serviceSysOperaLog.RemoveSysOperaLog(control)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(control.GetId(), "删除成功")
}
