package pkg

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

func StrToInt(err error, index string) int {
	result, err := strconv.Atoi(index)
	if err != nil {
		AssertErr(err, "string to int error"+err.Error(), -1)
	}
	return result
}

//加密
func Encrypt(e string) ([]byte, error) {
	//var s []byte
	s, err := bcrypt.GenerateFromPassword([]byte(e), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s, err
}

func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	return true, nil
}

// Assert 条件断言
// 当断言条件为 假 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
func Assert(condition bool, msg string, code ...int) {
	if !condition {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		panic("CustomErroe#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

// AssertErr 错误断言
// 当 error 不为 nil 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
// 若 msg 为空，则默认为 error 中的内容
func AssertErr(err error, msg string, code ...int) {
	if err != nil {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		if msg == "" {
			msg = err.Error()
		}
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}
