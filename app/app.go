package app

import (
	"dromatech/pos-backend/global"
	configdomain "dromatech/pos-backend/internal/domain/config"
	pinghandler "dromatech/pos-backend/internal/handler/ping"
	producthandler "dromatech/pos-backend/internal/handler/product"
	rolehandler "dromatech/pos-backend/internal/handler/role"
	sessionhandler "dromatech/pos-backend/internal/handler/session"
	webuserhandler "dromatech/pos-backend/internal/handler/webuser"
	configrepo "dromatech/pos-backend/internal/repo/config"
	productrepo "dromatech/pos-backend/internal/repo/product"
	rolerepo "dromatech/pos-backend/internal/repo/role"
	webuserrepo "dromatech/pos-backend/internal/repo/webuser"
	productusecase "dromatech/pos-backend/internal/usecase/product"
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
	productHandler *producthandler.Handler
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

	roleRepo := rolerepo.New()

	productRepo := productrepo.New()

	// init usecase
	sessionUsecase := sessionusecase.New(configRepo, webuserRepo, roleRepo)
	webUserUsecase := webuserusecase.New(webuserRepo)
	roleUsecase := roleusecase.New(roleRepo)
	productUsecase := productusecase.New(productRepo)

	// init Handler
	pingHandler := pinghandler.New()
	sessionHandler := sessionhandler.New(sessionUsecase)
	webUserHander := webuserhandler.New(webUserUsecase)
	rolehandler := rolehandler.New(roleUsecase)
	productHandler := producthandler.New(productUsecase)

	appHandler := AppHandler{
		pingHandler:    pingHandler,
		sessionHandler: sessionHandler,
		webUserHander:  webUserHander,
		roleHandler:    rolehandler,
		productHandler: productHandler,
	}

	router := newRoutes(appHandler)
	router.Run(fmt.Sprintf(global.CONFIG.Server.HTTP.Address))

	return nil
}
