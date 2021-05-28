package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models"
)

// SysConfigControl 增、改使用的结构体
type SysSettingControl struct {
	SettingsId int    `json:"settings_id" binding:"required"` // 头像
	Name       string `json:"name" binding:"required"`        // 名称
	Logo       string `json:"logo" binding:"required"`        // 头像
}

// Bind 映射上下文中的结构体数据
func (s *SysSettingControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)

	err := ctx.ShouldBind(s)
	if err != nil {
		log.Errorf("ShouldBind error: %s", err.Error())
		return err
	}

	err = ctx.ShouldBindUri(s)
	if err != nil {
		log.Errorf("ShouldBindUri error: %s", err.Error())
		return err
	}
	var jsonStr []byte
	jsonStr, err = json.Marshal(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	ctx.Set("body", string(jsonStr))
	return err
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *SysSettingControl) Generate() (*models.SysSetting, error) {
	return &models.SysSetting{
		SettingsId: s.SettingsId,
		Name:       s.Name,
		Logo:       s.Logo,
	}, nil
}

// GetId 获取数据对应的ID
func (s *SysSettingControl) GetId() interface{} {
	return s.SettingsId
}

// SysConfigById 获取单个或者删除的结构体
type SysSettingById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysSettingById) Generate() *SysSettingById {
	cp := *s
	return &cp
}

func (s *SysSettingById) GetId() interface{} {
	return s.Id
}

func (s *SysSettingById) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysSettingById) GenerateM() (*models.SysSetting, error) {
	return &models.SysSetting{}, nil
}
