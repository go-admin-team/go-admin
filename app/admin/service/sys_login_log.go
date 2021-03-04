package service

import (
	"errors"

	"gorm.io/gorm"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	common "go-admin/common/models"
	"go-admin/common/service"
)

type SysLoginLog struct {
	service.Service
}

// GetSysLoginLogPage 获取SysLoginLog列表
func (e *SysLoginLog) GetSysLoginLogPage(c *dto.SysLoginLogSearch, list *[]system.SysLoginLog, count *int64) error {
	var err error
	var data system.SysLoginLog

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

// GetSysLoginLog 获取SysLoginLog对象
func (e *SysLoginLog) GetSysLoginLog(d *dto.SysLoginLogById, model *system.SysLoginLog) error {
	var err error
	var data system.SysLoginLog

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

// InsertSysLoginLog 创建SysLoginLog对象
func (e *SysLoginLog) InsertSysLoginLog(model common.ActiveRecord) error {
	var err error
	var data system.SysLoginLog

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// UpdateSysLoginLog 修改SysLoginLog对象
func (e *SysLoginLog) UpdateSysLoginLog(c common.ActiveRecord) error {
	var err error
	var data system.SysLoginLog

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

// RemoveSysLoginLog 删除SysLoginLog
func (e *SysLoginLog) RemoveSysLoginLog(d *dto.SysLoginLogById, c common.ActiveRecord) error {
	var err error
	var data system.SysLoginLog

	db := e.Orm.Model(&data).
		Where(d.Ids).Delete(c)
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
