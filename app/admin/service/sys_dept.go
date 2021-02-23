package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/common/service"
	"go-admin/tools"
	"gorm.io/gorm"
)

type SysDept struct {
	service.Service
}

// GetSysDeptPage 获取SysDept列表
func (e *SysDept) GetSysDeptPage(c *dto.SysDeptSearch, list *[]models.SysDept) error {
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
	deptPath := "/" + tools.IntToString(model.DeptId)
	if model.ParentId != 0 {
		var deptP models.SysDepts
		e.Orm.Model(&data).First(&deptP, model.ParentId)
		deptPath = deptP.DeptPath + deptPath
	} else {
		deptPath = "/0" + deptPath
	}
	var mp = map[string]string{}
	mp["dept_path"] = deptPath
	if err := e.Orm.Model(&model).Update("dept_path",deptPath).Error; err != nil {
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

	//参数1:是要修改的数据
	//参数2:是修改的数据

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
func (e *SysDept) getSysDeptList(c *dto.SysDeptSearch, list *[]models.SysDept) error {
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
func (e *SysDept) SetDeptTree(c *dto.SysDeptSearch) (m []dto.DeptLabel, err error) {
	var list []models.SysDept
	err = e.getSysDeptList(c, &list)

	m = make([]dto.DeptLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.DeptLabel{}
		e.Id = list[i].DeptId
		e.Label = list[i].DeptName
		deptsInfo := deptTreeCall(&list, e)

		m = append(m, deptsInfo)
	}
	return
}

// Call 递归构造组织数据
func deptTreeCall(deptlist *[]models.SysDept, dept dto.DeptLabel) dto.DeptLabel {
	list := *deptlist
	min := make([]dto.DeptLabel, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi := dto.DeptLabel{Id: list[j].DeptId, Label: list[j].DeptName, Children: []dto.DeptLabel{}}
		ms := deptTreeCall(deptlist, mi)
		min = append(min, ms)
	}
	dept.Children = min
	return dept
}

// SetDeptPage 设置dept页面数据
func (e *SysDept) SetDeptPage(c *dto.SysDeptSearch) (m []models.SysDept, err error) {
	var list []models.SysDept
	err = e.getSysDeptList(c, &list)

	//m := make([]models.SysDept, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := e.deptPageCall(&list, list[i])

		m = append(m, info)
	}
	return
}

func  (e *SysDept) deptPageCall(deptlist *[]models.SysDept, menu models.SysDept) models.SysDept {
	list := *deptlist

	min := make([]models.SysDept, 0)
	for j := 0; j < len(list); j++ {

		if menu.DeptId != list[j].ParentId {
			continue
		}
		mi := models.SysDept{}
		mi.DeptId = list[j].DeptId
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.DeptName = list[j].DeptName
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.CreatedAt = list[j].CreatedAt
		mi.Children = []models.SysDept{}
		ms := e.deptPageCall(deptlist, mi)
		min = append(min, ms)

	}
	menu.Children = min
	return menu
}
