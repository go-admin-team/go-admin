package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
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

				// 根据接口方法注释里的@Summary填充接口名称，适用于代码生成器
				// 可在此处增加配置路径前缀的if判断，只对代码生成的自建应用进行定向的接口名称填充
				jsonFile, _ := ioutil.ReadFile("docs/swagger.json")
				jsonData, _ := simplejson.NewFromReader(bytes.NewReader(jsonFile))
				urlPath := v.RelativePath
				idPatten := "(.*)/:(\\w+)" // 正则替换，把:id换成{id}
				reg, _ := regexp.Compile(idPatten)
				if reg.MatchString(urlPath) {
					urlPath = reg.ReplaceAllString(v.RelativePath, "${1}/{${2}}") // 把:id换成{id}
				}
				apiTitle, _ := jsonData.Get("paths").Get(urlPath).Get(strings.ToLower(v.HttpMethod)).Get("summary").String()

				err := d.Debug().Where(SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
					Attrs(SysApi{Handle: v.Handler, Title: apiTitle}).
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
