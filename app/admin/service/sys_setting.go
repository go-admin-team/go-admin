package service

import (
	"errors"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/common/service"
)

type SysSetting struct {
	service.Service
}

// GetSysSetting 获取SysSetting对象
func (e *SysSetting) GetSysSetting(model *models.SysSetting) error {
	var err error
	var data models.SysSetting

	db := e.Orm.Model(&data).
		First(model)
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// UpdateSysSetting 修改SysSetting对象
func (e *SysSetting) UpdateSysSetting(c *models.SysSetting) error {
	var err error
	var data models.SysSetting

	db := e.Orm.Model(&data).
		Where(c.GetId()).Updates(c)
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}
