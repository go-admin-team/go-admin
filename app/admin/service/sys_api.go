package service

import (
	"errors"
	"fmt"

	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysApi struct {
	service.Service
}

// GetSysApiPage 获取SysApi列表
func (e *SysApi) GetSysApiPage(c *dto.SysApiSearch, p *actions.DataPermission, list *[]models.SysApi, count *int64) error {
	var err error
	var data models.SysApi

	err = e.Orm.Debug().Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("Service GetSysApiPage error:%s", err)
		return err
	}
	return nil
}

// GetSysApi 获取SysApi对象
func (e *SysApi) GetSysApi(d *dto.SysApiById, p *actions.DataPermission, model *models.SysApi) *SysApi {
	var data models.SysApi

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysApi error:%s", err)
		e.AddError(err)
		return e
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		e.AddError(err)
		return e
	}
	return nil
}

// InsertSysApi 创建SysApi对象
func (e *SysApi) CheckStorageSysApi(c *[]runtime.Router) error {
	for _, v := range *c {
		err := e.Orm.Debug().Where(models.SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
			Attrs(models.SysApi{Handle: v.Handler}).
			FirstOrCreate(&models.SysApi{}).Error
		if err != nil {
			err := fmt.Errorf("Service CheckStorageSysApi error: %s \r\n ", err.Error())
			return err
		}
	}
	return nil
}

// UpdateSysApi 修改SysApi对象
func (e *SysApi) UpdateSysApi(c *models.SysApi, p *actions.DataPermission) error {
	db := e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Where(c.GetId()).Updates(c)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysApi error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// RemoveSysApi 删除SysApi
func (e *SysApi) RemoveSysApi(d *dto.SysApiById, p *actions.DataPermission) error {
	var data models.SysApi

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysApi error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}