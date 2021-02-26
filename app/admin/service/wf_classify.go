package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/app/admin/service/dto"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type WfProcessClassify struct {
	service.Service
}

// GetWfProcessClassifyPage 获取WfProcessClassify列表
func (e *WfProcessClassify) GetWfProcessClassifyPage(c *dto.WfProcessClassifySearch, p *actions.DataPermission, list *[]models.WfProcessClassify, count *int64) error {
	var err error
	var data models.WfProcessClassify
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

// GetWfProcessClassify 获取WfProcessClassify对象
func (e *WfProcessClassify) GetWfProcessClassify(d *dto.WfProcessClassifyById, p *actions.DataPermission, model *models.WfProcessClassify) error {
	var err error
	var data models.WfProcessClassify
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

// InsertWfProcessClassify 创建WfProcessClassify对象
func (e *WfProcessClassify) InsertWfProcessClassify(model *models.WfProcessClassify) error {
	var err error
	var data models.WfProcessClassify
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateWfProcessClassify 修改WfProcessClassify对象
func (e *WfProcessClassify) UpdateWfProcessClassify(c *models.WfProcessClassify, p *actions.DataPermission) error {
	var err error
	var data models.WfProcessClassify
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

// RemoveWfProcessClassify 删除WfProcessClassify
func (e *WfProcessClassify) RemoveWfProcessClassify(d *dto.WfProcessClassifyById, p *actions.DataPermission) error {
	var err error
	var data models.WfProcessClassify
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
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