package apis

import (
	"go-admin/app/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysLoginLog struct {
	api.Api
}

func (e SysLoginLog) GetSysLoginLogList(c *gin.Context) {
	s := new(service.SysLoginLog)
	d := new(dto.SysLoginLogSearch)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(d).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	list := make([]models.SysLoginLog, 0)
	var count int64

	err = s.GetSysLoginLogPage(d, &list, &count)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e SysLoginLog) GetSysLoginLog(c *gin.Context) {
	s := new(service.SysLoginLog)
	control := new(dto.SysLoginLogById)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	var object models.SysLoginLog
	err = s.GetSysLoginLog(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysLoginLog) DeleteSysLoginLog(c *gin.Context) {
	s := service.SysLoginLog{}
	control := new(dto.SysLoginLogById)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		e.Logger.Error(err)
		return
	}

	// 设置编辑人
	control.SetUpdateBy(user.GetUserId(c))

	err = s.RemoveSysLoginLog(control)
	if err != nil {
		e.Logger.Errorf("RemoveSysLoginLog error, %s", err)
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(control.GetId(), "删除成功")
}