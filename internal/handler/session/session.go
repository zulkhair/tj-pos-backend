package sessionhandler

import (
	sessiondomain "dromatech/pos-backend/internal/domain/session"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type sessionUsecase interface {
	Login(username string, password string) (*sessiondomain.Session, error)
	Logout(token string)
	AuthCheck(token string, requestorPath string) (string, int, *sessiondomain.Session)
	GetSession(token string) *sessiondomain.Session
}

// Handler defines the handler
type Handler struct {
	sessionUc sessionUsecase
}

var WhitelistPath = map[string]bool{
	"/api/auth/login":           true,
	"/api/auth/logout":          true,
	"/api/auth/getmenu":         true,
	"/api/user/edit":            true,
	"/api/user/change-password": true,
	"/api/ping":                 true,
}

func New(sessionUsecase sessionUsecase) *Handler {
	return &Handler{
		sessionUc: sessionUsecase,
	}
}

func (h *Handler) AuthCheck(c *gin.Context) {
	path := c.FullPath()
	token := c.GetHeader("token")
	if isWhitelistedPath(path) {
		if token != "" {
			session := h.sessionUc.GetSession(token)
			if session != nil {
				restutil.SetSession(c, session)
			}
		}
		c.Next()
		return
	}

	_, status, session := h.sessionUc.AuthCheck(token, path)
	if status == 200 {
		restutil.SetSession(c, session)
		c.Next()
		return
	} else {
		c.AbortWithStatus(status)
	}
}

func (h *Handler) CheckPermission(c *gin.Context) {
	permission := c.Query("permission")
	if permission == "" {
		restutil.SendResponseFail(c, "")
		return
	}
	session := restutil.GetSession(c)
	for _, menu := range session.Menu {
		for _, subMenu := range menu.SubMenu {
			for _, per := range subMenu.Permissions {
				if per == permission {
					restutil.SendResponseOk(c, "", nil)
					return
				}
			}
		}
	}
	restutil.SendResponseFail(c, "")
}

func (h *Handler) Login(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	username := gjson.Get(string(jsonData), "username")
	if !username.Exists() || username.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap isi username", nil))
		return
	}
	password := gjson.Get(string(jsonData), "password")
	if !password.Exists() || password.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap isi password", nil))
		return
	}

	session, err := h.sessionUc.Login(username.String(), password.String())
	if err != nil {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, restutil.CreateResponse(0, "", session))
}

func (h *Handler) Logout(c *gin.Context) {
	token := c.GetHeader("token")
	h.sessionUc.Logout(token)
	restutil.SendResponseOk(c, "Berhasil logout", nil)
}

func (h *Handler) GetMenu(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap melakukan login terlebih dahulu", nil))
		return
	}

	session := h.sessionUc.GetSession(token)
	if session == nil {
		c.JSON(http.StatusOK, restutil.CreateResponse(1, "Harap melakukan login terlebih dahulu", nil))
		return
	}

	restutil.SendResponseOk(c, "", session)
}

func isWhitelistedPath(path string) bool {
	if _, exists := WhitelistPath[path]; exists {
		return true
	}

	return false
}
