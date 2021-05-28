package service

import (
	"errors"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"

	"go-admin/app/admin/models/system"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	common "go-admin/common/models"
	"go-admin/common/service"
)

type SysUser struct {
	service.Service
}

// GetSysUserPage 获取SysUser列表
func (e *SysUser) GetSysUserPage(c cDto.Index, p *actions.DataPermission, list *[]system.SysUser, count *int64) error {
	var err error
	var data system.SysUser

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Preload("Dept").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// GetSysUser 获取SysUser对象
func (e *SysUser) GetSysUser(d cDto.Control, p *actions.DataPermission, model *system.SysUser) error {
	var err error
	var data system.SysUser

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

// InsertSysUser 创建SysUser对象
func (e *SysUser) InsertSysUser(model common.ActiveRecord) error {
	var err error
	var data system.SysUser

	err = e.Orm.Model(&data).
		Create(model).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}

// UpdateSysUser 修改SysUser对象
func (e *SysUser) UpdateSysUser(c common.ActiveRecord, p *actions.DataPermission) error {
	var err error

	db := e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Where(c.GetId()).Updates(c)
	if db.Error != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	return nil
}

// RemoveSysUser 删除SysUser
func (e *SysUser) RemoveSysUser(d cDto.Control, c common.ActiveRecord, p *actions.DataPermission) error {
	var err error
	var data system.SysUser

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Where(d.GetId()).Delete(c)
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

// UpdateSysUserPwd 修改SysUser对象密码
func (e *SysUser) UpdateSysUserPwd(id int, oldPassword, newPassword string, p *actions.DataPermission) error {
	var err error

	if newPassword == "" {
		return nil
	}
	c := &system.SysUser{}

	err = e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Where(id).Select("UserId", "Password", "Salt").First(c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权更新该数据")
		}
		e.Log.Errorf("db error: %s", err)
		return err
	}
	var ok bool
	ok, err = pkg.CompareHashAndPassword(c.Password, oldPassword)
	if err != nil {
		e.Log.Errorf("CompareHashAndPassword error, %s", err.Error())
		return err
	}
	if !ok {
		err = errors.New("incorrect Password")
		e.Log.Warnf("user[%d] %s", id, err.Error())
		return err
	}
	c.Password = newPassword
	db := e.Orm.Model(c).Where(id).Select("Password", "Salt").Updates(c)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("set password error")
		log.Warnf("db update error")
		return err
	}
	return nil
}

func (e *SysUser) GetSysUserProfile(id int, user *system.SysUser, roles *[]system.SysRole, posts *[]system.SysPost) error {
	err := e.Orm.Preload("Dept").First(user, "user_id = ?", id).Error
	if err != nil {
		return err
	}
	err = e.Orm.Find(roles, user.RoleId).Error
	if err != nil {
		return err
	}
	err = e.Orm.Find(posts, user.PostIds).Error
	if err != nil {
		return err
	}

	return nil
}
