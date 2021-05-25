package apis

import (
	"github.com/gin-gonic/gin/binding"
	"go-admin/app/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysOperaLog struct {
	api.Api
}

// GetSysOperaLogList 操作日志列表
func (e SysOperaLog) GetSysOperaLogList(c *gin.Context) {
	s := service.SysOperaLog{}
	req := new(dto.SysOperaLogSearch)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	list := make([]models.SysOperaLog, 0)
	var count int64

	err = s.GetSysOperaLogPage(req, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetSysOperaLog  操作日志通过id获取
func (e SysOperaLog) GetSysOperaLog(c *gin.Context) {
	s := new(service.SysOperaLog)
	req := new(dto.SysOperaLogById)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	var object models.SysOperaLog

	err = s.GetSysOperaLog(req, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

// DeleteSysOperaLog 操作日志删除
func (e SysOperaLog) DeleteSysOperaLog(c *gin.Context) {
	s := new(service.SysOperaLog)
	req := new(dto.SysOperaLogById)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	err = s.RemoveSysOperaLog(req)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	e.OK(req.GetId(), "删除成功")
}