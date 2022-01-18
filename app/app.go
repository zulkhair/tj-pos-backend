package app

import (
	"dromatech/pos-backend/global"
	pinghandler "dromatech/pos-backend/internal/handler/ping"
	"fmt"
)

type AppHandler struct {
	pingHandler *pinghandler.Handler
}

func StartApp() {
	pingHandler := pinghandler.New()

	appHandler := AppHandler{
		pingHandler: pingHandler,
	}

	router := newRoutes(appHandler)
	router.Run(fmt.Sprintf(global.CONFIG.Server.HTTP.Address))
}
