package test

import (
	"fmt"
	"go-admin/models"
	"testing"
)

func TestBackendUserEncrypt(t *testing.T) {
	u := models.BackendUser{"chengxiao", "Cheng",
		"Xiao", "admin", "master", "weiyi1314", "cheng8984@gmail.com"}
	u.Encrypt()
	//ValidePassword("weiyi")
	fmt.Println(u)

}

//func ValidePassword(password string) bool {
//	rightpass := "$2a$10$J6O9GuMikIuFuo4t4X28X.hjwVevCkotvW4.4z3BH4WAkG49dqhyi"
//	err := bcrypt.CompareHashAndPassword([]byte(rightpass), []byte(password))
//	if err != nil {
//		fmt.Println("pw wrong")
//		return false
//	} else {
//		fmt.Println("pw ok")
//		return true
//	}
//}
