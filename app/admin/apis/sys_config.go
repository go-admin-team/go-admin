package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
)

type SysConfig struct {
	api.Api
}

// GetSysConfigList 获取配置管理列表
// @Summary 获取配置管理列表
// @Description 获取配置管理列表
// @Tags 配置管理
// @Param configName query string false "名称"
// @Param configKey query string false "key"
// @Param configType query string false "类型"
// @Param isFrontend query int false "是否前端"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysApi}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-config [get]
// @Security Bearer
func (e SysConfig) GetSysConfigList(c *gin.Context) {
	s := service.SysConfig{}
	d := new(dto.SysConfigSearch)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(d, binding.Query).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	list := make([]models.SysConfig, 0)
	var count int64
	err = s.GetSysConfigPage(d, &list, &count)
	if err != nil {
		e.Logger.Errorf("GetSysConfigPage 查询失败, error:%s", err)
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

// GetSysConfigBySysApp 获取系统配置信息，主要注意这里不在验证数据权限
func (e SysConfig) GetSysConfigBySysApp(c *gin.Context) {
	d := new(dto.SysConfigSearch)
	s := service.SysConfig{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(d, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	// 控制只读前台的数据
	d.IsFrontend = 1
	list := make([]models.SysConfig, 0)
	err = s.GetSysConfigByKey(d, &list)
	if err != nil {
		e.Logger.Errorf("GetSysConfigPage 查询失败, error:%s", err)
		e.Error(500, err, "查询失败")
		return
	}
	mp := make(map[string]string)
	for i := 0; i < len(list); i++ {
		key := list[i].ConfigKey
		if key != "" {
			mp[key] = list[i].ConfigValue
		}
	}
	e.OK(mp, "查询成功")
}

// GetSysConfig 获取配置管理
// @Summary 获取配置管理
// @Description 获取配置管理
// @Tags 配置管理
// @Param id path string false "id"
// @Success 200 {object} response.Response{data=models.SysApi} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys_api/{id} [get]
// @Security Bearer
func (e SysConfig) GetSysConfig(c *gin.Context) {
	control := new(dto.SysConfigById)
	s := service.SysConfig{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	var object models.SysConfig

	err = s.GetSysConfig(control, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		e.Logger.Errorf("Orm获取失败, error:%s", err)
		e.Error(500, err, "Orm获取失败")
		return
	}

	e.OK(object, "查看成功")
}

func (e SysConfig) InsertSysConfig(c *gin.Context) {
	s := service.SysConfig{}
	control := new(dto.SysConfigControl)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	// 设置创建人
	control.SetCreateBy(user.GetUserId(c))

	err = s.InsertSysConfig(control)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, "创建失败")
		return
	}
	e.OK(control.GetId(), "创建成功")
}

// UpdateSysConfig 修改配置管理
// @Summary 修改配置管理
// @Description 修改配置管理
// @Tags 配置管理
// @Accept application/json
// @Product application/json
// @Param data body dto.SysConfigControl true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-config/{id} [put]
// @Security Bearer
func (e SysConfig) UpdateSysConfig(c *gin.Context) {
	s := service.SysConfig{}
	control := new(dto.SysConfigControl)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	control.SetUpdateBy(user.GetUserId(c))

	err = s.UpdateSysConfig(control)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "更新失败")
		e.Logger.Errorf("Orm获取失败, error:%s", err)
		return
	}
	e.OK(control.GetId(), "更新成功")
}

// SetSysConfig 设置配置
// @Summary 设置配置
// @Description 界面操作设置配置值
// @Tags 配置管理
// @Accept application/json
// @Product application/json
// @Param data body []dto.SysConfigSetReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/set-config [put]
// @Security Bearer
func (e SysConfig) SetSysConfig(c *gin.Context) {
	s := service.SysConfig{}
	req := make([]dto.SysConfigSetReq, 0)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	err = s.SetSysConfig(&req)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "更新失败")
		e.Logger.Errorf("Orm获取失败, error:%s", err)
		return
	}
	e.OK("", "更新成功")
}

func (e SysConfig) GetSetSysConfig(c *gin.Context) {
	s := service.SysConfig{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	d := make([]dto.SysConfigSetReq, 0)
	err = s.GetSetSysConfig(&d)
	if err != nil {
		e.Logger.Errorf("GetSetSysConfig 查询失败, error:%s", err)
		e.Error(500, err, "查询失败")
		return
	}
	m := make(map[string]interface{}, 0)
	for _, req := range d {

		if req.ConfigKey == "sys_app_logo" {
			mp := make(map[string]interface{}, 0)
			mp["url"] = req.ConfigValue
			l := make([]map[string]interface{}, 0)
			l = append(l, mp)
			m[req.ConfigKey] = l
		} else {
			m[req.ConfigKey] = req.ConfigValue
		}

	}
	e.OK(m, "查询成功")
}

// DeleteSysConfig 删除配置管理
// @Summary 删除配置管理
// @Description 删除配置管理
// @Tags 配置管理
// @Param ids body []int false "ids"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-config [delete]
// @Security Bearer
func (e SysConfig) DeleteSysConfig(c *gin.Context) {
	s := service.SysConfig{}
	control := new(dto.SysConfigById)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	// 设置编辑人
	control.SetUpdateBy(user.GetUserId(c))

	err = s.RemoveSysConfig(control)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "删除失败")
		e.Logger.Errorf("Orm获取失败, error:%s", err)
		return
	}
	e.OK(control.GetId(), "删除成功")
}

// GetSysConfigByKEYForService 根据Key获取SysConfig的Service
func (e SysConfig) GetSysConfigByKEYForService(c *gin.Context) {
	var s = new(service.SysConfig)
	var control = new(dto.SysConfigByKeyReq)
	err := e.MakeContext(c).
		MakeOrm().
		Bind(control).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	err = s.GetSysConfigByKEY(control)
	if err != nil {
		e.Logger.Errorf("通过Key获取配置失败, error:%s", err)
		e.Error(500, err, "")
		return
	}
	e.OK(control, s.Msg)
}