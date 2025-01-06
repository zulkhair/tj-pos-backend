package webuserhandler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	webuserdomain "dromatech/pos-backend/internal/domain/webuser"
	restutil "dromatech/pos-backend/internal/util/rest"
)

type webuserUsecase interface {
	EditUser(userId, name, username, role, status string) error
	ChangePassword(userId, password1, password2 string) error
	RegisterUser(creatorId, name, username, password, roleId string) error
	FindAllUser() ([]*webuserdomain.WebUser, error)
	ForceChangePassword(userId, password string) error
	ChangeStatus(userId string, status bool)
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

func (h *Handler) EditName(c *gin.Context) {
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
	h.webuserUsecase.EditUser(session.UserID, name.String(), "", "", "")
	restutil.SendResponseOk(c, "Nama berhasil diubah", nil)
}

func (h *Handler) EditUser(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	userId := gjson.Get(string(jsonData), "userId")
	name := gjson.Get(string(jsonData), "name")
	username := gjson.Get(string(jsonData), "username")
	role := gjson.Get(string(jsonData), "role")
	active := gjson.Get(string(jsonData), "active")

	err = h.webuserUsecase.EditUser(userId.String(), name.String(), username.String(), role.String(), active.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}
	restutil.SendResponseOk(c, "Data berhasil diubah", nil)
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

func (h *Handler) FindAllUser(c *gin.Context) {
	users, err := h.webuserUsecase.FindAllUser()
	if err != nil {
		logrus.Error(err.Error())
	}
	restutil.SendResponseOk(c, "", users)
}

func (h *Handler) ForceChangePassword(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	userId := gjson.Get(string(jsonData), "userId")
	if !userId.Exists() || userId.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap pilih user yang akan diedit", nil))
		return
	}

	password1 := gjson.Get(string(jsonData), "password1")
	password2 := gjson.Get(string(jsonData), "password2")
	if !password1.Exists() || password1.String() == "" {
		restutil.SendResponseFail(c, "Harap isi kata sandi baru")
		return
	}
	if !password2.Exists() || password2.String() == "" {
		restutil.SendResponseFail(c, "Harap isi ulangi kata sandi baru")
		return
	}
	if password2.String() != password1.String() {
		restutil.SendResponseFail(c, "Ulangi kata sandi baru harus sesuai")
		return
	}

	//session := restutil.GetSession(c)
	err = h.webuserUsecase.ForceChangePassword(userId.String(), password1.String())
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}
	restutil.SendResponseOk(c, "Kata sandi berhasil diubah", nil)
}

func (h *Handler) ChangeStatus(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	userId := gjson.Get(string(jsonData), "userId")
	if !userId.Exists() || userId.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap pilih user yang akan diedit", nil))
		return
	}

	active := gjson.Get(string(jsonData), "active")
	if !active.Exists() {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap isi status", nil))
		return
	}

	//session := restutil.GetSession(c)
	h.webuserUsecase.ChangeStatus(userId.String(), active.Bool())
	restutil.SendResponseOk(c, "Status berhasil diubah", nil)
}
