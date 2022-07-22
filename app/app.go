package app

import (
	"dromatech/pos-backend/global"
	configdomain "dromatech/pos-backend/internal/domain/config"
	customerhandler "dromatech/pos-backend/internal/handler/customer"
	kontrabonhandler "dromatech/pos-backend/internal/handler/kontrabon"
	pinghandler "dromatech/pos-backend/internal/handler/ping"
	pricehandler "dromatech/pos-backend/internal/handler/price"
	producthandler "dromatech/pos-backend/internal/handler/product"
	rolehandler "dromatech/pos-backend/internal/handler/role"
	sessionhandler "dromatech/pos-backend/internal/handler/session"
	supplierhandler "dromatech/pos-backend/internal/handler/supplier"
	transactionhandler "dromatech/pos-backend/internal/handler/transaction"
	unithandler "dromatech/pos-backend/internal/handler/unit"
	webuserhandler "dromatech/pos-backend/internal/handler/webuser"
	configrepo "dromatech/pos-backend/internal/repo/config"
	customerrepo "dromatech/pos-backend/internal/repo/customer"
	kontrabonrepo "dromatech/pos-backend/internal/repo/kontrabon"
	pricerepo "dromatech/pos-backend/internal/repo/price"
	productrepo "dromatech/pos-backend/internal/repo/product"
	rolerepo "dromatech/pos-backend/internal/repo/role"
	sequencerepo "dromatech/pos-backend/internal/repo/sequence"
	supplierrepo "dromatech/pos-backend/internal/repo/supplier"
	transactionrepo "dromatech/pos-backend/internal/repo/transaction"
	unitrepo "dromatech/pos-backend/internal/repo/unit"
	webuserrepo "dromatech/pos-backend/internal/repo/webuser"
	customerusecase "dromatech/pos-backend/internal/usecase/customer"
	kontrabonusecase "dromatech/pos-backend/internal/usecase/kontrabon"
	priceusecase "dromatech/pos-backend/internal/usecase/price"
	productusecase "dromatech/pos-backend/internal/usecase/product"
	roleusecase "dromatech/pos-backend/internal/usecase/role"
	sessionusecase "dromatech/pos-backend/internal/usecase/session"
	supplierusecase "dromatech/pos-backend/internal/usecase/supplier"
	transactionusecase "dromatech/pos-backend/internal/usecase/transaction"
	unitusecase "dromatech/pos-backend/internal/usecase/unit"
	webuserusecase "dromatech/pos-backend/internal/usecase/webuser"
	"fmt"
	"strconv"
)

type AppHandler struct {
	pingHandler        *pinghandler.Handler
	sessionHandler     *sessionhandler.Handler
	webUserHander      *webuserhandler.Handler
	roleHandler        *rolehandler.Handler
	productHandler     *producthandler.Handler
	supplierHandler    *supplierhandler.Handler
	customerHandler    *customerhandler.Handler
	unitHandler        *unithandler.Handler
	transactionHandler *transactionhandler.Handler
	kontrabonHandler   *kontrabonhandler.Handler
	priceHandler       *pricehandler.Handler
}

func StartApp() error {
	// init config repo
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

	// init repo
	webuserRepo := webuserrepo.New()
	roleRepo := rolerepo.New()
	productRepo := productrepo.New()
	supplierRepo := supplierrepo.New()
	customerRepo := customerrepo.New()
	unitRepo := unitrepo.New()
	sequenceRepo := sequencerepo.New()
	transactionRepo := transactionrepo.New()
	kontrabonRepo := kontrabonrepo.New()
	priceRepo := pricerepo.New()

	// init usecase
	sessionUsecase := sessionusecase.New(configRepo, webuserRepo, roleRepo)
	webUserUsecase := webuserusecase.New(webuserRepo)
	roleUsecase := roleusecase.New(roleRepo)
	productUsecase := productusecase.New(productRepo)
	supplierUsecase := supplierusecase.New(supplierRepo)
	custmerUsecase := customerusecase.New(customerRepo)
	unitUsecase := unitusecase.New(unitRepo)
	transactionusecase := transactionusecase.New(transactionRepo, sequenceRepo, supplierRepo, customerRepo)
	kontrabonUseccase := kontrabonusecase.New(kontrabonRepo, sequenceRepo, customerRepo)
	priceUsecase := priceusecase.New(priceRepo, productRepo, customerRepo)

	// init Handler
	pingHandler := pinghandler.New()
	sessionHandler := sessionhandler.New(sessionUsecase)
	webUserHander := webuserhandler.New(webUserUsecase)
	rolehandler := rolehandler.New(roleUsecase)
	productHandler := producthandler.New(productUsecase)
	supplierHandler := supplierhandler.New(supplierUsecase)
	customerHandler := customerhandler.New(custmerUsecase)
	unitHandler := unithandler.New(unitUsecase)
	transactionHandler := transactionhandler.New(transactionusecase)
	kontrabonHandler := kontrabonhandler.New(kontrabonUseccase)
	priceHandler := pricehandler.New(priceUsecase)

	appHandler := AppHandler{
		pingHandler:        pingHandler,
		sessionHandler:     sessionHandler,
		webUserHander:      webUserHander,
		roleHandler:        rolehandler,
		productHandler:     productHandler,
		supplierHandler:    supplierHandler,
		customerHandler:    customerHandler,
		unitHandler:        unitHandler,
		transactionHandler: transactionHandler,
		kontrabonHandler:   kontrabonHandler,
		priceHandler:       priceHandler,
	}

	router := newRoutes(appHandler)
	router.Run(fmt.Sprintf(global.CONFIG.Server.HTTP.Address))

	return nil
}
