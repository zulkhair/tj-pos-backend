package global

import (
	"dromatech/pos-backend/config"

	"gorm.io/gorm"
)

var CONFIG *config.Config
var DBCON *gorm.DB
var LOGIN_URL string
var FORBIDDEN_URL string
var UNAUTHORIZED_URL string
var SESSION_TIMEOUT_MINUTE int
