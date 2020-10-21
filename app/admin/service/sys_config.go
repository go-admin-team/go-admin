package service

import (
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service/dto"
	"go-admin/common/log"
	"go-admin/common/service"
)

type SysConfig struct {
	service.Service
}

// GetSysConfigByKEY 根据Key获取SysConfig
func (e *SysConfig) GetSysConfigByKEY(c *dto.SysConfigControl) error {
	var err error
	var data system.SysConfig
	msgID := e.MsgID
	data.ConfigKey = c.ConfigKey
	err = e.Orm.Table(data.TableName()).Where("config_key = ?", data.ConfigKey).First(c).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}

	return nil
}
