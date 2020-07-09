package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	jwt "go-admin/pkg/jwtauth"
	"net/http"
)

func NoFound(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Printf("NoRoute claims: %#v\n", claims)
	c.JSON(http.StatusOK, gin.H{
		"code":    "404",
		"message": "not found",
	})
}
