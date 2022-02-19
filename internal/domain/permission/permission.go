package permissiondomain

import "sync"

type Permission struct {
	ID              string `json:"id"`
	MenuID          string `json:"menuId"`
	Name            string `json:"name"`
	PermissionOrder int64  `json:"permissionOrder"`
	Outcome         string `json:"outcome"`
	Paths           string `json:"paths"`
	Icon            string `json:"icon"`
}

type PermissionCache struct {
	sync.RWMutex
	DataMap map[string]*Permission
}

type RolePermissionCache struct {
	sync.RWMutex
	DataMap map[string][]*Permission
}
