package roleusecase

import (
	roledomain "dromatech/pos-backend/internal/domain/role"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
)

type RoleUsecase interface {
	GetActiveRole() ([]*roledomain.RoleResponseModel, error)
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
