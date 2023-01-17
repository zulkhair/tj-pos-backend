package pricehandler

import (
	pricedomain "dromatech/pos-backend/internal/domain/price"
	priceusecase "dromatech/pos-backend/internal/usecase/price"
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
	priceUsecase priceusecase.PriceUsecase
}

func New(priceUsecase priceusecase.PriceUsecase) *Handler {
	return &Handler{
		priceUsecase: priceUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	name := c.Query("name")

	prices, err := h.priceUsecase.Find(name)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", prices)
}

func (h *Handler) FindDetail(c *gin.Context) {
	templateId := c.Query("templateId")

	prices, err := h.priceUsecase.FindDetail(templateId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", prices)
}

func (h *Handler) Create(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	name := gjson.Get(string(jsonData), "name")
	if !name.Exists() || name.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama template")
		return
	}

	err = h.priceUsecase.Create(name.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Template berhasil ditambahkan", nil)
}

func (h *Handler) EditPrice(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	price := gjson.Get(string(jsonData), "price")
	if !price.Exists() {
		restutil.SendResponseFail(c, "Harap isi harga produk")
		return
	}

	templateId := gjson.Get(string(jsonData), "templateId")
	if !templateId.Exists() || templateId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih template")
		return
	}

	productId := gjson.Get(string(jsonData), "productId")
	if !productId.Exists() || productId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih produk")
		return
	}

	err = h.priceUsecase.EditPrice(templateId.String(), productId.String(), price.Float())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Harga berhasil diubah", nil)
}

func (h *Handler) ApplyToCustomer(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	request := &pricedomain.ApplyToCustomerReq{}
	err = json.Unmarshal(jsonData, request)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if request.TemplateID == "" {
		restutil.SendResponseFail(c, "Harap pilih template")
		return
	}

	if len(request.CustomerIDs) == 0 {
		restutil.SendResponseFail(c, "Harap pilih customer")
		return
	}

	userId := restutil.GetSession(c).UserID
	err = h.priceUsecase.ApplyToCustomer(request.TemplateID, request.CustomerIDs, userId)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Harga berhasil diterapkan", nil)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	request := &pricedomain.DeleteTemplateReq{}
	err = json.Unmarshal(jsonData, request)
	if request.TemplateID != "" {
		h.priceUsecase.DeleteTemplate(request.TemplateID)
		restutil.SendResponseOk(c, "Template berhasil dihapus", nil)
		return
	}
	restutil.SendResponseFail(c, "Harap pilih template")
}

func (h *Handler) CopyTemplate(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	templateId := gjson.Get(string(jsonData), "templateId")
	if !templateId.Exists() || templateId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih template")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !name.Exists() || name.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama template")
		return
	}

	err = h.priceUsecase.CopyTemplate(templateId.String(), name.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Template berhasil diduplikasi", nil)
}

func (h *Handler) Download(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	request := &pricedomain.Download{}
	err = json.Unmarshal(jsonData, request)
	if err != nil {
		logrus.Error(err.Error())
		c.AbortWithError(400, fmt.Errorf("bad request"))
		return
	}

	if request != nil && request.TemplateDetailIDs != nil && len(request.TemplateDetailIDs) > 0 {
		h.priceUsecase.Download(*request)
		restutil.SendResponseOk(c, "Template berhasil diunduh", nil)
		return
	}
	restutil.SendResponseFail(c, "Harap pilih template")
}