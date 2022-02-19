package menurepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	menudomain "dromatech/pos-backend/internal/domain/menu"
	"github.com/sirupsen/logrus"
)

type MenuRepo interface {
	ReInitCache() error
	Find(id string) *menudomain.Menu
}

type Repo struct {
	menuCache menudomain.MenuCache
}

func New() (*Repo, error) {
	repo := &Repo{}
	err := repo.ReInitCache()
	return repo, err
}

func (r *Repo) ReInitCache() error {
	r.menuCache.Lock()
	r.menuCache.DataMap = make(map[string]*menudomain.Menu)

	rows, err := global.DBCON.Raw("SELECT id, name, menu_order, menu_path, icon FROM menu").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ID sql.NullString
		var name sql.NullString
		var menuOrder sql.NullInt64
		var menuPath sql.NullString
		var icon sql.NullString

		rows.Scan(&ID, &name, &menuOrder, &menuPath, icon)

		menu := &menudomain.Menu{}
		if ID.Valid {
			menu.ID = ID.String
		}

		if name.Valid {
			menu.Name = name.String
		}

		if menuOrder.Valid {
			menu.MenuOrder = menuOrder.Int64
		}

		if menuPath.Valid {
			menu.MenuPath = menuPath.String
		}

		if icon.Valid {
			menu.Icon = icon.String
		}

		r.menuCache.DataMap[menu.ID] = menu
	}

	r.menuCache.Unlock()

	return nil
}

func (r *Repo) Find(id string) *menudomain.Menu {
	if menu, ok := r.menuCache.DataMap[id]; ok {
		return menu
	}
	return nil
}
