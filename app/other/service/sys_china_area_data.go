package service

import (
	"errors"
	"go-admin/app/other/models"
	"go-admin/app/other/service/dto"
	common "go-admin/common/models"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysChinaAreaData struct {
	service.Service
}

// GetPage 获取SysChinaAreaData列表
func (e *SysChinaAreaData) GetPage(c *dto.SysChinaAreaDataSearch, p *actions.DataPermission, list *[]models.SysChinaAreaData, count *int64) error {
	var data models.SysChinaAreaData

	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("Service GetSysChinaAreaDataPage error:%s", err)
		return err
	}
	return nil
}

// Get 获取SysChinaAreaData对象
func (e *SysChinaAreaData) Get(d *dto.SysChinaAreaDataById, p *actions.DataPermission, model *models.SysChinaAreaData) error {
	var data models.SysChinaAreaData

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("Service GetSysChinaAreaData error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysChinaAreaData对象
func (e *SysChinaAreaData) Insert(model *models.SysChinaAreaData) error {
	var data models.SysChinaAreaData

	err := e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("Service InsertSysChinaAreaData error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysChinaAreaData对象
func (e *SysChinaAreaData) Update(c *models.SysChinaAreaData, p *actions.DataPermission) error {
	db := e.Orm.Model(&models.SysChinaAreaData{
		Model: common.Model{
			Id: c.GetId().(int),
		}}).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Updates(c)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysChinaAreaData error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SysChinaAreaData
func (e *SysChinaAreaData) Remove(d *dto.SysChinaAreaDataById, p *actions.DataPermission) error {

	var data models.SysChinaAreaData

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())

	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysChinaAreaData error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
