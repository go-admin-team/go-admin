package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
)

type SysDictType struct {
	service.Service
}

// GetPage 获取列表
func (e *SysDictType) GetPage(c *dto.SysDictTypeGetPageReq, list *[]models.SysDictType, count *int64) error {
	var err error
	var data models.SysDictType

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Get 获取对象
func (e *SysDictType) Get(d *dto.SysDictTypeGetReq, model *models.SysDictType) error {
	var err error

	db := e.Orm.First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if db.Error != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Insert 创建对象
func (e *SysDictType) Insert(c *dto.SysDictTypeInsertReq) error {
	var err error
	var data models.SysDictType
	c.Generate(&data)
	var count int64
	e.Orm.Model(&data).Where("dict_type = ?", data.DictType).Count(&count)
	if count > 0 {
		return errors.New(fmt.Sprintf("当前字典类型[%s]已经存在！", data.DictType))
	}
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// Update 修改对象
func (e *SysDictType) Update(c *dto.SysDictTypeUpdateReq) error {
	var err error
	var model = models.SysDictType{}
	e.Orm.First(&model, c.GetId())
	c.Generate(&model)
	db := e.Orm.Save(&model)
	if db.Error != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// Remove 删除
func (e *SysDictType) Remove(d *dto.SysDictTypeDeleteReq) error {
	var err error
	var data models.SysDictType

	db := e.Orm.Delete(&data, d.GetId())
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

// GetAll 获取所有
func (e *SysDictType) GetAll(c *dto.SysDictTypeGetPageReq, list *[]models.SysDictType) error {
	var err error
	var data models.SysDictType

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}
