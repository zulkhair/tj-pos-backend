package global

import (
	"dromatech/pos-backend/config"

	"gorm.io/gorm"
)

var CONFIG *config.Config
var DBCON *gorm.DB
