package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
)

// SysPostSearch 列表或者搜索使用结构体
type SysPostSearch struct {
	dto.Pagination `search:"-"`
	PostId         int    `form:"postId" search:"type:exact;column:post_id;table:sys_post" comment:"id"`     // id
	PostName       string `form:"postName" search:"type:exact;column:post_name;table:sys_post" comment:"名称"` // 名称
	PostCode       string `form:"postCode" search:"type:exact;column:post_code;table:sys_post" comment:"编码"` // 编码
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_post" comment:"排序"`          // 排序
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_post" comment:"状态"`      // 状态
	Remark         string `form:"remark" search:"type:exact;column:remark;table:sys_post" comment:"备注"`      // 备注
}

func (m *SysPostSearch) GetNeedSearch() interface{} {
	return *m
}

// Bind 映射上下文中的结构体数据
func (m *SysPostSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

// SysConfigControl 增、改使用的结构体
type SysPostControl struct {
	PostId   int    `uri:"id"  comment:"id"`        // id
	PostName string `form:"postName"  comment:"名称"` // 名称
	PostCode string `form:"postCode" comment:"编码"`  // 编码
	Sort     int    `form:"sort" comment:"排序"`      // 排序
	Status   int    `form:"status"   comment:"状态"`  // 状态
	Remark   string `form:"remark"   comment:"备注"`  // 备注
}

// Bind 映射上下文中的结构体数据
func (s *SysPostControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBindBodyWith(s, binding.JSON)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	var jsonStr []byte
	jsonStr, err = json.Marshal(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	ctx.Set("body", string(jsonStr))
	return err
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *SysPostControl) Generate() (*system.SysPost, error) {
	return &system.SysPost{
		PostId:   s.PostId,
		PostName: s.PostName,
		PostCode: s.PostCode,
		Sort:     s.Sort,
		Status:   s.Status,
		Remark:   s.Remark,
	}, nil
}

// GetId 获取数据对应的ID
func (s *SysPostControl) GetId() interface{} {
	return s.PostId
}

// SysConfigById 获取单个或者删除的结构体
type SysPostById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysPostById) Generate() *SysPostById {
	cp := *s
	return &cp
}

func (s *SysPostById) GetId() interface{} {
	return s.Id
}

func (s *SysPostById) Bind(ctx *gin.Context) error {
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

func (s *SysPostById) GenerateM() (*system.SysPost, error) {
	return &system.SysPost{}, nil
}
