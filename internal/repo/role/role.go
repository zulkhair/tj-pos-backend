package rolerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	roledomain "dromatech/pos-backend/internal/domain/role"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
	"github.com/sirupsen/logrus"
)

type RoleRepo interface {
	ReInitCache() error
	Find(id string) *roledomain.Role
	FindMenu(roleId string) ([]sessiondomain.Menu, error)
}

type Repo struct {
	roleCache roledomain.RoleCache
}

func New() (*Repo, error) {
	repo := &Repo{}
	err := repo.ReInitCache()
	return repo, err
}

func (r *Repo) ReInitCache() error {
	r.roleCache.Lock()
	r.roleCache.DataMap = make(map[string]*roledomain.Role)

	rows, err := global.DBCON.Raw("SELECT id, name, active FROM role").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	defer rows.Close()

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

		r.roleCache.DataMap[role.ID] = role
	}

	r.roleCache.Unlock()

	return nil
}

func (r *Repo) Find(id string) *roledomain.Role {
	if role, ok := r.roleCache.DataMap[id]; ok {
		return role
	}
	return nil
}

func (r *Repo) FindMenu(roleId string) ([]*sessiondomain.Menu, error) {
	rows, err := global.DBCON.Raw("SELECT m.name, m.icon, p.name, p.outcome,p.icon FROM menu m "+
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
		var perName sql.NullString
		var perOutcome sql.NullString
		var perIcon sql.NullString

		rows.Scan(&menuName, &menuIcon, &perName, &perOutcome, &perIcon)

		if menu.Name != menuName.String {
			menu = &sessiondomain.Menu{
				Name:    menuName.String,
				Icon:    menuIcon.String,
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
