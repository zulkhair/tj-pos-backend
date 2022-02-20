package webuserhandler

import (
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type webuserUsecase interface {
	EditUser(userId string, name string)
	ChangePassword(userId string, password1 string, password2 string) error
}

// Handler defines the handler
type Handler struct {
	webuserUsecase webuserUsecase
}

func New(webuserUsecase webuserUsecase) *Handler {
	return &Handler{
		webuserUsecase: webuserUsecase,
	}
}

func (h *Handler) EditUser(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	name := gjson.Get(string(jsonData), "name")
	if !name.Exists() || name.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap isi nama", nil))
		return
	}

	session := restutil.GetSession(c)
	h.webuserUsecase.EditUser(session.UserID, name.String())
	restutil.SendResponseOk(c, "Nama berhasil diubah", nil)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	password1 := gjson.Get(string(jsonData), "password1")
	password2 := gjson.Get(string(jsonData), "password2")
	password3 := gjson.Get(string(jsonData), "password3")
	if !password1.Exists() || password1.String() == "" {
		restutil.SendResponseOk(c, "Harap isi kata sandi lama", nil)
		return
	}
	if !password2.Exists() || password2.String() == "" {
		restutil.SendResponseOk(c, "Harap isi kata sandi baru", nil)
		return
	}
	if !password3.Exists() || password3.String() == "" {
		restutil.SendResponseOk(c, "Harap isi ulangi kata sandi baru", nil)
		return
	}
	if password2.String() != password3.String() {
		restutil.SendResponseOk(c, "Ulangi kata sandi baru harus sesuai", nil)
		return
	}

	session := restutil.GetSession(c)
	h.webuserUsecase.ChangePassword(session.UserID, password1.String(), password2.String())
	restutil.SendResponseOk(c, "Kata sandi berhasil diubah", nil)
}
