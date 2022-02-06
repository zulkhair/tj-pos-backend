package rolerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	roledomain "dromatech/pos-backend/internal/domain/role"
	"github.com/sirupsen/logrus"
)

type UserRepo interface {
	ReInitCache() error
}

type Repo struct {
	RoleCache roledomain.RoleCache
}

func New() (*Repo, error) {
	repo := &Repo{}
	err := repo.ReInitCache()
	return repo, err
}

func (r *Repo) ReInitCache() error {
	r.RoleCache.Lock()
	r.RoleCache.DataMap = make(map[string]roledomain.Role)

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

		role := roledomain.Role{}
		if ID.Valid {
			role.ID = ID.String
		}

		if Name.Valid {
			role.Name = Name.String
		}

		if Active.Valid {
			role.Active = Active.Bool
		}

		r.RoleCache.DataMap[role.ID] = role
	}

	r.RoleCache.Unlock()

	return nil
}
