package supplierhandler

import (
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type supplierUsecase interface {
	Find(id, code, name string) ([]*supplierdomain.Supplier, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
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

	products, err := h.supplierUsecase.Find(id, code, name)
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
