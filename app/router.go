package app

import (
	"github.com/gin-gonic/gin"
)

func newRoutes(appHandler AppHandler) *gin.Engine {
	router := gin.Default()
	router.Use(appHandler.sessionHandler.AuthCheck)

	router.GET("/ping", appHandler.pingHandler.Ping)
	router.POST("/auth/login", appHandler.sessionHandler.Login)
	router.POST("/auth/logout", appHandler.sessionHandler.Logout)
	router.GET("/auth/getmenu", appHandler.sessionHandler.GetMenu)

	return router
}
