package orm

import (
	"github.com/jinzhu/gorm"
)

var Eloquent *gorm.DB
var Source string
var Driver string
var DBName string
