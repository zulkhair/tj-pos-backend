package productrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	productdomain "dromatech/pos-backend/internal/domain/product"
	"fmt"
	"github.com/sirupsen/logrus"
)

type ProductRepo interface {
	Find(params map[string]interface{}) ([]*productdomain.Product, error)
	Create(product *productdomain.Product) error
	Edit(product *productdomain.Product) error
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params map[string]interface{}) ([]*productdomain.Product, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT p.id, p.code, p.name, p.description, p.active, u.id, u.code FROM product p JOIN unit u ON (u.id = p.unit_id) %s ORDER BY p.name", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var products []*productdomain.Product

	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var Name sql.NullString
		var Description sql.NullString
		var Active sql.NullBool
		var UnitID sql.NullString
		var UnitCode sql.NullString

		rows.Scan(&ID, &Code, &Name, &Description, &Active, &UnitID,  &UnitCode)

		product := &productdomain.Product{}
		if ID.Valid && ID.String != "" {
			product.ID = ID.String
		} else {
			return nil, nil
		}

		if Code.Valid {
			product.Code = Code.String
		}

		if Name.Valid {
			product.Name = Name.String
		}

		if Description.Valid {
			product.Description = Description.String
		}

		if Active.Valid {
			product.Active = Active.Bool
		}

		product.UnitID = UnitID.String
		product.UnitCode = UnitCode.String

		products = append(products, product)
	}

	return products, nil
}

func (r *Repo) Create(product *productdomain.Product) error {
	return global.DBCON.Exec("INSERT INTO public.product(id, code, name, description, active, unit_id) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		product.ID, product.Code, product.Name, product.Description, product.Active, product.UnitID).Error
}

func (r *Repo) Edit(product *productdomain.Product) error {
	return global.DBCON.Exec("UPDATE public.product "+
		"SET code=?, name=?, description=?, active=?, unit_id = ? "+
		"WHERE id=?;", product.Code, product.Name, product.Description, product.Active, product.UnitID, product.ID).Error
}
