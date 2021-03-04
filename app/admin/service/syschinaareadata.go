package service

import (
	"errors"

	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/service"
)

type SysChinaAreaData struct {
	service.Service
}

// GetSysChinaAreaDataPage 获取SysChinaAreaData列表
func (e *SysChinaAreaData) GetSysChinaAreaDataPage(c *dto.SysChinaAreaDataSearch, p *actions.DataPermission, list *[]models.SysChinaAreaData, count *int64) error {
	var err error
	var data models.SysChinaAreaData

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetSysChinaAreaData 获取SysChinaAreaData对象
func (e *SysChinaAreaData) GetSysChinaAreaData(d *dto.SysChinaAreaDataById, p *actions.DataPermission, model *models.SysChinaAreaData) error {
	var err error
	var data models.SysChinaAreaData

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
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

// InsertSysChinaAreaData 创建SysChinaAreaData对象
func (e *SysChinaAreaData) InsertSysChinaAreaData(model *models.SysChinaAreaData) error {
	var err error
	var data models.SysChinaAreaData

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// UpdateSysChinaAreaData 修改SysChinaAreaData对象
func (e *SysChinaAreaData) UpdateSysChinaAreaData(c *models.SysChinaAreaData, p *actions.DataPermission) error {
	var err error
	var data models.SysChinaAreaData

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(c.GetId()).Updates(c)
	if db.Error != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysChinaAreaData 删除SysChinaAreaData
func (e *SysChinaAreaData) RemoveSysChinaAreaData(d *dto.SysChinaAreaDataById, p *actions.DataPermission) error {
	var err error
	var data models.SysChinaAreaData

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
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
