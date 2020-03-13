package nacos_error

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-14 11:22
**/

type NacosError struct {
	errorCode   string
	errMsg      string
	originError error
}

func NewNacosError(errorCode string, errMsg string, originError error) *NacosError {
	return &NacosError{
		errorCode:   errorCode,
		errMsg:      errMsg,
		originError: originError,
	}

}

func (err *NacosError) Error() (str string) {
	nacosErrMsg := fmt.Sprintf("[%s] %s", err.ErrorCode(), err.errMsg)
	if err.originError != nil {
		return nacosErrMsg + "\ncaused by:\n" + err.originError.Error()
	}
	return nacosErrMsg
}

func (err *NacosError) ErrorCode() string {
	if err.errorCode == "" {
		return constant.DefaultClientErrorCode
	} else {
		return err.errorCode
	}
}
