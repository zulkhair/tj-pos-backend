package app

import (
	"dromatech/pos-backend/global"
	configdomain "dromatech/pos-backend/internal/domain/config"
	pinghandler "dromatech/pos-backend/internal/handler/ping"
	sessionhandler "dromatech/pos-backend/internal/handler/session"
	configrepo "dromatech/pos-backend/internal/repo/config"
	permissionrepo "dromatech/pos-backend/internal/repo/permission"
	rolerepo "dromatech/pos-backend/internal/repo/role"
	webuserrepo "dromatech/pos-backend/internal/repo/webuser"
	sessionusecase "dromatech/pos-backend/internal/usecase/session"
	"fmt"
	"strconv"
)

type AppHandler struct {
	pingHandler    *pinghandler.Handler
	sessionHandler *sessionhandler.Handler
}

func StartApp() error {
	// init repo
	configRepo, err := configrepo.New()
	if err != nil {
		return err
	}

	// set global config
	global.LOGIN_URL = configRepo.GetValue(configdomain.LOGIN_URL)
	global.UNAUTHORIZED_URL = configRepo.GetValue(configdomain.UNAUTHORIZED_URL)
	global.FORBIDDEN_URL = configRepo.GetValue(configdomain.FORBIDDEN_URL)

	sessionTimeout, _ := strconv.Atoi(configRepo.GetValue(configdomain.SESSION_TIMEOUT_MINUTE))
	global.SESSION_TIMEOUT_MINUTE = sessionTimeout

	webuserRepo, err := webuserrepo.New()
	if err != nil {
		return err
	}

	permissionRepo, err := permissionrepo.New()
	if err != nil {
		return err
	}

	roleRepo, err := rolerepo.New()
	if err != nil {
		return err
	}

	// init usecase
	sessionUsecase := sessionusecase.New(configRepo, webuserRepo, permissionRepo, roleRepo)

	// init Handler
	pingHandler := pinghandler.New()
	sessionHandler := sessionhandler.New(sessionUsecase)

	appHandler := AppHandler{
		pingHandler:    pingHandler,
		sessionHandler: sessionHandler,
	}

	router := newRoutes(appHandler)
	router.Run(fmt.Sprintf(global.CONFIG.Server.HTTP.Address))

	return nil
}
