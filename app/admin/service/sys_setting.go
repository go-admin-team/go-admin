package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/common/log"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysSetting struct {
	service.Service
}

// GetSysSetting 获取SysSetting对象
func (e *SysSetting) GetSysSetting(model *models.SysSetting) error {
	var err error
	var data models.SysSetting
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		First(model)
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateSysSetting 修改SysSetting对象
func (e *SysSetting) UpdateSysSetting(c *models.SysSetting) error {
	var err error
	var data models.SysSetting
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Where(c.GetId()).Updates(c)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}
