package roleusecase

import (
	roledomain "dromatech/pos-backend/internal/domain/role"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
	"fmt"
)

type RoleUsecase interface {
	GetActiveRole() ([]*roledomain.RoleResponseModel, error)
	FindPermissions(roleId string) ([]*roledomain.Permission, error)
	RegisterRole(roleName string, permissions []string) error
	GetAllRole() ([]*roledomain.Role, error)
	EditRole(roleId string, roleName string, active bool, permissions []string) error
}

type Usecase struct {
	rolerepo roleRepo
}

type roleRepo interface {
	Find(id string) *roledomain.Role
	FindMenu(roleId string) ([]*sessiondomain.Menu, error)
	FindAll() ([]*roledomain.Role, error)
	FindActive() ([]*roledomain.RoleResponseModel, error)
	FindActivePermissionPaths(roleId string) ([]string, error)
	FindPermissions() ([]*roledomain.Permission, error)
	FindPermissionsByRoleId(roleId string) ([]*roledomain.Permission, error)
	RegisterRole(name string, permissions []string)
	FindByName(name string) *roledomain.Role
	EditRole(roleId string, name string, active bool, permissions []string)
}

func New(rolerepo roleRepo) *Usecase {
	uc := &Usecase{
		rolerepo: rolerepo,
	}

	return uc
}

func (uc *Usecase) GetActiveRole() ([]*roledomain.RoleResponseModel, error) {
	return uc.rolerepo.FindActive()
}

func (uc *Usecase) GetAllRole() ([]*roledomain.Role, error) {
	return uc.rolerepo.FindAll()
}

func (uc *Usecase) FindPermissions(roleId string) ([]*roledomain.Permission, error) {
	if roleId != "" {
		return uc.rolerepo.FindPermissionsByRoleId(roleId)
	}
	return uc.rolerepo.FindPermissions()
}

func (uc *Usecase) RegisterRole(roleName string, permissions []string) error {
	role := uc.rolerepo.FindByName(roleName)
	if role != nil {
		return fmt.Errorf("Role dengan nama %s sudah terdaftar", roleName)
	}

	uc.rolerepo.RegisterRole(roleName, permissions)
	return nil
}

func (uc *Usecase) EditRole(roleId string, roleName string, active bool, permissions []string) error {
	role := uc.rolerepo.FindByName(roleName)
	if role != nil && role.ID != roleId {
		return fmt.Errorf("Role dengan nama %s sudah terdaftar", roleName)
	}

	uc.rolerepo.EditRole(roleId, roleName, active, permissions)
	return nil
}
