package webuserdomain

import (
	roledomain "dromatech/pos-backend/internal/domain/role"
	"sync"
)

type WebUser struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Username              string          `json:"username"`
	PasswordSalt          string          `json:"passwordSalt"`
	Email                 string          `json:"email"`
	Role                  roledomain.Role `json:"role"`
	Aactive               bool            `json:"active"`
	RegistrationTimestamp string          `json:"registrationTimestamp"`
	CreatedBy             string          `json:"createdBy"`
}

type WebUserCache struct {
	sync.RWMutex
	DataMap map[string]WebUser
}
