package dto

import (

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type TApiZlGetPageReq struct {
	dto.Pagination     `search:"-"`
    TApiZlOrder
}

type TApiZlOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:t_api_zl"`
    Code string `form:"codeOrder"  search:"type:order;column:code;table:t_api_zl"`
    Handle string `form:"handleOrder"  search:"type:order;column:handle;table:t_api_zl"`
    Title string `form:"titleOrder"  search:"type:order;column:title;table:t_api_zl"`
    Path string `form:"pathOrder"  search:"type:order;column:path;table:t_api_zl"`
    Type string `form:"typeOrder"  search:"type:order;column:type;table:t_api_zl"`
    Action string `form:"actionOrder"  search:"type:order;column:action;table:t_api_zl"`
    Req string `form:"reqOrder"  search:"type:order;column:req;table:t_api_zl"`
    Res string `form:"resOrder"  search:"type:order;column:res;table:t_api_zl"`
    ResError string `form:"resErrorOrder"  search:"type:order;column:res_error;table:t_api_zl"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:t_api_zl"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:t_api_zl"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:t_api_zl"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:t_api_zl"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:t_api_zl"`
    
}

func (m *TApiZlGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type TApiZlInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Code string `json:"code" comment:"业务编码"`
    Handle string `json:"handle" comment:"handle"`
    Title string `json:"title" comment:"标题"`
    Path string `json:"path" comment:"地址"`
    Type string `json:"type" comment:"接口类型"`
    Action string `json:"action" comment:"请求类型"`
    Req string `json:"req" comment:"请求入参"`
    Res string `json:"res" comment:"响应参数"`
    ResError string `json:"resError" comment:"错误返回"`
    common.ControlBy
}

func (s *TApiZlInsertReq) Generate(model *models.TApiZl)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Code = s.Code
    model.Handle = s.Handle
    model.Title = s.Title
    model.Path = s.Path
    model.Type = s.Type
    model.Action = s.Action
    model.Req = s.Req
    model.Res = s.Res
    model.ResError = s.ResError
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *TApiZlInsertReq) GetId() interface{} {
	return s.Id
}

type TApiZlUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Code string `json:"code" comment:"业务编码"`
    Handle string `json:"handle" comment:"handle"`
    Title string `json:"title" comment:"标题"`
    Path string `json:"path" comment:"地址"`
    Type string `json:"type" comment:"接口类型"`
    Action string `json:"action" comment:"请求类型"`
    Req string `json:"req" comment:"请求入参"`
    Res string `json:"res" comment:"响应参数"`
    ResError string `json:"resError" comment:"错误返回"`
    common.ControlBy
}

func (s *TApiZlUpdateReq) Generate(model *models.TApiZl)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Code = s.Code
    model.Handle = s.Handle
    model.Title = s.Title
    model.Path = s.Path
    model.Type = s.Type
    model.Action = s.Action
    model.Req = s.Req
    model.Res = s.Res
    model.ResError = s.ResError
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *TApiZlUpdateReq) GetId() interface{} {
	return s.Id
}

// TApiZlGetReq 功能获取请求参数
type TApiZlGetReq struct {
     Id int `uri:"id"`
}
func (s *TApiZlGetReq) GetId() interface{} {
	return s.Id
}

// TApiZlDeleteReq 功能删除请求参数
type TApiZlDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *TApiZlDeleteReq) GetId() interface{} {
	return s.Ids
}
