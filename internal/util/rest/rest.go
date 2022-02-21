package restutil

import (
	"dromatech/pos-backend/global"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateResponse(status int, message string, data interface{}) Response {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	return response
}

func CreateResponseJson(status int, message string, data interface{}) string {
	json, _ := json.Marshal(CreateResponse(status, message, data))
	return string(json)
}

func CreateResponseOk(msg string, data interface{}) Response {
	return CreateResponse(0, msg, data)
}

func SendResponseOk(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, CreateResponseOk(msg, data))
}

func SendResponseFail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, CreateResponse(1, msg, nil))
}

func RedirectToLogin(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, global.LOGIN_URL)
}

func RedirectToUnuthorized(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, global.UNAUTHORIZED_URL)
}

func RedirectToForbidden(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, global.FORBIDDEN_URL)
}

func CloseContext(c *gin.Context) {
	if !c.IsAborted() {
		c.Abort()
	}
}

func SetSession(c *gin.Context, session *sessiondomain.Session) {
	c.Set("session", session)
}

func GetSession(c *gin.Context) *sessiondomain.Session {
	if session, ok := c.Get("session"); ok {
		return session.(*sessiondomain.Session)
	}

	return nil
}
