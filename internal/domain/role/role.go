package roledomain

import (
	"sync"
)

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Active      bool     `json:"active"`
	Permissions []string `json:"permissions"`
}

type RoleResponseModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleCache struct {
	sync.RWMutex
	DataMap map[string]*Role
}

type Permission struct {
	ID      string `json:"id"`
	Menu    string `json:"menu"`
	SubMenu string `json:"subMenu"`
	Name    string `json:"name"`
}
