package supplierhandler

import (
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	restutil "dromatech/pos-backend/internal/util/rest"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strconv"
)

type supplierUsecase interface {
	Find(id, code, name string, active *bool) ([]*supplierdomain.Supplier, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
	GetBuyPrice(supplierId, unitId, date, productId string) ([]*supplierdomain.BuyPriceResponse, error)
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
	AddBuyPrice(entity supplierdomain.AddPriceRequest, userId string) error
	FindBuyPrice(customerId, unitId, productId string, latest *bool) ([]*supplierdomain.PriceResponse, error)
}

// Handler defines the handler
type Handler struct {
	supplierUsecase supplierUsecase
}

func New(supplierUsecase supplierUsecase) *Handler {
	return &Handler{
		supplierUsecase: supplierUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")
	active := c.Query("active")

	var activeBool *bool
	if active != "" {
		parsedBool, err := strconv.ParseBool(active)
		if err != nil {
			logrus.Error(err.Error())
			activeBool = nil
		} else {
			activeBool = &parsedBool
		}
	}

	products, err := h.supplierUsecase.Find(id, code, name, activeBool)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
	}

	restutil.SendResponseOk(c, "", products)
}

func (h *Handler) Create(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	code := gjson.Get(string(jsonData), "code")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kode supplier")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama supplier")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.supplierUsecase.Create(code.String(), name.String(), description.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Supplier berhasil ditambahkan", nil)
}

func (h *Handler) Edit(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	id := gjson.Get(string(jsonData), "id")
	if !id.Exists() || id.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih supplier yang akan diperbarui")
		return
	}

	code := gjson.Get(string(jsonData), "code")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kode supplier")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama supplier")
		return
	}

	active := gjson.Get(string(jsonData), "active")
	if !active.Exists() {
		restutil.SendResponseFail(c, "Harap pilih status")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.supplierUsecase.Edit(id.String(), code.String(), name.String(), description.String(), active.Bool())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "supplier berhasil diperbarui", nil)
}

func (h *Handler) GetBuyPrice(c *gin.Context) {
	supplierId := c.Query("supplierId")
	unitId := c.Query("unitId")
	date := c.Query("date")
	productId := c.Query("productId")

	response, err := h.supplierUsecase.GetBuyPrice(supplierId, unitId, date, productId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", response)
}

func (h *Handler) UpdateBuyPrice(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat memperbarui harga")
		return
	}

	request := supplierdomain.BuyPriceRequest{}
	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		logrus.Errorf(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat memperbarui harga")
		return
	}

	err = h.supplierUsecase.UpdateBuyPrice(request)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Daftar harga berhasil diperbarui", nil)
}

func (h *Handler) AddPrice(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat menambahkan data harga")
		return
	}

	request := supplierdomain.AddPriceRequest{}
	err = json.Unmarshal(jsonData, &request)
	if err != nil {
		logrus.Errorf(err.Error())
		restutil.SendResponseFail(c, "Ada kesalahan saat menambahkan data harga")
		return
	}

	userId := restutil.GetSession(c).UserID
	err = h.supplierUsecase.AddBuyPrice(request, userId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Data harga berhasil diperbarui", nil)
}

func (h *Handler) FindLatestPrice(c *gin.Context) {
	supplierId := c.Query("supplierId")
	unitId := c.Query("unitId")

	latest := true
	response, err := h.supplierUsecase.FindBuyPrice(supplierId, unitId, "", &latest)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", response)
}

func (h *Handler) FindPrice(c *gin.Context) {
	supplierId := c.Query("supplierId")
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

	response, err := h.supplierUsecase.FindBuyPrice(supplierId, unitId, productId, pointerBool)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", response)
}
