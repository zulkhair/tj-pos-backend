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
	router.GET("/api/auth/check", appHandler.sessionHandler.CheckPermission)

	router.POST("/api/user/edit", appHandler.webUserHander.EditUser)
	router.POST("/api/user/change-password", appHandler.webUserHander.ChangePassword)
	router.POST("/api/user/force-change-password", appHandler.webUserHander.ForceChangePassword)
	router.POST("/api/user/register-user", appHandler.webUserHander.RegisterUser)
	router.GET("/api/user/find-all", appHandler.webUserHander.FindAllUser)
	router.POST("/api/user/change-status", appHandler.webUserHander.ChangeStatus)

	router.GET("/api/role/active-list", appHandler.roleHandler.GetActive)
	router.GET("/api/role/find-all", appHandler.roleHandler.GetAll)
	router.GET("/api/role/permissions", appHandler.roleHandler.FindPermissions)
	router.POST("/api/role/create", appHandler.roleHandler.RegisterRole)
	router.POST("/api/role/edit", appHandler.roleHandler.EditRole)

	router.GET("/api/product/find", appHandler.productHandler.Find)
	router.POST("/api/product/edit", appHandler.productHandler.Edit)
	router.POST("/api/product/create", appHandler.productHandler.Create)

	router.GET("/api/supplier/find", appHandler.supplierHandler.Find)
	router.POST("/api/supplier/edit", appHandler.supplierHandler.Edit)
	router.POST("/api/supplier/create", appHandler.supplierHandler.Create)
	router.GET("/api/supplier/buy-price", appHandler.supplierHandler.GetBuyPrice)
	router.POST("/api/supplier/update-buy-price", appHandler.supplierHandler.UpdateBuyPrice)
	router.POST("/api/supplier/add-price", appHandler.supplierHandler.AddPrice)
	router.GET("/api/supplier/find-latest-price", appHandler.supplierHandler.FindLatestPrice)
	router.GET("/api/supplier/find-price", appHandler.supplierHandler.FindPrice)

	router.GET("/api/customer/find", appHandler.customerHandler.Find)
	router.POST("/api/customer/edit", appHandler.customerHandler.Edit)
	router.POST("/api/customer/create", appHandler.customerHandler.Create)
	router.GET("/api/customer/sell-price", appHandler.customerHandler.GetSellPrice)
	router.POST("/api/customer/update-sell-price", appHandler.customerHandler.UpdateSellPrice)
	router.POST("/api/customer/add-price", appHandler.customerHandler.AddPrice)
	router.GET("/api/customer/find-latest-price", appHandler.customerHandler.FindLatestPrice)
	router.GET("/api/customer/find-price", appHandler.customerHandler.FindPrice)

	router.GET("/api/unit/find", appHandler.unitHandler.Find)
	router.POST("/api/unit/edit", appHandler.unitHandler.Edit)
	router.POST("/api/unit/create", appHandler.unitHandler.Create)

	router.GET("/api/transaction/find", appHandler.transactionHandler.Find)
	router.POST("/api/transaction/create", appHandler.transactionHandler.Create)
	router.POST("/api/transaction/updateStatus", appHandler.transactionHandler.UpdateStatus)
	router.POST("/api/transaction/updateBuyPrice", appHandler.transactionHandler.UpdateBuyPrice)

	return router
}
