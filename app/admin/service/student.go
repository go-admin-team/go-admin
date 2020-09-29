package service

import (
	"errors"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/common/service"
	"gorm.io/gorm"
)

type Student struct {
	service.Service
}

// GetStudentPage 获取student列表
func (e *Student) GetStudentPage(c cDto.Index, p *actions.DataPermission, list *[]models.Student, count *int64) error {
	var err error
	var data models.Student
	msgID := e.MsgID

	err = e.Orm.Table(data.TableName()).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// GetStudent 获取student对象
func (e *Student) GetStudent(c *dto.StudentSearch, p *actions.DataPermission, model *common.ActiveRecord) error {
	var err error
	var data models.Student
	msgID := e.MsgID

	db := e.Orm.Table(data.TableName()).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
		).
		First(model)
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

// InsertStudent 创建student对象
func (e *Student) InsertStudent(model common.ActiveRecord) error {
	var err error
	var data models.Student
	msgID := e.MsgID

	err = e.Orm.Table(data.TableName()).
		Create(model).Error
	if err != nil {
		log.Errorf("msgID[%s] db error:%s", msgID, err)
		return err
	}
	return nil
}

// UpdateStudent 修改student对象
func (e *Student) UpdateStudent(c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data models.Student
	msgID := e.MsgID

	db := e.Orm.Table(data.TableName()).
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

// RemoveStudent 删除student
func (e *Student) RemoveStudent(c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data models.Student
	msgID := e.MsgID

	db := e.Orm.Table(data.TableName()).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(c.GetId()).Delete(c)
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
