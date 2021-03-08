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

type SysFileInfo struct {
	service.Service
}

// GetSysFileInfoPage 获取SysFileInfo列表
func (e *SysFileInfo) GetSysFileInfoPage(c *dto.SysFileInfoSearch, p *actions.DataPermission, list *[]models.SysFileInfo, count *int64) error {
	var err error
	var data models.SysFileInfo

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// GetSysFileInfo 获取SysFileInfo对象
func (e *SysFileInfo) GetSysFileInfo(d *dto.SysFileInfoById, p *actions.DataPermission, model *models.SysFileInfo) error {
	var err error
	var data models.SysFileInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId())
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

// InsertSysFileInfo 创建SysFileInfo对象
func (e *SysFileInfo) InsertSysFileInfo(model *dto.SysFileInfoControl) error {
	var err error
	var data *models.SysFileInfo

	data, err = model.Generate()
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	err = e.Orm.Model(&data).
		Create(data).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// UpdateSysFileInfo 修改SysFileInfo对象
func (e *SysFileInfo) UpdateSysFileInfo(c *dto.SysFileInfoControl, p *actions.DataPermission) error {
	var err error
	var data *models.SysFileInfo

	data, err = c.Generate()
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	err = e.Orm.Debug().Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where("id = ?", c.ID).Updates(&data).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysFileInfo 删除SysFileInfo
func (e *SysFileInfo) RemoveSysFileInfo(d *dto.SysFileInfoById, p *actions.DataPermission) error {
	var err error
	var data models.SysFileInfo

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(d.GetId()).Delete(&data)
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
