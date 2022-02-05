package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sessionUsecase interface {
	CheckTokenAuth(ctx context.Context, r *http.Request) (string, int)
	Logout(ctx context.Context, r *http.Request) error
}

// Handler defines the handler
type Handler struct {
	sessionUc sessionUsecase
}

var WhitelistPath = map[string]bool{
	"/api/auth":        true,
	"/ping":            true,
	"/api/auth/logout": true,
}

func New(sessionUsecase sessionUsecase) *Handler {
	return &Handler{
		sessionUc: sessionUsecase,
	}
}

func (h *Handler) CheckTokenAPI(c *gin.Context) {
	if isWhitelistedPath(c.FullPath()) {
		c.Next()
		return
	}

	token := c.Param("udtj_token")
	if token == "" {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
	}

	// 0 = OK
	// 1 = Token id is empty, need signin
	// 2 = Not registered in database
	_, res := h.sessionUc.CheckTokenAuth(ctx, r)
	if res == 0 {
		next(w, r)
	} else if res == 1 || res == 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else if res == 3 {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func (h *Handler) CheckTokenAuth(w http.ResponseWriter, r *http.Request) {
	var (
		span, ctx = tracer.StartFromRequest(r)
	)
	defer span.Finish()

	// 0 = OK
	// 1 = Token id is empty, need signin
	// 2 = Failed get info to google
	// 3 = Not registered in database
	redirectUrl, flag := h.sessionUc.CheckTokenAuth(ctx, r)
	switch flag {
	case 0:
		// OK continue
		if _, err := response.WriteJSONAPIData(w, r, http.StatusOK, ""); err != nil {
			h.commonLog.Errorf("[session.CheckAuthToken] error from WriteJSON: ", err)
		}
	case 1:
		// empty token redirect to login page
		if _, err := response.WriteJSONAPIData(w, r, http.StatusUnauthorized, redirectUrl); err != nil {
			h.commonLog.Errorf("[session.CheckAuthToken] error from WriteJSON: ", err)
		}
	case 2:
		// failed get info to google redirect to login page
		if _, err := response.WriteJSONAPIData(w, r, http.StatusUnauthorized, redirectUrl); err != nil {
			h.commonLog.Errorf("[session.CheckAuthToken] error from WriteJSON: ", err)
		}
	case 3:
		// forbidden redirect forbidden page
		if _, err := response.WriteJSONAPIData(w, r, http.StatusForbidden, redirectUrl); err != nil {
			h.commonLog.Errorf("[session.CheckAuthToken] error from WriteJSON: ", err)
		}
	default:
		// default redirect to login page
		if _, err := response.WriteJSONAPIData(w, r, http.StatusUnauthorized, redirectUrl); err != nil {
			h.commonLog.Errorf("[session.CheckAuthToken] error from WriteJSON: ", err)
		}
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var (
		span, ctx = tracer.StartFromRequest(r)
	)
	defer span.Finish()

	h.sessionUc.Logout(ctx, r)
}

func isWhitelistedPath(path string) bool {
	if _, exists := WhitelistPath[path]; exists {
		return true
	}

	return false
}
