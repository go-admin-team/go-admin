package app

import (
	pbErrors "github.com/go-admin-team/go-admin-core/errors"
)

type Response struct {
	pbErrors.Error
	// 数据集
	Data interface{} `json:"data"`
}

type Page struct {
	List      interface{} `json:"list"`
	Count     int         `json:"count"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

func (res *Response) ReturnError(code int32) *Response {
	res.Code = code
	return res
}
