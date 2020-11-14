package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/common/service"
	"go-admin/tools"
	"gorm.io/gorm"
)

type SysFileDir struct {
	service.Service
}

// GetSysFileDirPage 获取SysFileDir列表
func (e *SysFileDir) GetSysFileDirPage(c cDto.Index, list *[]models.SysFileDirL) error {
	var err error
	var data models.SysFileDir
	msgID := e.MsgID

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			//cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list). //Limit(-1).Offset(-1).
		Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetSysFileDir 获取SysFileDir对象
func (e *SysFileDir) GetSysFileDir(d cDto.Control, model *models.SysFileDir) error {
	var err error
	var data models.SysFileDir
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

// InsertSysFileDir 创建SysFileDir对象
func (e *SysFileDir) InsertSysFileDir(model common.ActiveRecord) error {
	var err error
	var data models.SysFileDir
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Create(model)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	id := model.GetId()
	path := "/" + tools.UIntToString(id.(uint))
	db = e.Orm.Model(&data).
		First(&data, model.GetId())
	err = db.Error

	if data.PId != 0 {
		var dept models.SysFileDir
		e.Orm.Model(&data).Where("id = ?", data.PId).First(&dept)
		path = dept.Path + path
	} else {
		path = "/0" + path
	}
	//var mp = map[string]string{}
	//mp["path"] = path
	if err := e.Orm.Model(&data).Where("id = ?", data.ID).Update("path", path).Error; err != nil {
		return err
	}

	return nil
}

// UpdateSysFileDir 修改SysFileDir对象
func (e *SysFileDir) UpdateSysFileDir(c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data models.SysFileDir
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(c.GetId()).Updates(c)
	if db.Error != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysFileDir 删除SysFileDir
func (e *SysFileDir) RemoveSysFileDir(d cDto.Control, c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data models.SysFileDir
	msgID := e.MsgID

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(d.GetId()).Delete(c)
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

func (e *SysFileDir) SetSysFileDir(c cDto.Index) (*[]models.SysFileDirL, error) {
	var list []models.SysFileDirL
	err := e.GetSysFileDirPage(c, &list)
	m := make([]models.SysFileDirL, 0)
	for i := 0; i < len(list); i++ {
		if list[i].PId != 0 {
			continue
		}
		info := SysFileDirCall(&list, list[i])
		m = append(m, info)
	}
	return &m, err
}

func SysFileDirCall(list *[]models.SysFileDirL, m models.SysFileDirL) models.SysFileDirL {
	listGroup := *list
	min := make([]models.SysFileDirL, 0)
	for j := 0; j < len(listGroup); j++ {
		if m.ID != listGroup[j].PId {
			continue
		}
		mi := models.SysFileDirL{}
		mi.ID = listGroup[j].ID
		mi.PId = listGroup[j].PId
		mi.Label = listGroup[j].Label
		//mi.Sort = listGroup[j].Sort
		mi.CreatedAt = listGroup[j].CreatedAt
		mi.UpdatedAt = listGroup[j].UpdatedAt
		mi.Children = []models.SysFileDirL{}
		ms := SysFileDirCall(list, mi)
		min = append(min, ms)
	}
	if len(min) > 0 {
		m.Children = min
	} else {
		m.Children = nil
	}

	return m
}
