package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/admin/models"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysLoginLog struct {
	api.Api
}

// GetPage 登录日志列表
// @Summary 登录日志列表
// @Description 获取JSON
// @Tags 登录日志
// @Param username query string false "用户名"
// @Param ipaddr query string false "ip地址"
// @Param loginLocation  query string false "归属地"
// @Param status query string false "状态"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-login-log [get]
// @Security Bearer
func (e SysLoginLog) GetPage(c *gin.Context) {
	s := service.SysLoginLog{}
	req :=dto.SysLoginLogGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	list := make([]models.SysLoginLog, 0)
	var count int64
	err = s.GetPage(&req, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 登录日志通过id获取
// @Summary 登录日志通过id获取
// @Description 获取JSON
// @Tags 登录日志
// @Param id path string false "id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-login-log/{id} [get]
// @Security Bearer
func (e SysLoginLog) Get(c *gin.Context) {
	s := service.SysLoginLog{}
	req :=dto.SysLoginLogGetReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.SysLoginLog
	err = s.Get(&req, &object)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(object, "查询成功")
}

// Delete 登录日志删除
// @Summary 登录日志删除
// @Description 登录日志删除
// @Tags 登录日志
// @Param data body dto.SysLoginLogDeleteReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-login-log [delete]
// @Security Bearer
func (e SysLoginLog) Delete(c *gin.Context) {
	s := service.SysLoginLog{}
	req :=dto.SysLoginLogDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	err = s.Remove(&req)
	if err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(req.GetId(), "删除成功")
}