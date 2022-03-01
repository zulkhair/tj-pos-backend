package rolehandler

import (
	roledomain "dromatech/pos-backend/internal/domain/role"
	restutil "dromatech/pos-backend/internal/util/rest"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type roleUsecase interface {
	GetActiveRole() ([]*roledomain.RoleResponseModel, error)
	FindPermissions(roleId string) ([]*roledomain.Permission, error)
	RegisterRole(roleName string, permissions []string) error
	GetAllRole() ([]*roledomain.Role, error)
	EditRole(roleId string, roleName string, active bool, permissions []string) error
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

func (h *Handler) GetAll(c *gin.Context) {
	roles, err := h.roleusecase.GetAllRole()
	if err != nil {
		restutil.SendResponseFail(c, "Terjadi kesalahan saat pengambilan data Role")
		return
	}
	restutil.SendResponseOk(c, "", roles)
}

func (h *Handler) FindPermissions(c *gin.Context) {
	roleId := c.Query("roleId")
	permissions, err := h.roleusecase.FindPermissions(roleId)
	if err != nil {
		restutil.SendResponseFail(c, "Terjadi kesalahan saat pengambilan data Role")
		return
	}
	restutil.SendResponseOk(c, "", permissions)
}

func (h *Handler) RegisterRole(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	roleName := gjson.Get(string(jsonData), "roleName")
	permissionArray := gjson.Get(string(jsonData), "permissions")

	if !roleName.Exists() || roleName.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama role")
		return
	}
	if !permissionArray.Exists() || len(permissionArray.Array()) <= 0 {
		restutil.SendResponseFail(c, "Harap pilih ability")
		return
	}

	var permissions []string
	for _, p := range permissionArray.Array() {
		permissions = append(permissions, p.String())
	}

	err = h.roleusecase.RegisterRole(roleName.String(), permissions)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}
	restutil.SendResponseOk(c, "Role baru berhasil ditambahkan", nil)
}

func (h *Handler) EditRole(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("bad request"))
	}

	roleId := gjson.Get(string(jsonData), "roleId")
	roleName := gjson.Get(string(jsonData), "roleName")
	active := gjson.Get(string(jsonData), "active")
	permissionArray := gjson.Get(string(jsonData), "permissions")

	if !roleId.Exists() || roleId.String() == "" {
		restutil.SendResponseFail(c, "Harap pilih role terlebih dahulu")
		return
	}
	if !roleName.Exists() || roleName.String() == "" {
		restutil.SendResponseFail(c, "Harap isi nama role")
		return
	}
	if !active.Exists() {
		restutil.SendResponseFail(c, "Harap pilih status")
		return
	}
	if !permissionArray.Exists() || len(permissionArray.Array()) <= 0 {
		restutil.SendResponseFail(c, "Harap pilih ability")
		return
	}

	var permissions []string
	for _, p := range permissionArray.Array() {
		permissions = append(permissions, p.String())
	}

	err = h.roleusecase.EditRole(roleId.String(), roleName.String(), active.Bool(), permissions)
	if err != nil {
		restutil.SendResponseFail(c, err.Error())
		return
	}
	restutil.SendResponseOk(c, "Role berhasil diperbarui", nil)
}
