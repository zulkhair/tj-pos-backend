package restutil

import (
	"dromatech/pos-backend/global"
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

func CreateResponseOk(data interface{}) Response {
	return CreateResponse(0, "", data)
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
