package app

import (
	"dromatech/pos-backend/global"
	configdomain "dromatech/pos-backend/internal/domain/config"
	pinghandler "dromatech/pos-backend/internal/handler/ping"
	rolehandler "dromatech/pos-backend/internal/handler/role"
	sessionhandler "dromatech/pos-backend/internal/handler/session"
	webuserhandler "dromatech/pos-backend/internal/handler/webuser"
	configrepo "dromatech/pos-backend/internal/repo/config"
	rolerepo "dromatech/pos-backend/internal/repo/role"
	webuserrepo "dromatech/pos-backend/internal/repo/webuser"
	roleusecase "dromatech/pos-backend/internal/usecase/role"
	sessionusecase "dromatech/pos-backend/internal/usecase/session"
	webuserusecase "dromatech/pos-backend/internal/usecase/webuser"
	"fmt"
	"strconv"
)

type AppHandler struct {
	pingHandler    *pinghandler.Handler
	sessionHandler *sessionhandler.Handler
	webUserHander  *webuserhandler.Handler
	roleHandler    *rolehandler.Handler
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

	webuserRepo := webuserrepo.New()
	if err != nil {
		return err
	}

	roleRepo := rolerepo.New()
	if err != nil {
		return err
	}

	// init usecase
	sessionUsecase := sessionusecase.New(configRepo, webuserRepo, roleRepo)
	webUserUsecase := webuserusecase.New(webuserRepo)
	roleUsecase := roleusecase.New(roleRepo)

	// init Handler
	pingHandler := pinghandler.New()
	sessionHandler := sessionhandler.New(sessionUsecase)
	webUserHander := webuserhandler.New(webUserUsecase)
	rolehandler := rolehandler.New(roleUsecase)

	appHandler := AppHandler{
		pingHandler:    pingHandler,
		sessionHandler: sessionHandler,
		webUserHander:  webUserHander,
		roleHandler:    rolehandler,
	}

	router := newRoutes(appHandler)
	router.Run(fmt.Sprintf(global.CONFIG.Server.HTTP.Address))

	return nil
}
