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
	"/auth/login":   true,
	"/auth/logout":  true,
	"/auth/getmenu": true,
	"/ping":         true,
}

func New(sessionUsecase sessionUsecase) *Handler {
	return &Handler{
		sessionUc: sessionUsecase,
	}
}

func (h *Handler) AuthCheck(c *gin.Context) {
	path := c.FullPath()
	if isWhitelistedPath(path) {
		c.Next()
		return
	}

	token := c.GetHeader("token")
	url, status, session := h.sessionUc.AuthCheck(token, path)
	if status == 200 {
		c.Set("session", session)
		c.Next()
		return
	} else {
		c.Redirect(status, url)
	}
}

func (h *Handler) Login(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	username := gjson.Get(string(jsonData), "username")
	if !username.Exists() || username.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponseJson(1, "Harap isi username", nil))
		return
	}
	password := gjson.Get(string(jsonData), "password")
	if !password.Exists() || password.String() == "" {
		c.JSON(http.StatusOK, restutil.CreateResponseJson(1, "Harap isi password", nil))
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
	token := c.Param("token")
	h.sessionUc.Logout(token)
}

func (h *Handler) GetMenu(c *gin.Context) {
	token := c.GetHeader("token")
	if token == "" {
		restutil.RedirectToUnuthorized(c)
		return
	}

	session := h.sessionUc.GetSession(token)
	if session == nil {
		restutil.RedirectToUnuthorized(c)
		return
	}

	c.JSON(http.StatusOK, restutil.CreateResponseOk(session))
}

func isWhitelistedPath(path string) bool {
	if _, exists := WhitelistPath[path]; exists {
		return true
	}

	return false
}
