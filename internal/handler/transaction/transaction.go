package transactionhandler

import (
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	restutil "dromatech/pos-backend/internal/util/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type transactionUsecase interface {
	CreateTransaction(transaction *transactiondomain.Transaction) error
	ViewTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.Transaction, error)
	UpdateStatus(transactionID, status string) error
	UpdateBuyPrice(transactionID, unitID, productID string, price float64) error
}

// Handler defines the handler
type Handler struct {
	transactionUsecase transactionUsecase
}

func New(transactionUsecase transactionUsecase) *Handler {
	return &Handler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	code := c.Query("code")
	stakeholderID := c.Query("stakeholderId")
	txType := c.Query("txType")
	status := c.Query("status")
	productID := c.Query("productId")

	transactions, err := h.transactionUsecase.ViewTransaction(startDate, endDate, code, stakeholderID, txType, status, productID)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
	}

	restutil.SendResponseOk(c, "", transactions)
}

func (h *Handler) Create(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	transaction := &transactiondomain.Transaction{}
	err = json.Unmarshal(jsonData, transaction)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	if transaction.StakeholderID == "" {
		restutil.SendResponseFail(c, "Harap isi kode supplier")
		return
	}

	if transaction.TransactionType == "" {
		restutil.SendResponseFail(c, "Harap isi kode tipe transaksi")
		return
	}

	if transaction.ReferenceCode == "" {
		restutil.SendResponseFail(c, "Harap isi kode referensi")
		return
	}

	err = h.transactionUsecase.CreateTransaction(transaction)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Transaksi berhasil ditambahkan", nil)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	id := gjson.Get(string(jsonData), "transactionId")
	if !id.Exists() || id.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih transaksi yang akan diperbarui")
		return
	}

	status := gjson.Get(string(jsonData), "status")
	if !status.Exists() || status.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih status")
		return
	}

	err = h.transactionUsecase.UpdateStatus(id.String(), status.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Transaksi berhasil diperbarui", nil)
}

func (h *Handler) UpdateBuyPrice(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	price := gjson.Get(string(jsonData), "price")
	if !price.Exists() {
		restutil.SendResponseFail(c, "Harap pilih transaksi yang akan diperbarui")
		return
	}

	transactionId := gjson.Get(string(jsonData), "transactionId")
	if !transactionId.Exists() || transactionId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih transaksi")
		return
	}

	unitId := gjson.Get(string(jsonData), "unitId")
	if !unitId.Exists() || unitId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih satuan")
		return
	}

	productId := gjson.Get(string(jsonData), "productId")
	if !productId.Exists() || productId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih produk")
		return
	}

	err = h.transactionUsecase.UpdateBuyPrice(transactionId.String(), unitId.String(), productId.String(), price.Float())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Transaksi berhasil diperbarui", nil)
}
