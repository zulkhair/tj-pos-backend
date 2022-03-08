package customerhandler

import (
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type customerUsecase interface {
	Find(id, code, name string) ([]*customerdomain.Customer, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
}

// Handler defines the handler
type Handler struct {
	customerUsecase customerUsecase
}

func New(customerUsecase customerUsecase) *Handler {
	return &Handler{
		customerUsecase: customerUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")

	products, err := h.customerUsecase.Find(id, code, name)
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
		restutil.SendResponseFail(c, "Harap isi kode customer")
		return
	}

	name := gjson.Get(string(jsonData), "name")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama customer")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.customerUsecase.Create(code.String(), name.String(), description.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "customer berhasil ditambahkan", nil)
}

func (h *Handler) Edit(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
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

	description := gjson.Get(string(jsonData), "description")

	err = h.customerUsecase.Edit(id.String(), code.String(), name.String(), description.String(), active.Bool())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Customer berhasil diperbarui", nil)
}
