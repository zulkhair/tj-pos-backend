package app

import (
	"github.com/gin-gonic/gin"
)

func newRoutes(appHandler AppHandler) *gin.Engine {
	router := gin.Default()
	router.Use(appHandler.sessionHandler.AuthCheck)

	router.GET("/api/ping", appHandler.pingHandler.Ping)
	router.POST("/api/auth/login", appHandler.sessionHandler.Login)
	router.POST("/api/auth/logout", appHandler.sessionHandler.Logout)
	router.GET("/api/auth/getmenu", appHandler.sessionHandler.GetMenu)

	router.POST("/api/user/edit", appHandler.webUserHander.EditUser)
	router.POST("/api/user/change-password", appHandler.webUserHander.ChangePassword)

	return router
}
