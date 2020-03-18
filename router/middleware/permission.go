package middleware

import (
	"github.com/gin-gonic/gin"
	mycasbin "go-admin/pkg/casbin"
	"go-admin/pkg/jwtauth"
	_ "go-admin/pkg/jwtauth"
	"log"
	"net/http"
)

//权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		//根据上下文获取载荷claims 从claims获得role
		//claims := c.MustGet("claims").(*jwt.MapClaims)
		data, _ := c.Get("JWT_PAYLOAD")
		v := data.(jwtauth.MapClaims)
		e, _ := mycasbin.Casbin()
		//检查权限
		res, err := e.Enforce(v["rolekey"], c.Request.URL.Path, c.Request.Method)
		log.Println("----------------", v["rolekey"], c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Println(err.Error())
		}
		if res {
			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": 401,
				"msg":    "对不起，您没有该接口访问权限，请联系管理员",
			})
			c.Abort()
			return
		}
	}
}
