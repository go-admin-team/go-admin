package service

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type TApiZl struct {
	service.Service
}

// GetPage 获取TApiZl列表
func (e *TApiZl) GetPage(c *dto.TApiZlGetPageReq, p *actions.DataPermission, list *[]models.TApiZl, count *int64) error {
	var err error
	var data models.TApiZl

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("TApiZlService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取TApiZl对象
func (e *TApiZl) Get(d *dto.TApiZlGetReq, p *actions.DataPermission, model *models.TApiZl) error {
	var data models.TApiZl

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetTApiZl error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建TApiZl对象
func (e *TApiZl) Insert(c *dto.TApiZlInsertReq) error {
    var err error
    var data models.TApiZl
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("TApiZlService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改TApiZl对象
func (e *TApiZl) Update(c *dto.TApiZlUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.TApiZl{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("TApiZlService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除TApiZl
func (e *TApiZl) Remove(d *dto.TApiZlDeleteReq, p *actions.DataPermission) error {
	var data models.TApiZl

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveTApiZl error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
