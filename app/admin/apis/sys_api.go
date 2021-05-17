package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go-admin/common/apis"
)

type SysApi struct {
	apis.Api
}

// GetSysApiList 获取接口管理列表
// @Summary 获取接口管理列表
// @Description 获取接口管理列表
// @Tags 接口管理
// @Param name query string false "名称"
// @Param title query string false "标题"
// @Param path query string false "地址"
// @Param action query string false "类型"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysApi}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys_api [get]
// @Security Bearer
func (e SysApi) GetSysApiList(c *gin.Context) {
	e.MakeContext(c)
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	d := new(dto.SysApiSearch)
	//查询列表
	err = d.Bind(c)
	if err != nil {
		log.Warnf("request body bind error, %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	list := make([]models.SysApi, 0)
	var count int64
	serviceSysApi := service.SysApi{}
	serviceSysApi.Log = log
	serviceSysApi.Orm = db
	err = serviceSysApi.GetSysApiPage(d, p, &list, &count)
	if err != nil {
		log.Errorf("Get SysApi Page error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

// GetSysApi 获取接口管理
// @Summary 获取接口管理
// @Description 获取接口管理
// @Tags 接口管理
// @Param id path string false "id"
// @Success 200 {object} response.Response{data=models.SysApi} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys_api/{id} [get]
// @Security Bearer
func (e SysApi) GetSysApi(c *gin.Context) {
	control := new(dto.SysApiById)
	s := new(service.SysApi)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	var object models.SysApi

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	err = s.GetSysApi(control, p, &object).Error
	if err != nil {
		e.Logger.Errorf("Get SysApi error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

// UpdateSysApi 修改接口管理
// @Summary 修改接口管理
// @Description 修改接口管理
// @Tags 接口管理
// @Accept application/json
// @Product application/json
// @Param data body dto.SysApiControl true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys_api/{id} [put]
// @Security Bearer
func (e SysApi) UpdateSysApi(c *gin.Context) {
	e.MakeContext(c)
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	control := new(dto.SysApiControl)

	//更新操作
	err = control.Bind(c)
	if err != nil {
		log.Warnf("request body bind error, %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.GenerateM()
	if err != nil {
		log.Errorf("generate SysApi model error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	//数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysApi := service.SysApi{}
	serviceSysApi.Orm = db
	serviceSysApi.Log = log
	err = serviceSysApi.UpdateSysApi(object.(*models.SysApi), p)
	if err != nil {
		log.Errorf("Update SysApi error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "更新失败")
		return
	}
	e.OK(object.GetId(), "更新成功")
}

// DeleteSysApi 删除接口管理
// @Summary 删除接口管理
// @Description 删除接口管理
// @Tags 接口管理
// @Param ids body []int false "ids"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys_api [delete]
// @Security Bearer
func (e SysApi) DeleteSysApi(c *gin.Context) {
	e.MakeContext(c)
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}

	control := new(dto.SysApiById)

	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Warnf("request body bind error, %s", err.Error())
		e.Error(http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	// 设置编辑人
	control.SetUpdateBy(user.GetUserId(c))

	// 数据权限检查
	p := actions.GetPermissionFromContext(c)

	serviceSysApi := service.SysApi{}
	serviceSysApi.Orm = db
	serviceSysApi.Log = log
	err = serviceSysApi.RemoveSysApi(control, p)
	if err != nil {
		log.Errorf("Remove SysApi error, %s", err.Error())
		e.Error(http.StatusInternalServerError, err, "删除失败")
		return
	}
	e.OK(control.GetId(), "删除成功")
}