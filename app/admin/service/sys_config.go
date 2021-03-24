package service

import (
	"errors"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysConfig struct {
	service.Service
}

// GetSysConfigPage 获取SysConfig列表
func (e *SysConfig) GetSysConfigPage(c *dto.SysConfigSearch, list *[]system.SysConfig, count *int64) error {
	var err error
	var data system.SysConfig

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *SysConfig) GetSysConfigByKey(c *dto.SysConfigSearch, list *[]system.SysConfig) error {
	var err error
	var data system.SysConfig

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetSysConfig 获取SysConfig对象
func (e *SysConfig) GetSysConfig(d *dto.SysConfigById, model *system.SysConfig) error {
	var err error
	var data system.SysConfig

	db := e.Orm.Model(&data).
		First(model, d.GetId())
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

// InsertSysConfig 创建SysConfig对象
func (e *SysConfig) InsertSysConfig(model *system.SysConfig) error {
	var err error
	var data system.SysConfig

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// UpdateSysConfig 修改SysConfig对象
func (e *SysConfig) UpdateSysConfig(c *system.SysConfig) error {
	var err error
	var data system.SysConfig

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

// RemoveSysConfig 删除SysConfig
func (e *SysConfig) RemoveSysConfig(d *dto.SysConfigById, c *system.SysConfig) error {
	var err error
	var data system.SysConfig

	db := e.Orm.Model(&data).
		Where(d.GetId()).Delete(c)
	if db.Error != nil {
		err = db.Error
		e.Log.Errorf("Delete error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	return nil
}

// GetSysConfigByKEY 根据Key获取SysConfig
func (e *SysConfig) GetSysConfigByKEY(c *dto.SysConfigControl) error {
	var err error
	var data system.SysConfig
	data.ConfigKey = c.ConfigKey
	err = e.Orm.Table(data.TableName()).Where("config_key = ?", data.ConfigKey).First(c).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	return nil
}
