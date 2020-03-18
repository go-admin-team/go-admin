package models

type Response struct {
	// 代码
	Code int `json:"code" example:"200"`
	// 数据集
	Data interface{} `json:"data"`
	// 错误消息
	Msg string `json:"msg"`
}

func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	return res
}
