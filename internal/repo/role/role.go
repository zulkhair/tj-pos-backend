package rolerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	roledomain "dromatech/pos-backend/internal/domain/role"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
	"github.com/sirupsen/logrus"
)

type RoleRepo interface {
	Find(id string) *roledomain.Role
	FindMenu(roleId string) ([]*sessiondomain.Menu, error)
	FindAll() ([]*roledomain.Role, error)
	FindActive() ([]*roledomain.RoleActive, error)
	FindActivePermissionPaths(roleId string) ([]string, error)
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

func (r *Repo) FindActive() ([]*roledomain.RoleActive, error) {
	rows, err := global.DBCON.Raw("SELECT id, name FROM role WHERE active = true").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var roles []*roledomain.RoleActive
	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString

		rows.Scan(&ID, &Name)

		role := &roledomain.RoleActive{}
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
	rows, err := global.DBCON.Raw("SELECT m.name, m.icon, m.menu_path, p.name, p.outcome, p.icon FROM menu m "+
		"JOIN menu_permission mp ON (m.id = mp.menu_id) "+
		"JOIN permission p ON (mp.permission_id = p.id) "+
		"JOIN role_permission rp ON (p.id = rp.permission_id) "+
		"JOIN role r ON (rp.role_id = r.id) "+
		"WHERE r.id = ? AND r.active = true "+
		"ORDER BY m.menu_order ASC, p.permission_order ASC", roleId).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	menu := &sessiondomain.Menu{}
	var menus []*sessiondomain.Menu
	for rows.Next() {
		var menuName sql.NullString
		var menuIcon sql.NullString
		var menuPath sql.NullString
		var perName sql.NullString
		var perOutcome sql.NullString
		var perIcon sql.NullString

		rows.Scan(&menuName, &menuIcon, &menuPath, &perName, &perOutcome, &perIcon)

		if menu.Name != menuName.String {
			menu = &sessiondomain.Menu{
				Name:    menuName.String,
				Icon:    menuIcon.String,
				Path:    menuPath.String,
				SubMenu: nil,
			}
			menus = append(menus, menu)
		}

		subMenu := sessiondomain.SubMenu{
			Name:    perName.String,
			Outcome: perOutcome.String,
			Icon:    perIcon.String,
		}

		menu.SubMenu = append(menu.SubMenu, subMenu)

	}

	return menus, nil
}

func (r *Repo) FindActivePermissionPaths(roleId string) ([]string, error) {
	rows, err := global.DBCON.Raw("SELECT p.paths FROM permission p "+
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
