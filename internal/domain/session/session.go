package sessiondomain

import (
	"sync"
	"time"
)

type Menu struct {
	Name    string    `json:"name"`
	Icon    string    `json:"icon"`
	SubMenu []SubMenu `json:"subMenu"`
}

type SubMenu struct {
	Name    string `json:"name"`
	Outcome string `json:"outcome"`
	Icon    string `json:"icon"`
}

type Session struct {
	Token       string    `json:"token"`
	ExpiredTime time.Time `json:"-"`
	UserID      string    `json:"-"`
	UserName    string    `json:"username"`
	Name        string    `json:"name"`
	RoleName    string    `json:"roleName"`
	Menu        []*Menu   `json:"menu"`
}

type SessionCache struct {
	sync.RWMutex
	DataMap map[string]*Session
}
