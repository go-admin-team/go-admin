package service

import (
	"errors"

	"gorm.io/gorm"

	"go-admin/app/admin/models/system"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/service"
)

type SysPost struct {
	service.Service
}

// GetSysPostPage 获取SysPost列表
func (e *SysPost) GetSysPostPage(c *dto.SysPostSearch, list *[]system.SysPost, count *int64) error {
	var err error
	var data system.SysPost

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

// GetSysPost 获取SysPost对象
func (e *SysPost) GetSysPost(d *dto.SysPostById, model *system.SysPost) error {
	var err error
	var data system.SysPost

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

// InsertSysPost 创建SysPost对象
func (e *SysPost) InsertSysPost(model *system.SysPost) error {
	var err error
	var data system.SysPost

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// UpdateSysPost 修改SysPost对象
func (e *SysPost) UpdateSysPost(c *system.SysPost) error {
	var err error
	var data system.SysPost

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

// RemoveSysPost 删除SysPost
func (e *SysPost) RemoveSysPost(d *dto.SysPostById) error {
	var err error
	var data system.SysPost

	db := e.Orm.Model(&data).Delete(&data, d.GetId())
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
