package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type SysDept struct {
	service.Service
}

// GetSysDeptPage 获取SysDept列表
func (e *SysDept) GetSysDeptPage(c *dto.SysDeptSearch, list *[]models.SysDept, count *int64) error {
	var err error
	var data models.SysDept
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetSysDept 获取SysDept对象
func (e *SysDept) GetSysDept(d *dto.SysDeptById, model *models.SysDept) error {
	var err error
	var data models.SysDept
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// InsertSysDept 创建SysDept对象
func (e *SysDept) InsertSysDept(model *models.SysDept) error {
	var err error
	var data models.SysDept
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateSysDept 修改SysDept对象
func (e *SysDept) UpdateSysDept(c *models.SysDept) error {
	var err error
	var data models.SysDept
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Where(c.GetId()).Updates(c)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysDept 删除SysDept
func (e *SysDept) RemoveSysDept(d *dto.SysDeptById) error {
	var err error
	var data models.SysDept
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Where(d.GetId()).Delete(&data)
	if db.Error != nil {
		err = db.Error
		log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("无权删除该数据")
		return err
	}
	return nil
}

// GetSysDeptList 获取组织数据
func (e *SysDept) GetSysDeptList(c *dto.SysDeptSearch, list *[]models.SysDept) error {
	var err error
	var data models.SysDept
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
		).
		Find(list).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// SetDeptTree 设置组织数据
func (e *SysDept) SetDeptTree() (c *dto.SysDeptSearch, m []dto.DeptLabel, err error) {
	var list []models.SysDept
	err = e.GetSysDeptList(c, &list)

	m = make([]dto.DeptLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.DeptLabel{}
		e.Id = list[i].DeptId
		e.Label = list[i].DeptName
		deptsInfo := DeptCall(&list, e)

		m = append(m, deptsInfo)
	}
	return
}

// Call 递归构造组织数据
func DeptCall(deptlist *[]models.SysDept, dept dto.DeptLabel) dto.DeptLabel {
	list := *deptlist
	min := make([]dto.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi := dto.DeptLabel{Id: list[j].DeptId, Label: list[j].DeptName, Children: []dto.DeptLabel{}}
		ms := DeptCall(deptlist, mi)
		min = append(min, ms)
	}
	dept.Children = min
	return dept
}
