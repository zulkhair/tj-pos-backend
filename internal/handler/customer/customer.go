package customerhandler

import (
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	customerusecase "dromatech/pos-backend/internal/usecase/customer"
	restutil "dromatech/pos-backend/internal/util/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
)

// Handler defines the handler
type Handler struct {
	customerUsecase customerusecase.CustmerUsecase
}

func New(customerUsecase customerusecase.CustmerUsecase) *Handler {
	return &Handler{
		customerUsecase: customerUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")
	active := c.Query("active")

	var pointerBool *bool
	latestBool, err := strconv.ParseBool(active)
	if err != nil {
		pointerBool = nil
	} else {
		pointerBool = &latestBool
	}

	products, err := h.customerUsecase.Find(id, code, name, pointerBool)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
	}

	restutil.SendResponseOk(c, "", products)
}

func (h *Handler) FindActive(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")

	active := true

	products, err := h.customerUsecase.Find(id, code, name, &active)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
	}

	restutil.SendResponseOk(c, "", products)
}

func (h *Handler) Create(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	code := gjson.Get(string(jsonData), "code")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kode customer")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama customer")
		return
	}

	initialCredit := int64(0)
	initialCreditResult := gjson.Get(string(jsonData), "initialCredit")
	if initialCreditResult.Exists() || initialCreditResult.Int() > 0 {
		initialCredit = initialCreditResult.Int()
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.customerUsecase.Create(code.String(), name.String(), description.String(), initialCredit)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "customer berhasil ditambahkan", nil)
}

func (h *Handler) Edit(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	id := gjson.Get(string(jsonData), "id")
	if !id.Exists() || id.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih customer yang akan diperbarui")
		return
	}

	code := gjson.Get(string(jsonData), "code")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kode customer")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama customer")
		return
	}

	active := gjson.Get(string(jsonData), "active")
	if !active.Exists() {
		restutil.SendResponseFail(c, "Harap pilih status")
		return
	}

	initialCredit := int64(0)
	initialCreditResult := gjson.Get(string(jsonData), "initialCredit")
	if initialCreditResult.Exists() || initialCreditResult.Int() > 0 {
		initialCredit = initialCreditResult.Int()
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.customerUsecase.Edit(id.String(), code.String(), name.String(), description.String(), active.Bool(), initialCredit)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Customer berhasil diperbarui", nil)
}

func (h *Handler) GetSellPrice(c *gin.Context) {
	supplierId := c.Query("customerId")
	unitId := c.Query("unitId")
	date := c.Query("date")
	productId := c.Query("productId")

	response, err := h.customerUsecase.GetSellPrice(supplierId, unitId, date, productId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", response)
}

func (h *Handler) AddPrice(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat menambahkan data harga")
		return
	}

	request := customerdomain.AddPriceRequest{}
	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		logrus.Errorf(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat menambahkan data harga")
		return
	}

	userId := restutil.GetSession(c).UserID
	err = h.customerUsecase.AddSellPrice(request, userId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Data harga berhasil diperbarui", nil)
}

func (h *Handler) UpdateSellPrice(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat memperbarui harga")
		return
	}

	request := customerdomain.SellPriceRequest{}
	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		logrus.Errorf(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat memperbarui harga")
		return
	}

	err = h.customerUsecase.UpdateSellPrice(request)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Daftar harga berhasil diperbarui", nil)
}

func (h *Handler) FindLatestPrice(c *gin.Context) {
	customerId := c.Query("customerId")
	unitId := c.Query("unitId")

	latest := true
	response, err := h.customerUsecase.FindSellPrice(customerId, unitId, "", &latest)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", response)
}

func (h *Handler) FindPrice(c *gin.Context) {
	customerId := c.Query("customerId")
	unitId := c.Query("unitId")
	productId := c.Query("productId")
	latest := c.Query("latest")

	var pointerBool *bool
	latestBool, err := strconv.ParseBool(latest)
	if err != nil {
		pointerBool = nil
	} else {
		pointerBool = &latestBool
	}

	response, err := h.customerUsecase.FindSellPrice(customerId, unitId, productId, pointerBool)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", response)
}
