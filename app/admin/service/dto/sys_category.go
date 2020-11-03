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

type SysCategorySearch struct {
	dto.Pagination     `search:"-"`
    Name string `form:"name" search:"type:exact;column:name;table:sys_category" comment:"名称"`

    Status string `form:"status" search:"type:exact;column:status;table:sys_category" comment:"状态"`

    
}

func (m *SysCategorySearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysCategorySearch) Bind(ctx *gin.Context) error {
    msgID := tools.GenerateMsgIDFromContext(ctx)
    err := ctx.ShouldBind(m)
    if err != nil {
    	log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
    }
    return err
}

func (m *SysCategorySearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysCategoryControl struct {
    
    ID uint `uri:"ID" comment:"标识"` // 标识

    Name string `json:"name" comment:"名称"`
    

    Img string `json:"img" comment:"图标"`
    

    Sort string `json:"sort" comment:"排序"`
    

    Status string `json:"status" comment:"状态"`
    

    Remark string `json:"remark" comment:"备注"`
    
}

func (s *SysCategoryControl) Bind(ctx *gin.Context) error {
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

func (s *SysCategoryControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysCategoryControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysCategory{
	
        Model:     gorm.Model{ID: s.ID},
        Name:  s.Name,
        Img:  s.Img,
        Sort:  s.Sort,
        Status:  s.Status,
        Remark:  s.Remark,
	}, nil
}

func (s *SysCategoryControl) GetId() interface{} {
	return s.ID
}

type SysCategoryById struct {
	dto.ObjectById
}

func (s *SysCategoryById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysCategoryById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysCategory{}, nil
}
