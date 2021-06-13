package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/storage"

	"go-admin/common/models"
)

type SysApi struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Handle string `json:"handle" gorm:"size:128;comment:handle"`
	Title  string `json:"title" gorm:"size:128;comment:标题"`
	Path   string `json:"path" gorm:"size:128;comment:地址"`
	Action string `json:"action" gorm:"size:16;comment:请求类型"`
	Type   string `json:"type" gorm:"size:16;comment:接口类型"`
	models.ModelTime
	models.ControlBy
}

func (SysApi) TableName() string {
	return "sys_api"
}

func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysApi) GetId() interface{} {
	return e.Id
}

func SaveSysApi(message storage.Messager) (err error) {
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		fmt.Errorf("json Marshal error, %s", err.Error())
		return err
	}

	var l runtime.Routers
	err = json.Unmarshal(rb, &l)
	if err != nil {
		fmt.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	dbList := sdk.Runtime.GetDb()
	for _, d := range dbList {
		for _, v := range l.List {
			if v.HttpMethod != "HEAD" ||
				strings.Contains(v.RelativePath, "/swagger/") ||
				strings.Contains(v.RelativePath, "/static/") ||
				strings.Contains(v.RelativePath, "/form-generator/") ||
				strings.Contains(v.RelativePath, "/sys/tables") {
				err := d.Debug().Where(SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
					Attrs(SysApi{Handle: v.Handler}).
					FirstOrCreate(&SysApi{}).
					//Update("handle", v.Handler).
					Error
				if err != nil {
					err := fmt.Errorf("Models SaveSysApi error: %s \r\n ", err.Error())
					return err
				}
			}
		}
	}
	return nil
}
