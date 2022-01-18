package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

// New creates profile handler
func New() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, map[string]interface{}{"status": "system is running", "msg": "hello world"})
}
