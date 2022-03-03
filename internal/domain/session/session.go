package sessiondomain

import (
	"sync"
	"time"
)

type Menu struct {
	Name    string     `json:"name"`
	Icon    string     `json:"icon"`
	Path    string     `json:"path"`
	SubMenu []*SubMenu `json:"subMenu"`
}

type SubMenu struct {
	Name        string   `json:"name"`
	Outcome     string   `json:"outcome"`
	Icon        string   `json:"icon"`
	Permissions []string `json:"-"`
}

type Session struct {
	Token       string    `json:"token"`
	ExpiredTime time.Time `json:"-"`
	UserID      string    `json:"-"`
	RoleID      string    `json:"-"`
	UserName    string    `json:"username"`
	Name        string    `json:"name"`
	RoleName    string    `json:"roleName"`
	Menu        []*Menu   `json:"menu"`
	Permissions []string  `json:"-"`
}

type SessionCache struct {
	sync.RWMutex
	DataMap map[string]*Session
}
