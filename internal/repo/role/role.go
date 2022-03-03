package rolerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	roledomain "dromatech/pos-backend/internal/domain/role"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type RoleRepo interface {
	Find(id string) *roledomain.Role
	FindMenu(roleId string) ([]*sessiondomain.Menu, error)
	FindAll() ([]*roledomain.Role, error)
	FindActive() ([]*roledomain.RoleResponseModel, error)
	FindActivePermissionPaths(roleId string) ([]string, error)
	FindPermissions() ([]*roledomain.Permission, error)
	FindPermissionsByRoleId(roleId string) ([]*roledomain.Permission, error)
	RegisterRole(name string, permissions []string)
	FindByName(name string) *roledomain.Role
	EditRole(roleId string, name string, active bool, permissions []string)
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(id string) *roledomain.Role {
	row := global.DBCON.Raw("SELECT id, name, active FROM role WHERE id = ? AND active = true", id).Row()
	var ID sql.NullString
	var Name sql.NullString
	var Active sql.NullBool

	row.Scan(&ID, &Name, &Active)

	role := &roledomain.Role{}
	if ID.Valid {
		role.ID = ID.String
	}

	if Name.Valid {
		role.Name = Name.String
	}

	if Active.Valid {
		role.Active = Active.Bool
	}

	return role
}

func (r *Repo) FindByName(name string) *roledomain.Role {
	row := global.DBCON.Raw("SELECT id, name, active FROM role WHERE LOWER(name) = ? AND active = true", strings.ToLower(name)).Row()
	var ID sql.NullString
	var Name sql.NullString
	var Active sql.NullBool

	row.Scan(&ID, &Name, &Active)

	role := &roledomain.Role{}
	if ID.Valid && ID.String != "" {
		role.ID = ID.String
	} else {
		return nil
	}

	if Name.Valid {
		role.Name = Name.String
	}

	if Active.Valid {
		role.Active = Active.Bool
	}

	return role
}

func (r *Repo) FindAll() ([]*roledomain.Role, error) {
	rows, err := global.DBCON.Raw("SELECT id, name, active FROM role").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var roles []*roledomain.Role
	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString
		var Active sql.NullBool

		rows.Scan(&ID, &Name, &Active)

		role := &roledomain.Role{}
		if ID.Valid {
			role.ID = ID.String
		}

		if Name.Valid {
			role.Name = Name.String
		}

		if Active.Valid {
			role.Active = Active.Bool
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *Repo) FindActive() ([]*roledomain.RoleResponseModel, error) {
	rows, err := global.DBCON.Raw("SELECT id, name FROM role WHERE active = true").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var roles []*roledomain.RoleResponseModel
	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString

		rows.Scan(&ID, &Name)

		role := &roledomain.RoleResponseModel{}
		if ID.Valid {
			role.ID = ID.String
		}

		if Name.Valid {
			role.Name = Name.String
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *Repo) FindMenu(roleId string) ([]*sessiondomain.Menu, error) {
	rows, err := global.DBCON.Raw("SELECT m.name, m.icon, m.path, s.name, s.outcome, s.icon, p.id FROM menu m "+
		"JOIN sub_menu s ON (m.id = s.menu_id) "+
		"JOIN permission p ON (s.id = p.sub_menu_id) "+
		"JOIN role_permission rp ON (p.id = rp.permission_id) "+
		"JOIN role r ON (rp.role_id = r.id) "+
		"WHERE r.id = ? AND r.active = true "+
		"ORDER BY m.seq_order ASC, s.seq_order, p.seq_order ASC", roleId).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	menu := &sessiondomain.Menu{}
	subMenu := &sessiondomain.SubMenu{}
	var menus []*sessiondomain.Menu
	for rows.Next() {
		var menuName sql.NullString
		var menuIcon sql.NullString
		var menuPath sql.NullString
		var subName sql.NullString
		var subOutcome sql.NullString
		var subIcon sql.NullString
		var perName sql.NullString

		rows.Scan(&menuName, &menuIcon, &menuPath, &subName, &subOutcome, &subIcon, &perName)

		if menu.Name != menuName.String {
			menu = &sessiondomain.Menu{
				Name:    menuName.String,
				Icon:    menuIcon.String,
				Path:    menuPath.String,
				SubMenu: nil,
			}
			menus = append(menus, menu)
		}

		if subMenu.Name != subName.String {
			subMenu = &sessiondomain.SubMenu{
				Name:    subName.String,
				Outcome: subOutcome.String,
				Icon:    subIcon.String,
			}
			menu.SubMenu = append(menu.SubMenu, subMenu)
		}

		subMenu.Permissions = append(subMenu.Permissions, perName.String)

	}

	return menus, nil
}

func (r *Repo) FindActivePermissionPaths(roleId string) ([]string, error) {
	rows, err := global.DBCON.Raw("SELECT p.apis FROM permission p "+
		"JOIN role_permission rp ON (p.id = rp.permission_id) "+
		"JOIN role r ON (rp.role_id = r.id) "+
		"WHERE r.id = ? AND r.active = true", roleId).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path sql.NullString

		rows.Scan(&path)

		if path.Valid {
			paths = append(paths, path.String)
		}
	}

	return paths, nil
}

func (r *Repo) FindPermissions() ([]*roledomain.Permission, error) {
	rows, err := global.DBCON.Raw("SELECT p.id, m.name, s.name, p.name FROM permission p " +
		"JOIN sub_menu s ON (s.id = p.sub_menu_id) " +
		"JOIN menu m ON (s.menu_id = m.id) " +
		"ORDER BY m.seq_order, s.seq_order, p.seq_order").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var permissions []*roledomain.Permission
	for rows.Next() {
		permission := &roledomain.Permission{}
		rows.Scan(&permission.ID, &permission.Menu, &permission.SubMenu, &permission.Name)

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *Repo) FindPermissionsByRoleId(roleId string) ([]*roledomain.Permission, error) {
	rows, err := global.DBCON.Raw("SELECT p.id, m.name, s.name, p.name FROM permission p "+
		"JOIN sub_menu s ON (p.sub_menu_id = s.id) "+
		"JOIN menu m ON (s.menu_id = m.id) "+
		"JOIN role_permission rp ON (rp.permission_id = p.id) "+
		"WHERE rp.role_id = ?"+
		"ORDER BY m.seq_order, s.seq_order, p.seq_order", roleId).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var permissions []*roledomain.Permission
	for rows.Next() {
		permission := &roledomain.Permission{}
		rows.Scan(&permission.ID, &permission.Menu, &permission.SubMenu, &permission.Name)

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *Repo) RegisterRole(name string, permissions []string) {
	tx := global.DBCON.Begin()
	roleId := strings.ReplaceAll(uuid.NewString(), "-", "")
	tx.Exec("INSERT INTO public.role(id, active, name) VALUES (?, ?, ?);",
		roleId, true, name)

	if tx.Error != nil {
		tx.Rollback()
		return
	}

	for _, id := range permissions {
		tx.Exec("INSERT INTO public.role_permission(role_id, permission_id) VALUES (?, ?);",
			roleId, id)

		if tx.Error != nil {
			tx.Rollback()
			return
		}
	}

	if tx.Error != nil {
		tx.Rollback()
	}
	tx.Commit()
}

func (r *Repo) EditRole(roleId string, name string, active bool, permissions []string) {
	tx := global.DBCON.Begin()
	tx.Exec("DELETE FROM public.role_permission WHERE role_id = ?;", roleId)

	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return
	}

	for _, id := range permissions {
		tx.Exec("INSERT INTO public.role_permission(role_id, permission_id) VALUES (?, ?);",
			roleId, id)

		if tx.Error != nil {
			tx.Rollback()
			logrus.Error(tx.Error.Error())
			return
		}
	}

	tx.Exec("UPDATE public.role SET active=?, name=? WHERE id=?;",
		active, name, roleId)

	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return
	}

	tx.Commit()
}
