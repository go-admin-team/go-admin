package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jwt "go-admin/pkg/jwtauth"
	"log"
)

func ExtractClaims(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}

func GetUserId(c *gin.Context) int {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return int((data["identity"]).(float64))
	}
	log.Println("****************************** 路径：" + c.Request.URL.Path + "  请求方法：" + c.Request.Method + "  说明：缺少identity")
	return 0
}

func GetUserIdStr(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return Int64ToString(int64((data["identity"]).(float64)))
	}
	log.Println("****************************** 路径：" + c.Request.URL.Path + "  请求方法：" + c.Request.Method + "  缺少identity")
	return ""
}

func GetUserName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["nice"] != nil {
		return (data["nice"]).(string)
	}
	fmt.Println("****************************** 路径：" + c.Request.URL.Path + "  请求方法：" + c.Request.Method + "  缺少nice")
	return ""
}

func GetRoleName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["rolekey"] != nil {
		return (data["rolekey"]).(string)
	}
	fmt.Println("****************************** 路径：" + c.Request.URL.Path + "  请求方法：" + c.Request.Method + "  缺少rolekey")
	return ""
}

func GetRoleId(c *gin.Context) int {
	data := ExtractClaims(c)
	if data["roleid"] != nil {
		i := int((data["roleid"]).(float64))
		return i
	}
	fmt.Println("****************************** 路径：" + c.Request.URL.Path + "  请求方法：" + c.Request.Method + "  缺少roleid")
	return 0
}
