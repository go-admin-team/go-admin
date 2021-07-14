package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
)

type SysPost struct {
	service.Service
}

// GetPage 获取SysPost列表
func (e *SysPost) GetPage(c *dto.SysPostPageReq, list *[]models.SysPost, count *int64) error {
	var err error
	var data models.SysPost

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error:%s \r", err)
		return err
	}
	return nil
}

// Get 获取SysPost对象
func (e *SysPost) Get(d *dto.SysPostGetReq, model *models.SysPost) error {
	var err error
	var data models.SysPost

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

// Insert 创建SysPost对象
func (e *SysPost) Insert(c *dto.SysPostInsertReq) error {
	var err error
	var data models.SysPost
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysPost对象
func (e *SysPost) Update(c *dto.SysPostUpdateReq) error {
	var err error
	var model = models.SysPost{}
	e.Orm.First(&model, c.GetId())
	c.Generate(&model)

	db := e.Orm.Save(&model)
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// Remove 删除SysPost
func (e *SysPost) Remove(d *dto.SysPostDeleteReq) error {
	var err error
	var data models.SysPost

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