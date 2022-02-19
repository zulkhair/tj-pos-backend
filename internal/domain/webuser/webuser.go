package webuserdomain

import (
	"sync"
	"time"
)

type WebUser struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Username              string    `json:"username"`
	PasswordHash          string    `json:"PasswordHash"`
	PasswordSalt          string    `json:"passwordSalt"`
	Email                 string    `json:"email"`
	RoleId                string    `json:"roleId"`
	Active                bool      `json:"active"`
	RegistrationTimestamp time.Time `json:"registrationTimestamp"`
	CreatedBy             string    `json:"createdBy"`
}

type WebUserCache struct {
	sync.RWMutex
	DataMap map[string]*WebUser
}
