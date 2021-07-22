package service

import (
	"errors"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
)

type SysConfig struct {
	service.Service
}

// GetPage 获取SysConfig列表
func (e *SysConfig) GetPage(c *dto.SysConfigGetPageReq, list *[]models.SysConfig, count *int64) error {
	err := e.Orm.
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

// Get 获取SysConfig对象
func (e *SysConfig) Get(d *dto.SysConfigGetReq, model *models.SysConfig) error {
	err := e.Orm.First(model, d.GetId()).Error
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

// Insert 创建SysConfig对象
func (e *SysConfig) Insert(c *dto.SysConfigControl) error {
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

// Update 修改SysConfig对象
func (e *SysConfig) Update(c *dto.SysConfigControl) error {
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

// SetSysConfig 修改SysConfig对象
func (e *SysConfig) SetSysConfig(c *[]dto.GetSetSysConfigReq) error {
	var err error
	for _, req := range *c {
		var model = models.SysConfig{}
		e.Orm.Where("config_key = ?", req.ConfigKey).First(&model)
		if model.Id != 0 {
			req.Generate(&model)
			db := e.Orm.Save(&model)
			err = db.Error
			if err != nil {
				e.Log.Errorf("Service SetSysConfig error:%s", err)
				return err
			}
			if db.RowsAffected == 0 {
				return errors.New("无权更新该数据")
			}
		}
	}
	return nil
}

func (e *SysConfig) GetForSet(c *[]dto.GetSetSysConfigReq) error {
	var err error
	var data models.SysConfig

	err = e.Orm.Model(&data).
		Find(c).Error
	if err != nil {
		e.Log.Errorf("Service GetSysConfigPage error:%s", err)
		return err
	}
	return nil
}

func (e *SysConfig) UpdateForSet(c *[]dto.GetSetSysConfigReq) error {
	m := *c
	for _, req := range m {
		var data models.SysConfig
		if err := e.Orm.Where("config_key = ?", req.ConfigKey).
			First(&data).Error; err != nil {
			e.Log.Errorf("Service GetSysConfigPage error:%s", err)
			return err
		}
		if data.ConfigValue != req.ConfigValue {
			data.ConfigValue = req.ConfigValue

			if err := e.Orm.Save(&data).Error; err != nil {
				e.Log.Errorf("Service GetSysConfigPage error:%s", err)
				return err
			}
		}

	}

	return nil
}

// Remove 删除SysConfig
func (e *SysConfig) Remove(d *dto.SysConfigDeleteReq) error {
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

// GetWithKey 根据Key获取SysConfig
func (e *SysConfig) GetWithKey(c *dto.SysConfigByKeyReq, resp *dto.GetSysConfigByKEYForServiceResp) error {
	var err error
	var data models.SysConfig
	err = e.Orm.Table(data.TableName()).Where("config_key = ?", c.ConfigKey).First(resp).Error
	if err != nil {
		e.Log.Errorf("At Service GetSysConfigByKEY Error:%s", err)
		return err
	}

	return nil
}

func (e *SysConfig) GetWithKeyList(c *dto.SysConfigGetToSysAppReq, list *[]models.SysConfig) error {
	var err error
	err = e.Orm.
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
