package dto

import (
	"go-admin/app/other/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/common/dto"
)

type SysChinaAreaDataSearch struct {
	dto.Pagination `search:"-"`
	PId            string `form:"pId" search:"type:exact;column:p_id;table:sys_china_area_data" comment:"上级编码"`
	Name           string `form:"name" search:"type:exact;column:name;table:sys_china_area_data" comment:"名称"`
}

func (m *SysChinaAreaDataSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysChinaAreaDataSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

type SysChinaAreaDataControl struct {
	Id         int       `uri:"id" comment:"编码"` // 编码
	PId        int       `json:"pId" comment:"上级编码"`
	Name       string    `json:"name" comment:"名称"`
	CreateTime time.Time `json:"createTime" comment:""`
	UpdateTime time.Time `json:"updateTime" comment:""`
	DeleteTime time.Time `json:"deleteTime" comment:""`
}

func (s *SysChinaAreaDataControl) Bind(ctx *gin.Context) error {
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

func (s *SysChinaAreaDataControl) Generate() (*models.SysChinaAreaData, error) {
	return &models.SysChinaAreaData{
		Id:   s.Id,
		PId:  s.PId,
		Name: s.Name,
	}, nil
}

func (s *SysChinaAreaDataControl) GetId() interface{} {
	return s.Id
}

type SysChinaAreaDataById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysChinaAreaDataById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysChinaAreaDataById) Bind(ctx *gin.Context) error {
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

func (s *SysChinaAreaDataById) SetUpdateBy(id int) {

}
