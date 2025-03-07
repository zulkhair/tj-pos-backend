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

	router.POST("/api/user/edit", appHandler.webUserHander.EditName)
	router.POST("/api/user/edit-user", appHandler.webUserHander.EditUser)
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
	router.GET("/api/product/findActive", appHandler.productHandler.FindActive)
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
	router.GET("/api/customer/findActive", appHandler.customerHandler.Find)
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
	router.GET("/api/unit/findActive", appHandler.unitHandler.FindActive)

	router.GET("/api/transaction/find", appHandler.transactionHandler.Find)
	router.POST("/api/transaction/create", appHandler.transactionHandler.Create)
	router.POST("/api/transaction/updateStatus", appHandler.transactionHandler.UpdateStatus)
	router.POST("/api/transaction/updateBuyPrice", appHandler.transactionHandler.UpdateBuyPrice)
	router.POST("/api/transaction/cancelTrx", appHandler.transactionHandler.CancelTrx)
	router.POST("/api/transaction/update", appHandler.transactionHandler.Update)
	router.GET("/api/transaction/report", appHandler.transactionHandler.FindReport)
	router.POST("/api/transaction/updateHargaBeli", appHandler.transactionHandler.UpdateHargaBeli)
	router.POST("/api/transaction/insertTransactionBuy", appHandler.transactionHandler.InsertTransactionBuy)
	router.GET("/api/transaction/findCustomerCredit", appHandler.transactionHandler.FindCustomerCredit)
	router.GET("/api/transaction/findCustomerReport", appHandler.transactionHandler.FindCustomerReport)

	router.GET("/api/kontrabon/find", appHandler.kontrabonHandler.Find)
	router.GET("/api/kontrabon/findTransaction", appHandler.kontrabonHandler.FindTransaction)
	router.POST("/api/kontrabon/create", appHandler.kontrabonHandler.Create)
	router.POST("/api/kontrabon/add", appHandler.kontrabonHandler.Add)
	router.POST("/api/kontrabon/remove", appHandler.kontrabonHandler.Remove)
	router.POST("/api/kontrabon/update-lunas", appHandler.kontrabonHandler.UpdateLunas)

	router.GET("/api/price/template/find", appHandler.priceHandler.Find)
	router.GET("/api/price/template/findDetail", appHandler.priceHandler.FindDetail)
	router.POST("/api/price/template/create", appHandler.priceHandler.Create)
	router.POST("/api/price/template/edit-price", appHandler.priceHandler.EditPrice)
	router.POST("/api/price/template/apply", appHandler.priceHandler.ApplyToCustomer)
	router.POST("/api/price/template/delete", appHandler.priceHandler.DeleteTemplate)
	router.POST("/api/price/template/copy", appHandler.priceHandler.CopyTemplate)
	router.POST("/api/price/template/download", appHandler.priceHandler.Download)

	router.GET("/api/price/buytemplate/find", appHandler.priceHandler.FindBuy)
	router.GET("/api/price/buytemplate/findDetail", appHandler.priceHandler.FindBuyDetail)
	router.POST("/api/price/buytemplate/create", appHandler.priceHandler.CreateBuy)
	router.POST("/api/price/buytemplate/edit-price", appHandler.priceHandler.EditBuyPrice)
	router.POST("/api/price/buytemplate/apply", appHandler.priceHandler.ApplyToTrx)
	router.POST("/api/price/buytemplate/delete", appHandler.priceHandler.DeleteBuyTemplate)
	router.POST("/api/price/buytemplate/copy", appHandler.priceHandler.CopyBuyTemplate)
	router.POST("/api/price/buytemplate/download", appHandler.priceHandler.DownloadBuy)

	router.GET("/api/mobile/dana/find", appHandler.transactionHandler.FindDana)
	router.POST("/api/mobile/dana/create", appHandler.transactionHandler.CreateDana)
	router.POST("/api/mobile/dana/update", appHandler.transactionHandler.UpdateDana)
	router.POST("/api/mobile/dana/send", appHandler.transactionHandler.SendDana)
	router.POST("/api/mobile/dana/approve", appHandler.transactionHandler.ApproveDana)
	router.POST("/api/mobile/dana/reject", appHandler.transactionHandler.RejectDana)
	router.POST("/api/mobile/dana/cancel", appHandler.transactionHandler.CancelSendDana)
	router.GET("/api/mobile/dana/find-user", appHandler.transactionHandler.FindUserMobile)

	router.POST("/api/mobile/penjualan/create", appHandler.transactionHandler.CreatePenjualan)
	router.POST("/api/mobile/penjualan/delete", appHandler.transactionHandler.DeletePenjualan)
	router.GET("/api/mobile/penjualan/find", appHandler.transactionHandler.FindPenjualan)

	router.POST("/api/mobile/belanja/create", appHandler.transactionHandler.CreateBelanja)
	router.POST("/api/mobile/belanja/delete", appHandler.transactionHandler.DeleteBelanja)
	router.GET("/api/mobile/belanja/find", appHandler.transactionHandler.FindBelanja)

	router.POST("/api/mobile/operasional/create", appHandler.transactionHandler.CreateOperasional)
	router.POST("/api/mobile/operasional/delete", appHandler.transactionHandler.DeleteOperasional)
	router.GET("/api/mobile/operasional/find", appHandler.transactionHandler.FindOperasional)
	router.GET("/api/mobile/operasional/find-description", appHandler.transactionHandler.FindDescriptionOperasional)

	router.GET("/api/mobile/saldo", appHandler.transactionHandler.FindSaldo)
	router.GET("/api/mobile/rekapitulasi", appHandler.transactionHandler.FindRekapitulasi)

	return router
}
