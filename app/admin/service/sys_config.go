package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"gorm.io/gorm"
)

type SysConfig struct {
	service.Service
}

// GetSysConfigPage 获取SysConfig列表
func (e *SysConfig) GetSysConfigPage(c *dto.SysConfigSearch, list *[]models.SysConfig, count *int64) error {
	var err error
	var data models.SysConfig

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("Service GetSysConfigPage error:%s", err)
		return err
	}
	return nil
}

func (e *SysConfig) GetSysConfigByKey(c *dto.SysConfigSearch, list *[]models.SysConfig) error {
	var err error
	var data models.SysConfig

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("Service GetSysConfigByKey error:%s", err)
		return err
	}
	return nil
}

// GetSysConfig 获取SysConfig对象
func (e *SysConfig) GetSysConfig(d *dto.SysConfigById, model *models.SysConfig) error {
	var data models.SysConfig

	err := e.Orm.Model(&data).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysConfigPage error:%s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("Service GetSysConfig error:%s", err)
		return err
	}
	return nil
}

// InsertSysConfig 创建SysConfig对象
func (e *SysConfig) InsertSysConfig(c *dto.SysConfigControl) error {
	var err error
	var data models.SysConfig
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("Service InsertSysConfig error:%s", err)
		return err
	}
	return nil
}

// UpdateSysConfig 修改SysConfig对象
func (e *SysConfig) UpdateSysConfig(c *dto.SysConfigControl) error {
	var err error
	var model = models.SysConfig{}
	e.Orm.First(&model, c.GetId())
	c.Generate(&model)
	db := e.Orm.Save(&model)
	err = db.Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysConfig error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysConfig 删除SysConfig
func (e *SysConfig) RemoveSysConfig(d *dto.SysConfigById) error {
	var err error
	var data models.SysConfig

	db := e.Orm.Delete(&data, d.Ids)
	if db.Error != nil {
		err = db.Error
		e.Log.Errorf("Service RemoveSysConfig error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	return nil
}

// GetSysConfigByKEY 根据Key获取SysConfig
func (e *SysConfig) GetSysConfigByKEY(c *dto.SysConfigByKeyReq) error {
	var err error
	var data models.SysConfig
	data.ConfigKey = c.ConfigKey
	err = e.Orm.Table(data.TableName()).Where("config_key = ?", data.ConfigKey).First(c).Error
	if err != nil {
		e.Log.Errorf("Service GetSysConfigByKEY error:%s", err)
		return err
	}

	return nil
}