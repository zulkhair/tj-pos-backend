package unitrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	unitdomain "dromatech/pos-backend/internal/domain/unit"
	"fmt"
	"github.com/sirupsen/logrus"
)

type UnitRepo interface {
	Find(params map[string]interface{}) ([]*unitdomain.Unit, error)
	Create(product *unitdomain.Unit) error
	Edit(product *unitdomain.Unit) error
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params map[string]interface{}) ([]*unitdomain.Unit, error) {
	where := ""
	var values []interface{}
	for key, value := range params {
		if where != "" {
			where += "AND "
		}
		where += key + " = ? "
		values = append(values, value)
	}

	if where != "" {
		where = "WHERE " + where
	}

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT id, code, description FROM unit %s ORDER BY code", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*unitdomain.Unit

	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString

		rows.Scan(&ID, &Code)

		entity := &unitdomain.Unit{}
		if ID.Valid && ID.String != "" {
			entity.ID = ID.String
		} else {
			return nil, nil
		}

		if Code.Valid {
			entity.Code = Code.String
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) Create(entity *unitdomain.Unit) error {
	return global.DBCON.Exec("INSERT INTO public.supplier(id, code, description) "+
		"VALUES (?, ?, ?)",
		entity.ID, entity.Code, entity.Description).Error
}

func (r *Repo) Edit(entity *unitdomain.Unit) error {
	return global.DBCON.Exec("UPDATE public.supplier "+
		"SET code=?, description=? "+
		"WHERE id=?;", entity.Code, entity.Description, entity.ID).Error
}
