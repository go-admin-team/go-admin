package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

/**
*
* @description :
*
* @author : codezhang
*
* @create : 2019-01-08 16:08
**/

func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}
