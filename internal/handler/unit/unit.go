package unithandler

import (
	unitdomain "dromatech/pos-backend/internal/domain/unit"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type unitUsecase interface {
	Find(id, code string) ([]*unitdomain.Unit, error)
	Create(code, description string) error
	Edit(id, code, description string) error
}

// Handler defines the handler
type Handler struct {
	unitUsecase unitUsecase
}

func New(unitUsecase unitUsecase) *Handler {
	return &Handler{
		unitUsecase: unitUsecase,
	}
}

func (h *Handler) Find(c *gin.Context) {
	id := c.Query("id")
	code := c.Query("code")

	products, err := h.unitUsecase.Find(id, code)
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
		restutil.SendResponseFail(c, "Harap isi kode unit")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.unitUsecase.Create(code.String(), description.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Unit berhasil ditambahkan", nil)
}

func (h *Handler) Edit(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	id := gjson.Get(string(jsonData), "id")
	if !id.Exists() || id.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih unit yang akan diperbarui")
		return
	}

	code := gjson.Get(string(jsonData), "code")
	if !code.Exists() || code.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kode unit")
		return
	}

	description := gjson.Get(string(jsonData), "description")

	err = h.unitUsecase.Edit(id.String(), code.String(), description.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Unit berhasil diperbarui", nil)
}
