package file_store

import (
	"testing"
)

func TestOSSUpload(t *testing.T) {
	// 打括号内填写自己的测试信息即可
	e := OXS{}
	var oxs = e.Setup(AliYunOSS)
	err := oxs.UpLoad("test.png", "./test.png")
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}
