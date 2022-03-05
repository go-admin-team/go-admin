package file_store

import (
	"testing"
)

func TestOBSUpload(t *testing.T) {
	e := OXS{"", "", "", ""}
	var oxs = e.Setup(HuaweiOBS)
	err := oxs.UpLoad("test.png", "./test.png")
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}
