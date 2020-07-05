package global

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var GinEngine *gin.Engine

var Eloquent *gorm.DB