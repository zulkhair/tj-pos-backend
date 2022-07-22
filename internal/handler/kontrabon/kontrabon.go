package kontrabonhandler

import (
	kontrabondomain "dromatech/pos-backend/internal/domain/kontrabon"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	kontrabonusecase "dromatech/pos-backend/internal/usecase/kontrabon"
	restutil "dromatech/pos-backend/internal/util/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// Handler defines the handler
type Handler struct {
	kontrabonUsecase kontrabonusecase.KontrabonUsecase
}

func New(kontrabonUsecase kontrabonusecase.KontrabonUsecase) *Handler {
	return &Handler{
		kontrabonUsecase: kontrabonUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	code := c.Query("code")

	kontrabons, err := h.kontrabonUsecase.Find(code, startDate, endDate)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", kontrabons)
}

func (h *Handler) FindTransaction(c *gin.Context) {
	kontrabonId := c.Query("kontrabonId")

	kontrabons, err := h.kontrabonUsecase.FindTransaction(kontrabonId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", kontrabons)
}

func (h *Handler) Create(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	var kontrabon kontrabondomain.CreateRequest

	err = json.Unmarshal(jsonData, &kontrabon)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if len(kontrabon.TransactionIDs) <= 0 {
		restutil.SendResponseFail(c, "Harap pilih transaksi yang akan ditambahkan ke kontrabon")
		return
	}

	err = h.kontrabonUsecase.Create(kontrabon.CustomerID, kontrabon.TransactionIDs)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Kontrabon berhasil ditambahkan",nil)
}

func (h *Handler) Add(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	request := &kontrabondomain.UpdateRequest{}

	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if request.KontrabonID == ""{
		restutil.SendResponseFail(c, "Harap pilih kontrabon yang akan diperbarui")
		return
	}

	if len(request.TransactionIDs) <= 0{
		restutil.SendResponseFail(c, "Harap pilih transaksi yang akan ditambahkan")
		return
	}

	err = h.kontrabonUsecase.Update(request.KontrabonID, request.TransactionIDs, transactiondomain.TRANSACTION_KONTRABON)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Kontrabon berhasil diperbarui",nil)
}

func (h *Handler) Remove(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	request := &kontrabondomain.UpdateRequest{}

	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if request.KontrabonID == ""{
		restutil.SendResponseFail(c, "Harap pilih kontrabon yang akan diperbarui")
		return
	}

	if len(request.TransactionIDs) <= 0{
		restutil.SendResponseFail(c, "Harap pilih transaksi yang akan ditambahkan")
		return
	}

	err = h.kontrabonUsecase.Update(request.KontrabonID, request.TransactionIDs, transactiondomain.TRANSACTION_STATUS_PEMBUATAN)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Kontrabon berhasil diperbarui",nil)
}

func (h *Handler) UpdateLunas(c *gin.Context) {
	kontrabonID, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if string(kontrabonID) == ""{
		restutil.SendResponseFail(c, "Harap pilih kontrabon yang akan diperbarui")
		return
	}

	err = h.kontrabonUsecase.UpdateLunas(string(kontrabonID))
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Status kontrabon berhasil diperbarui",nil)
}