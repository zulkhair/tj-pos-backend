package unithandler

import (
	unitdomain "dromatech/pos-backend/internal/domain/unit"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strconv"
)

type unitUsecase interface {
	Find(id, code string, active *bool) ([]*unitdomain.Unit, error)
	Create(code, description string) error
	Edit(id, code, description string, active *bool) error
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

	products, err := h.unitUsecase.Find(id, code, activeBool)
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
		restutil.SendResponseFail(c, "Harap isi kode satuan")
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

	var activeBool *bool
	active := gjson.Get(string(jsonData), "active")
	if active.Exists() {
		activeAddress := active.Bool()
		activeBool = &activeAddress
	}

	code := gjson.Get(string(jsonData), "code")
	description := gjson.Get(string(jsonData), "description")

	err = h.unitUsecase.Edit(id.String(), code.String(), description.String(), activeBool)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}

	restutil.SendResponseOk(c, "Unit berhasil diperbarui", nil)
}
