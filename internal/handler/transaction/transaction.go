package transactionhandler

import (
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	transactionusecase "dromatech/pos-backend/internal/usecase/transaction"
	restutil "dromatech/pos-backend/internal/util/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

// Handler defines the handler
type Handler struct {
	transactionUsecase transactionusecase.TransactoionUsecase
}

func New(transactionUsecase transactionusecase.TransactoionUsecase) *Handler {
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

	transactions, err := h.transactionUsecase.ViewSellTransaction(startDate, endDate, code, stakeholderID, txType, status, productID)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", transactions)
}

func (h *Handler) Create(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	transaction := &transactiondomain.Transaction{}
	err = json.Unmarshal(jsonData, transaction)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if transaction.StakeholderID == "" {
		restutil.SendResponseFail(c, "Harap isi kode stakeholder")
		return
	}

	if transaction.TransactionType == "" {
		restutil.SendResponseFail(c, "Harap isi kode tipe transaksi")
		return
	}

	transaction.UserId = restutil.GetSession(c).UserID
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
		return
	}

	id := gjson.Get(string(jsonData), "transactionId")
	if !id.Exists() || id.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih transaksi yang akan diperbarui")
		return
	}

	err = h.transactionUsecase.UpdateStatus(id.String())
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
		return
	}

	buyPrice := gjson.Get(string(jsonData), "buyPrice")
	if !buyPrice.Exists() {
		restutil.SendResponseFail(c, "Harap isi harga beli")
		return
	}

	sellPrice := gjson.Get(string(jsonData), "sellPrice")
	if !sellPrice.Exists() {
		restutil.SendResponseFail(c, "Harap isi harga jual")
		return
	}

	quantity := gjson.Get(string(jsonData), "quantity")
	if !buyPrice.Exists() {
		restutil.SendResponseFail(c, "Harap isi jumlah jual")
		return
	}

	buyQuantity := gjson.Get(string(jsonData), "buy_quantity")
	if !sellPrice.Exists() {
		restutil.SendResponseFail(c, "Harap isi harga beli")
		return
	}

	transactionId := gjson.Get(string(jsonData), "transactionId")
	if !transactionId.Exists() || transactionId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih transaksi")
		return
	}

	productId := gjson.Get(string(jsonData), "productId")
	if !productId.Exists() || productId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih produk")
		return
	}

	err = h.transactionUsecase.UpdateBuyPrice(transactionId.String(), productId.String(), buyPrice.Float(), sellPrice.Float(), quantity.Int(), buyQuantity.Int())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Transaksi berhasil diperbarui", nil)
}
