package models

type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data interface{} `json:"data"`
	// 消息
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
}

type Page struct {
	List      interface{} `json:"list"`
	Count     int         `json:"count"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

// ReturnOK 正常返回
func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

// ReturnError 错误返回
func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	return res
}