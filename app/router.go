package app

import (
	"github.com/gin-gonic/gin"
)

func newRoutes(appHandler AppHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", appHandler.pingHandler.Ping)

	return router
}
