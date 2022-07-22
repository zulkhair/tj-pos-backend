package producthandler

import (
	productusecase "dromatech/pos-backend/internal/usecase/product"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strconv"
)

// Handler defines the handler
type Handler struct {
	productUsecase productusecase.ProductUsecase
}

func New(productUsecase productusecase.ProductUsecase) *Handler {
	return &Handler{
		productUsecase: productUsecase,
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

	products, err := h.productUsecase.Find(id, code, name, activeBool)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "", products)
}

func (h *Handler) FindActive(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")

	var activeBool = true

	products, err := h.productUsecase.Find(id, code, name, &activeBool)
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
		restutil.SendResponseFail(c, "Harap isi kode produk")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama produk")
		return
	}

	unitId := gjson.Get(string(jsonData), "unitId")
	if !unitId.Exists() || unitId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih satuan")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.productUsecase.Create(code.String(), name.String(), description.String(), unitId.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Produk berhasil ditambahkan", nil)
}

func (h *Handler) Edit(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	id := gjson.Get(string(jsonData), "id")
	if !id.Exists() || id.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih produk yang akan diperbarui")
		return
	}

	code := gjson.Get(string(jsonData), "code")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kode produk")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama produk")
		return
	}

	active := gjson.Get(string(jsonData), "active")
	if !active.Exists() {
		restutil.SendResponseFail(c, "Harap pilih status")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.productUsecase.Edit(id.String(), code.String(), name.String(), description.String(), active.Bool())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Produk berhasil diperbarui", nil)
}
