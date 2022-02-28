package roledomain

import (
	"sync"
)

type Role struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type RoleResponseModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleCache struct {
	sync.RWMutex
	DataMap map[string]*Role
}
