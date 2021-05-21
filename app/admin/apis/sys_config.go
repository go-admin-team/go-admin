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
	var control = new(dto.SysConfigControl)
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