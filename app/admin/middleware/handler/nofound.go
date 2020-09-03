package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	jwt "go-admin/pkg/jwtauth"
)

func NoFound(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Printf("NoRoute claims: %#v\n", claims)
	c.JSON(http.StatusOK, gin.H{
		"code":    "404",
		"message": "not found",
	})
}
