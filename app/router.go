package app

import (
	"github.com/gin-gonic/gin"
)

func newRoutes(appHandler AppHandler) *gin.Engine {
	router := gin.Default()
	router.Use(appHandler.pingHandler.Ping)

	router.GET("/ping", appHandler.pingHandler.Ping)

	return router
}
