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
	EditUser(userId, name string)
	ChangePassword(userId, password1, password2 string) error
	RegisterUser(creatorId, name, username, password, roleId string) error
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
		restutil.SendResponseFail(c, "Harap isi kata sandi lama")
		return
	}
	if !password2.Exists() || password2.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kata sandi baru")
		return
	}
	if !password3.Exists() || password3.String() == "" {
		restutil.SendResponseFail(c, "Harap isi ulangi kata sandi baru")
		return
	}
	if password2.String() != password3.String() {
		restutil.SendResponseFail(c, "Ulangi kata sandi baru harus sesuai")
		return
	}

	session := restutil.GetSession(c)
	err = h.webuserUsecase.ChangePassword(session.UserID, password1.String(), password2.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}
	restutil.SendResponseOk(c, "Kata sandi berhasil diubah", nil)
}

func (h *Handler) RegisterUser(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	name := gjson.Get(string(jsonData), "name")
	username := gjson.Get(string(jsonData), "username")
	password := gjson.Get(string(jsonData), "password")
	roleId := gjson.Get(string(jsonData), "roleId")

	if !name.Exists() || name.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama")
		return
	}
	if !username.Exists() || username.String() == "" {
		restutil.SendResponseFail(c, "Harap isi username")
		return
	}
	if !password.Exists() || password.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kata sandi")
		return
	}
	if !roleId.Exists() || roleId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih role")
		return
	}

	session := restutil.GetSession(c)
	err = h.webuserUsecase.RegisterUser(session.UserID, name.String(), username.String(), password.String(), roleId.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}
	restutil.SendResponseOk(c, "Pengguna baru berhasil ditambahkan", nil)
}
