package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysLoginLogSearch struct {
	dto.Pagination `search:"-"`
	Username       string `form:"username" search:"type:exact;column:username;table:sys_login_log" comment:"用户名"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_login_log" comment:"状态"`
	Ipaddr         string `form:"ipaddr" search:"type:exact;column:ipaddr;table:sys_login_log" comment:"ip地址"`
	LoginLocation  string `form:"loginLocation" search:"type:exact;column:login_location;table:sys_login_log" comment:"归属地"`
}

func (m *SysLoginLogSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysLoginLogSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

type SysLoginLogControl struct {
	ID            int       `uri:"Id" comment:"主键"` // 主键
	Username      string    `json:"username" comment:"用户名"`
	Status        string    `json:"status" comment:"状态"`
	Ipaddr        string    `json:"ipaddr" comment:"ip地址"`
	LoginLocation string    `json:"loginLocation" comment:"归属地"`
	Browser       string    `json:"browser" comment:"浏览器"`
	Os            string    `json:"os" comment:"系统"`
	Platform      string    `json:"platform" comment:"固件"`
	LoginTime     time.Time `json:"loginTime" comment:"登录时间"`
	Remark        string    `json:"remark" comment:"备注"`
	Msg           string    `json:"msg" comment:"信息"`
}

func (s *SysLoginLogControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysLoginLogControl) Generate() (*system.SysLoginLog, error) {
	return &system.SysLoginLog{
		Model:         common.Model{Id: s.ID},
		Username:      s.Username,
		Status:        s.Status,
		Ipaddr:        s.Ipaddr,
		LoginLocation: s.LoginLocation,
		Browser:       s.Browser,
		Os:            s.Os,
		Platform:      s.Platform,
		LoginTime:     s.LoginTime,
		Remark:        s.Remark,
		Msg:           s.Msg,
	}, nil
}

func (s *SysLoginLogControl) GetId() interface{} {
	return s.ID
}

type SysLoginLogById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysLoginLogById) GetId() interface{} {
	return s.Id
}

func (s *SysLoginLogById) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(&s.Ids)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysLoginLogById) Generate() *SysLoginLogById {
	cp := *s
	return &cp
}

func (s *SysLoginLogById) GenerateM() (*system.SysLoginLog, error) {
	return &system.SysLoginLog{}, nil
}
