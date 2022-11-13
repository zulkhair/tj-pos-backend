package pricerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	pricedomain "dromatech/pos-backend/internal/domain/price"
	stringutil "dromatech/pos-backend/internal/util/string"
	"fmt"
	"github.com/sirupsen/logrus"
)

type PriceRepo interface {
	Find(params map[string]interface{}) ([]*pricedomain.PriceTemplate, error)
	FindDetail(params map[string]interface{}) ([]*pricedomain.PriceTemplateDetail, error)
	Create(name string) error
	AddPrice(priceTemplateId string, productId string, price float64) error
	EditPrice(priceTemplateId string, productId string, price float64) error
	DeleteTemplate(templateId string) error
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params map[string]interface{}) ([]*pricedomain.PriceTemplate, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT pt.id, pt.name FROM price_template pt %s", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*pricedomain.PriceTemplate
	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString

		rows.Scan(&ID, &Name)

		entity := &pricedomain.PriceTemplate{}
		if ID.Valid && ID.String != "" {
			entity.ID = ID.String
		} else {
			return nil, nil
		}

		entity = &pricedomain.PriceTemplate{}
		entity.ID = ID.String
		entity.Name = Name.String

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) FindDetail(params map[string]interface{}) ([]*pricedomain.PriceTemplateDetail, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT ptd.product_id, ptd.price FROM public.price_template_detail ptd %s", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*pricedomain.PriceTemplateDetail
	for rows.Next() {
		var ProductID sql.NullString
		var Price sql.NullFloat64

		rows.Scan(&ProductID, &Price)

		entity := &pricedomain.PriceTemplateDetail{}
		if ProductID.Valid && ProductID.String != "" {
			entity.ProductID = ProductID.String
		} else {
			return nil, nil
		}

		entity = &pricedomain.PriceTemplateDetail{}
		entity.ProductID = ProductID.String
		entity.Price = Price.Float64

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) Create(name string) error {
	ID := stringutil.GenerateUUID()

	return global.DBCON.Exec("INSERT INTO public.price_template(id, name) VALUES (?, ?);", ID, name).Error
}

func (r *Repo) AddPrice(priceTemplateId string, productId string, price float64) error{
	ID := stringutil.GenerateUUID()

	return global.DBCON.Exec("INSERT INTO public.price_template_detail(id, price_template_id, product_id, price) VALUES (?, ?, ?, ?);", ID, priceTemplateId, productId, price).Error
}

func (r *Repo) EditPrice(priceTemplateId string, productId string, price float64) error {
	return global.DBCON.Exec("UPDATE public.price_template_detail SET price=? WHERE price_template_id=? AND product_id=?;", price, priceTemplateId, productId).Error
}

func (r *Repo) DeleteTemplate(templateId string) error {
	tx := global.DBCON.Begin()

	tx.Exec("DELETE FROM public.price_template_detail WHERE price_template_id = ?;", templateId)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("DELETE FROM public.price_template WHERE id = ?;", templateId)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	return tx.Commit().Error
}