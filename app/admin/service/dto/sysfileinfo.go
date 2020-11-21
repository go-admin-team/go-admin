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

type SysFileInfoSearch struct {
	dto.Pagination `search:"-"`
	ID             uint `form:"ID" search:"type:exact;column:id;table:sys_file_info" comment:"标识"`

	Type string `form:"type" search:"type:exact;column:type;table:sys_file_info" comment:"类型"`

	Name string `form:"name" search:"type:exact;column:name;table:sys_file_info" comment:"名称"`

	PId string `form:"pId" search:"type:exact;column:p_id;table:sys_file_info" comment:"目录"`

	Source string `form:"source" search:"type:exact;column:source;table:sys_file_info" comment:"来源"`
}

func (m *SysFileInfoSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysFileInfoSearch) Bind(ctx *gin.Context) error {
	msgID := tools.GenerateMsgIDFromContext(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
	}
	return err
}

func (m *SysFileInfoSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysFileInfoControl struct {
	ID      uint   `uri:"ID" comment:"标识"` // 标识
	Type    string `json:"type" comment:"类型"`
	Name    string `json:"name" comment:"名称"`
	Size    string `json:"size" comment:"大小"`
	PId     uint   `json:"pId" comment:"目录"`
	Source  string `json:"source" comment:"来源"`
	Url     string `json:"url" comment:"地址"`
	FullUrl string `json:"fullUrl" comment:"全地址"`
}

func (s *SysFileInfoControl) Bind(ctx *gin.Context) error {
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

func (s *SysFileInfoControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysFileInfoControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysFileInfo{

		Model:   gorm.Model{ID: s.ID},
		Type:    s.Type,
		Name:    s.Name,
		Size:    s.Size,
		PId:     s.PId,
		Source:  s.Source,
		Url:     s.Url,
		FullUrl: s.FullUrl,
	}, nil
}

func (s *SysFileInfoControl) GetId() interface{} {
	return s.ID
}

type SysFileInfoById struct {
	dto.ObjectById
}

func (s *SysFileInfoById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysFileInfoById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysFileInfo{}, nil
}
