package handler

import (
	"github.com/gin-gonic/gin"
	jwt "go-admin/pkg/jwtauth"
	"log"
)

func NoFound(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Printf("NoRoute claims: %#v\n", claims)
	c.JSON(404, gin.H{
		"code":    "NOT_FOUND",
		"message": "not found",
	})
}
