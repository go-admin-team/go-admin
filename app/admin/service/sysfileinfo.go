package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysFileInfo struct {
	service.Service
}

// GetSysFileInfoPage 获取SysFileInfo列表
func (e *SysFileInfo) GetSysFileInfoPage(c cDto.Index, p *actions.DataPermission, list *[]models.SysFileInfo, count *int64) error {
	var err error
	var data models.SysFileInfo
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetSysFileInfo 获取SysFileInfo对象
func (e *SysFileInfo) GetSysFileInfo(d cDto.Control, p *actions.DataPermission, model *models.SysFileInfo) error {
	var err error
	var data models.SysFileInfo
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
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

// InsertSysFileInfo 创建SysFileInfo对象
func (e *SysFileInfo) InsertSysFileInfo(model common.ActiveRecord) error {
	var err error
	var data models.SysFileInfo
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateSysFileInfo 修改SysFileInfo对象
func (e *SysFileInfo) UpdateSysFileInfo(c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data models.SysFileInfo
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(c.GetId()).Updates(c)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysFileInfo 删除SysFileInfo
func (e *SysFileInfo) RemoveSysFileInfo(d cDto.Control, c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data models.SysFileInfo
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(d.GetId()).Delete(c)
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