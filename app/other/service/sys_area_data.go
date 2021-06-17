package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/other/models"
	"go-admin/app/other/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysAreaData struct {
	service.Service
}

// GetPage 获取SysAreaData列表
func (e *SysAreaData) GetPage(c *dto.SysAreaDataSearch, list *[]models.SysAreaData) *SysAreaData {
	var data models.SysAreaData
	var areaDataList []models.SysAreaData
	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(&areaDataList).Error
	if err != nil {
		e.Log.Errorf("Service GetSysAreaDataPage error:%s", err)
		_ = e.AddError(err)
		return e
	}
	for i := 0; i < len(areaDataList); i++ {
		if areaDataList[i].PId != 86 {
			continue
		}
		menusInfo := areaDataCall(&areaDataList, areaDataList[i])
		*list = append(*list, menusInfo)
	}
	return e
}

// areaDataCall 构建树
func areaDataCall(areaDataList *[]models.SysAreaData, areaData models.SysAreaData) models.SysAreaData {
	list := *areaDataList

	min := make([]models.SysAreaData, 0)
	for j := 0; j < len(list); j++ {
		if areaData.Id != list[j].PId {
			continue
		}
		mi := models.SysAreaData{}
		mi = list[j]
		mi.Children = []models.SysAreaData{}
		ms := areaDataCall(areaDataList, mi)
		min = append(min, ms)
	}
	areaData.Children = min
	return areaData
}

// Get 获取SysAreaData对象
func (e *SysAreaData) Get(d *dto.SysAreaDataById, p *actions.DataPermission, model *models.SysAreaData) error {
	var data models.SysAreaData

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
		e.Log.Errorf("Service GetSysAreaData error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysAreaData对象
func (e *SysAreaData) Insert(model *models.SysAreaData) error {
	var data models.SysAreaData

	err := e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("Service InsertSysAreaData error:%s", err)
		return err
	}
	return nil
}

// Update 修改SysAreaData对象
func (e *SysAreaData) Update(c *models.SysAreaData, p *actions.DataPermission) error {
	db := e.Orm.Model(&models.SysAreaData{
		Id: c.GetId().(int),
	}).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Updates(c)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysAreaData error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SysAreaData
func (e *SysAreaData) Remove(d *dto.SysAreaDataById, p *actions.DataPermission) error {
	var data models.SysAreaData
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysAreaData error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}