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

type SysFileDirSearch struct {
	dto.Pagination `search:"-"`

	ID    uint   `form:"ID" search:"type:exact;column:id;table:sys_file_dir" comment:"标识"`
	Label string `form:"label" search:"type:exact;column:label;table:sys_file_dir" comment:"目录名称"`
	PId   string `form:"pId" search:"type:exact;column:p_id;table:sys_file_dir" comment:"上级目录"`
	//Sort  string `form:"sort" search:"type:exact;column:sort;table:sys_file_dir" comment:"排序"`
	Path  string `form:"path" search:"type:exact;column:path;table:sys_file_dir" comment:"路径"`
}

func (m *SysFileDirSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysFileDirSearch) Bind(ctx *gin.Context) error {
	msgID := tools.GenerateMsgIDFromContext(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
	}
	return err
}

func (m *SysFileDirSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysFileDirControl struct {
	ID    uint   `uri:"ID" comment:"标识"` // 标识
	Label string `json:"label" comment:"目录名称"`
	PId   uint   `json:"pId" comment:"上级目录"`
	//Sort  string `json:"sort" comment:"排序"`
	Path  string `json:"path" comment:"路径"`
}

func (s *SysFileDirControl) Bind(ctx *gin.Context) error {
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

func (s *SysFileDirControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysFileDirControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysFileDir{
		Model: gorm.Model{ID: s.ID},
		Label: s.Label,
		PId:   s.PId,
		//Sort:  s.Sort,
		Path:  s.Path,
	}, nil
}

func (s *SysFileDirControl) GetId() interface{} {
	return s.ID
}

type SysFileDirById struct {
	dto.ObjectById
}

func (s *SysFileDirById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysFileDirById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysFileDir{}, nil
}
