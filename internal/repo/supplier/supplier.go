package supplierrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	"fmt"
	"github.com/sirupsen/logrus"
)

type SupplierRepo interface {
	Find(params map[string]interface{}) ([]*supplierdomain.Supplier, error)
	Create(product *supplierdomain.Supplier) error
	Edit(product *supplierdomain.Supplier) error
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params map[string]interface{}) ([]*supplierdomain.Supplier, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT id, code, name, description, active FROM supplier %s ORDER BY code", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*supplierdomain.Supplier

	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var Name sql.NullString
		var Description sql.NullString
		var Active sql.NullBool

		rows.Scan(&ID, &Code, &Name, &Description, &Active)

		entity := &supplierdomain.Supplier{}
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

func (r *Repo) Create(entity *supplierdomain.Supplier) error {
	return global.DBCON.Exec("INSERT INTO public.supplier(id, code, name, description, active) "+
		"VALUES (?, ?, ?, ?, ?)",
		entity.ID, entity.Code, entity.Name, entity.Description, entity.Active).Error
}

func (r *Repo) Edit(entity *supplierdomain.Supplier) error {
	return global.DBCON.Exec("UPDATE public.supplier "+
		"SET code=?, name=?, description=?, active=? "+
		"WHERE id=?;", entity.Code, entity.Name, entity.Description, entity.Active, entity.ID).Error
}

func (r *Repo) UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error {
	tx := global.DBCON.Begin()

	for _, detail := range request.Prices {
		tx.Exec("INSERT INTO public.buy_price(date, supplier_id, unit_id, product_id, price) "+
			"VALUES (?, ?, ?, ?, ?)", request.Date, request.SupplierId, request.UnitId, detail.ProductID, detail.Price)

		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	return tx.Error
}
