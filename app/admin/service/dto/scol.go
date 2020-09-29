package dto

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/tools"
)

type ScolSearch struct {
	dto.Pagination     `search:"-"`
    Name string `form:"name" search:"type:contains;column:name;table:scol" comment:"名称"`

    
}

func (m *ScolSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *ScolSearch) Bind(ctx *gin.Context) error {
    msgID := tools.GenerateMsgIDFromContext(ctx)
    err := ctx.ShouldBind(m)
    if err != nil {
    	log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
    }
    return err
}

func (m *ScolSearch) Generate() dto.Index {
	o := *m
	return &o
}

type ScolControl struct {
    
    ID uint `uri:"ID" comment:""` // 

    Name string `json:"name" comment:"名称"`
    
}

func (s *ScolControl) Bind(ctx *gin.Context) error {
    msgID := tools.GenerateMsgIDFromContext(ctx)
    err := ctx.ShouldBindUri(s)
    if err != nil {
        log.Debugf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
        return err
    }
    err = ctx.ShouldBind(s)
    if err != nil {
        log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
    }
    return err
}

func (s *ScolControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *ScolControl) GenerateM() (common.ActiveRecord, error) {
	return &models.Scol{
	
        Model:     gorm.Model{ID: s.ID},
        Name:  s.Name,
	}, nil
}

func (s *ScolControl) GetId() interface{} {
	return s.ID
}

type ScolById struct {
	dto.ObjectById
}

func (s *ScolById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *ScolById) GenerateM() (common.ActiveRecord, error) {
	return &models.Scol{}, nil
}
