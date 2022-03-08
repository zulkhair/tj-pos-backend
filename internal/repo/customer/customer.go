package customerrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	"fmt"
	"github.com/sirupsen/logrus"
)

type CustomerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params map[string]interface{}) ([]*customerdomain.Customer, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT id, code, name, description, active FROM customer %s ORDER BY code", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*customerdomain.Customer

	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var Name sql.NullString
		var Description sql.NullString
		var Active sql.NullBool

		rows.Scan(&ID, &Code, &Name, &Description, &Active)

		entity := &customerdomain.Customer{}
		if ID.Valid && ID.String != "" {
			entity.ID = ID.String
		} else {
			return nil, nil
		}

		if Code.Valid {
			entity.Code = Code.String
		}

		if Name.Valid {
			entity.Name = Name.String
		}

		if Description.Valid {
			entity.Description = Description.String
		}

		if Active.Valid {
			entity.Active = Active.Bool
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) Create(entity *customerdomain.Customer) error {
	return global.DBCON.Exec("INSERT INTO public.customer(id, code, name, description, active) "+
		"VALUES (?, ?, ?, ?, ?)",
		entity.ID, entity.Code, entity.Name, entity.Description, entity.Active).Error
}

func (r *Repo) Edit(entity *customerdomain.Customer) error {
	return global.DBCON.Exec("UPDATE public.customer "+
		"SET code=?, name=?, description=?, active=? "+
		"WHERE id=?;", entity.Code, entity.Name, entity.Description, entity.Active, entity.ID).Error
}
