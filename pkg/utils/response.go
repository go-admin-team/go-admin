package utils

import (
	"net/http"
	"time"
)

// api结构体
type APIException struct {
	Code      int         `json:"code"`
	Success   bool        `json:"success"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Result    interface{} `json:"result"`
}

// 实现接口
func (e *APIException) Error() string {
	return e.Msg
}

func newAPIException(code int, msg string, data interface{}, success bool) *APIException {
	return &APIException{
		Code:      code,
		Success:   success,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
		Result:    data,
	}
}

// 500 错误处理
func ServerError() *APIException {
	return newAPIException(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil, false)
}

// 404 错误
func NotFound() *APIException {
	return newAPIException(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil, false)
}

// 未知错误
func UnknownError(message string) *APIException {
	return newAPIException(http.StatusForbidden, message, nil, false)
}

// 参数错误
func ParameterError(message string) *APIException {
	return newAPIException(http.StatusBadRequest, message, nil, false)
}

// 授权错误
func AuthError(message string) *APIException {
	return newAPIException(http.StatusBadRequest, message, nil, false)
}

// 200
func ResponseJson(message string, data interface{}, success bool) *APIException {
	return newAPIException(http.StatusOK, message, data, success)
}
