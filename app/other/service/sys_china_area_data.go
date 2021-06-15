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

type SysChinaAreaData struct {
	service.Service
}

// GetPage 获取SysChinaAreaData列表
func (e *SysChinaAreaData) GetPage(c *dto.SysChinaAreaDataSearch, list *[]models.SysChinaAreaData) *SysChinaAreaData {
	var data models.SysChinaAreaData
	var areaDataList []models.SysChinaAreaData
	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(&areaDataList).Error
	if err != nil {
		e.Log.Errorf("Service GetSysChinaAreaDataPage error:%s", err)
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
func areaDataCall(areaDataList *[]models.SysChinaAreaData, areaData models.SysChinaAreaData) models.SysChinaAreaData {
	list := *areaDataList

	min := make([]models.SysChinaAreaData, 0)
	for j := 0; j < len(list); j++ {
		if areaData.Id != list[j].PId {
			continue
		}
		mi := models.SysChinaAreaData{}
		mi = list[j]
		mi.Children = []models.SysChinaAreaData{}
		ms := areaDataCall(areaDataList, mi)
		min = append(min, ms)
	}
	areaData.Children = min
	return areaData
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
		Id: c.GetId().(int),
	}).
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