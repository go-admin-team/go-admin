package service

import (
	"errors"
	"go-admin/app/admin/models/system"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysLoginlog struct {
	service.Service
}

// GetSysLoginlogPage 获取SysLoginlog列表
func (e *SysLoginlog) GetSysLoginlogPage(c cDto.Index, list *[]system.SysLoginLog, count *int64) error {
	var err error
	var data system.SysLoginLog
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetSysLoginlog 获取SysLoginlog对象
func (e *SysLoginlog) GetSysLoginlog(d cDto.Control, model *system.SysLoginLog) error {
	var err error
	var data system.SysLoginLog
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		First(model, d.GetId())
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

// InsertSysLoginlog 创建SysLoginlog对象
func (e *SysLoginlog) InsertSysLoginlog(model common.ActiveRecord) error {
	var err error
	var data system.SysLoginLog
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateSysLoginlog 修改SysLoginlog对象
func (e *SysLoginlog) UpdateSysLoginlog(c common.ActiveRecord) error {
	var err error
	var data system.SysLoginLog
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

// RemoveSysLoginlog 删除SysLoginlog
func (e *SysLoginlog) RemoveSysLoginlog(d cDto.Control, c common.ActiveRecord) error {
	var err error
	var data system.SysLoginLog
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Where(d.GetId()).Delete(c)
	if db.Error != nil {
		err = db.Error
		log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	return nil
}