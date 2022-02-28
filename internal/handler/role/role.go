package rolehandler

import (
	roledomain "dromatech/pos-backend/internal/domain/role"
	restutil "dromatech/pos-backend/internal/util/rest"
	"github.com/gin-gonic/gin"
)

type roleUsecase interface {
	GetActiveRole() ([]*roledomain.RoleResponseModel, error)
}

// Handler defines the handler
type Handler struct {
	roleusecase roleUsecase
}

func New(roleusecase roleUsecase) *Handler {
	return &Handler{
		roleusecase: roleusecase,
	}
}

func (h *Handler) GetActive(c *gin.Context) {
	roles, err := h.roleusecase.GetActiveRole()
	if err != nil {
		restutil.SendResponseFail(c, "Terjadi kesalahan saat pengambilan data Role")
		return
	}
	restutil.SendResponseOk(c, "", roles)
}
