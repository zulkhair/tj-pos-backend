package permissionrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	permissiondomain "dromatech/pos-backend/internal/domain/permission"
	"github.com/sirupsen/logrus"
)

type PermissionRepo interface {
	ReInitCache() error
	Find(id string) *permissiondomain.Permission
	FindByRoleId(roleId string) *permissiondomain.Permission
}

type Repo struct {
	permissionCache     permissiondomain.PermissionCache
	rolePermissionCache permissiondomain.RolePermissionCache
}

func New() (*Repo, error) {
	repo := &Repo{}
	err := repo.ReInitCache()
	return repo, err
}

func (r *Repo) ReInitCache() error {
	err := r.gatherPermission()
	if err != nil {
		return err
	}

	return r.gatherRolePermission()
}

func (r *Repo) gatherPermission() error {
	rows, err := global.DBCON.Raw("SELECT id, menu_id, name, permission_order, outcome, paths, icon FROM permission").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	defer rows.Close()

	r.permissionCache.Lock()
	r.permissionCache.DataMap = make(map[string]*permissiondomain.Permission)

	for rows.Next() {
		var ID sql.NullString
		var menuID sql.NullString
		var name sql.NullString
		var permissionOrder sql.NullInt64
		var outcome sql.NullString
		var paths sql.NullString
		var icon sql.NullString

		rows.Scan(&ID, &menuID, &name, &permissionOrder, &outcome, &paths, icon)

		permission := &permissiondomain.Permission{}
		if ID.Valid {
			permission.ID = ID.String
		}

		if menuID.Valid {
			permission.MenuID = menuID.String
		}

		if name.Valid {
			permission.Name = name.String
		}

		if permissionOrder.Valid {
			permission.PermissionOrder = permissionOrder.Int64
		}

		if outcome.Valid {
			permission.Outcome = outcome.String
		}

		if paths.Valid {
			permission.Paths = paths.String
		}

		if icon.Valid {
			permission.Icon = icon.String
		}

		r.permissionCache.DataMap[permission.ID] = permission
	}

	r.permissionCache.Unlock()

	return nil
}

func (r *Repo) gatherRolePermission() error {
	rows, err := global.DBCON.Raw("SELECT role_id, permission_id FROM role_permission").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	defer rows.Close()

	r.rolePermissionCache.Lock()
	r.rolePermissionCache.DataMap = make(map[string][]*permissiondomain.Permission)

	for rows.Next() {
		var RoleID sql.NullString
		var PermissionID sql.NullString

		rows.Scan(&RoleID, &PermissionID)

		if RoleID.Valid && PermissionID.Valid {
			r.rolePermissionCache.DataMap[RoleID.String] = append(r.rolePermissionCache.DataMap[RoleID.String], r.Find(PermissionID.String))
		}
	}

	r.rolePermissionCache.Unlock()

	return nil
}

func (r *Repo) Find(id string) *permissiondomain.Permission {
	if permission, ok := r.permissionCache.DataMap[id]; ok {
		return permission
	}
	return nil
}

func (r *Repo) FindByRoleId(roleId string) []*permissiondomain.Permission {
	if permissions, ok := r.rolePermissionCache.DataMap[roleId]; ok {
		return permissions
	}
	return nil
}
