package service

import (
	"errors"
	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysRole struct {
	service.Service
}

// GetSysRolePage 获取SysRole列表
func (e *SysRole) GetSysRolePage(c *dto.SysRoleSearch, list *[]system.SysRole, count *int64) error {
	var err error
	var data system.SysRole
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

// GetSysRole 获取SysRole对象
func (e *SysRole) GetSysRole(d *dto.SysRoleById, model *system.SysRole) error {
	var err error
	var data system.SysRole
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

// InsertSysRole 创建SysRole对象
func (e *SysRole) InsertSysRole(model *system.SysRole) error {
	var err error
	var data system.SysRole
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateSysRole 修改SysRole对象
func (e *SysRole) UpdateSysRole(c *system.SysRole) error {
	var err error
	var data system.SysRole
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

// RemoveSysRole 删除SysRole
func (e *SysRole) RemoveSysRole(d *dto.SysRoleById) error {
	var err error
	var data system.SysRole
	msgID := e.MsgID

	db := e.Orm.Model(&data).Delete(&data, d.GetId())
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
