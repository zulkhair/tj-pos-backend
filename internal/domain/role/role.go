package roledomain

import (
	permissiondomain "dromatech/pos-backend/internal/domain/permission"
	"sync"
)

type Role struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type RoleCache struct {
	sync.RWMutex
	DataMap map[string]*Role
}

type RoleMenuPermissionCache struct {
	sync.RWMutex
	DataMap map[string][]*permissiondomain.Permission
}
