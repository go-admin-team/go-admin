package service

import (
	"errors"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
)

type SysOperaLog struct {
	service.Service
}

// GetPage 获取SysOperaLog列表
func (e *SysOperaLog) GetPage(c *dto.SysOperaLogGetPageReq, list *[]models.SysOperaLog, count *int64) error {
	var err error
	var data models.SysOperaLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("Service GetSysOperaLogPage error:%s", err.Error())
		return err
	}
	return nil
}

// Get 获取SysOperaLog对象
func (e *SysOperaLog) Get(d *dto.SysOperaLogGetReq, model *models.SysOperaLog) error {
	var data models.SysOperaLog

	err := e.Orm.Model(&data).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysOperaLog error:%s", err.Error())
		return err
	}
	if err != nil {
		e.Log.Errorf("Service GetSysOperaLog error:%s", err.Error())
		return err
	}
	return nil
}

// Insert 创建SysOperaLog对象
func (e *SysOperaLog) Insert(model *models.SysOperaLog) error {
	var err error
	var data models.SysOperaLog

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("Service InsertSysOperaLog error:%s", err.Error())
		return err
	}
	return nil
}

// Remove 删除SysOperaLog
func (e *SysOperaLog) Remove(d *dto.SysOperaLogDeleteReq) error {
	var err error
	var data models.SysOperaLog

	db := e.Orm.Model(&data).Delete(&data, d.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysOperaLog error:%s", err.Error())
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
