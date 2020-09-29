package dto

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
)

type StudentSearch struct {
	dto.Pagination     `search:"-"`
    Name string `form:"name" search:"type:exact;column:name;table:student" comment:"姓名"`
}

func (m *StudentSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *StudentSearch) Bind(ctx *gin.Context) error {
    err := ctx.Bind(m)
    if err != nil {
    	log.Errorf("MsgID[%s] Bind error: %#v", err)
    }
    return err
}

func (m *StudentSearch) Generate() dto.Index {
	o := *m
	return &o
}

type StudentControl struct {
    ID uint `uri:"ID" comment:""` //
    Name string `json:"name" comment:"姓名"`
}

func (s *StudentControl) Bind(ctx *gin.Context) error {
    err := ctx.BindUri(s)
    if err != nil {
        log.Errorf("MsgID[%s] Bind error: %#v", err)
        return err
    }
    err = ctx.Bind(s)
    if err != nil {
        log.Errorf("MsgID[%s] Bind error: %#v", err)
    }
    return err
}

func (s *StudentControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *StudentControl) GenerateM() (common.ActiveRecord, error) {
	return &models.Student{
	
        Model:     gorm.Model{ID: s.ID},
        Name:  s.Name,
	}, nil
}

func (s *StudentControl) GetId() interface{} {
	return s.ID
}

type StudentById struct {
	dto.ObjectById
}

func (s *StudentById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *StudentById) GenerateM() (common.ActiveRecord, error) {
	return &models.Student{}, nil
}

type StudentItem struct {
	gorm.Model
	common.ControlBy
	Name string `json:"name" gorm:"type:varchar(255);comment:姓名"` //
}
