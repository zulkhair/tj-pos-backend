package customerrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/sirupsen/logrus"
)

type CustomerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
	GetSellPrice(params []queryutil.Param) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
	DeleteSellPrice(supplierId, unitId, date string) error
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

func (r *Repo) GetSellPrice(params []queryutil.Param) ([]*customerdomain.SellPriceResponse, error) {
	where := ""
	var values []interface{}
	for _, param := range params {
		if where != "" {
			logic := "AND "
			if param.Logic != "" {
				logic = param.Logic + " "
			}
			where += logic
		}
		where += param.Field + " " + param.Operator + " ? "
		values = append(values, param.Value)
	}

	if where != "" {
		where = "WHERE " + where
	}

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT p.id, p.code, p.name, p.description, bp.price FROM product p "+
		"JOIN sell_price bp ON (p.id = bp.product_id) "+
		"%s ORDER BY p.code", where), values...).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*customerdomain.SellPriceResponse

	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var Name sql.NullString
		var Description sql.NullString
		var Price sql.NullFloat64

		rows.Scan(&ID, &Code, &Name, &Description, &Price)

		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		entity := &customerdomain.SellPriceResponse{
			ProductID:   ID.String,
			ProductCode: Code.String,
			ProductName: Name.String,
			ProductDesc: Description.String,
			Price:       Price.Float64,
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) UpdateSellPrice(request customerdomain.SellPriceRequest) error {
	tx := global.DBCON.Begin()

	for _, detail := range request.Prices {
		tx.Exec("INSERT INTO public.sell_price(date, customer_id, unit_id, product_id, price) "+
			"VALUES (?, ?, ?, ?, ?)", request.Date, request.CustomerId, request.UnitId, detail.ProductID, detail.Price)

		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}
	tx.Commit()
	return tx.Error
}

func (r *Repo) DeleteSellPrice(customerId, unitId, date string) error {
	tx := global.DBCON.Exec("DELETE from public.sell_price WHERE customer_id = ? AND unit_id = ? AND date = ? ",
		customerId, unitId, date)

	return tx.Error
}
