package permissiondomain

import (
	menudomain "dromatech/pos-backend/internal/domain/menu"
	roledomain "dromatech/pos-backend/internal/domain/role"
)

type Permission struct {
	ID              string          `json:"id"`
	Menu            menudomain.Menu `json:"menu"`
	Name            string          `json:"name"`
	PermissionOrder string          `json:"permissionOrder"`
	Outcome         string          `json:"outcome"`
	Paths           roledomain.Role `json:"paths"`
	Icon            bool            `json:"icon"`
}
