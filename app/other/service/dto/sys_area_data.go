package dto

import (
	"go-admin/app/other/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/common/dto"
)

type SysAreaDataSearch struct {
	dto.Pagination `search:"-"`
	PId            string `form:"pId" search:"type:exact;column:p_id;table:sys_china_area_data" comment:"上级编码"`
	Name           string `form:"name" search:"type:exact;column:name;table:sys_china_area_data" comment:"名称"`
}

func (m *SysAreaDataSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysAreaDataSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

type SysAreaDataControl struct {
	Id         int       `uri:"id" comment:"编码"` // 编码
	PId        int       `json:"pId" comment:"上级编码"`
	Name       string    `json:"name" comment:"名称"`
	CreateTime time.Time `json:"createTime" comment:""`
	UpdateTime time.Time `json:"updateTime" comment:""`
	DeleteTime time.Time `json:"deleteTime" comment:""`
}

func (s *SysAreaDataControl) Bind(ctx *gin.Context) error {
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

func (s *SysAreaDataControl) Generate() (*models.SysAreaData, error) {
	return &models.SysAreaData{
		Id:   s.Id,
		PId:  s.PId,
		Name: s.Name,
	}, nil
}

func (s *SysAreaDataControl) GetId() interface{} {
	return s.Id
}

type SysAreaDataById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysAreaDataById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysAreaDataById) Bind(ctx *gin.Context) error {
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

func (s *SysAreaDataById) SetUpdateBy(id int) {

}
