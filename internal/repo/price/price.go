package pricerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	pricedomain "dromatech/pos-backend/internal/domain/price"
	stringutil "dromatech/pos-backend/internal/util/string"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type PriceRepo interface {
	Find(params map[string]interface{}) ([]*pricedomain.PriceTemplate, error)
	FindDetail(params map[string]interface{}) ([]*pricedomain.PriceTemplateDetail, error)
	Create(name string) (string, error)
	AddPrice(priceTemplateId string, productId string, price float64) error
	EditPrice(priceTemplateId string, productId string, price float64) error
	DeleteTemplate(templateId string) error
	UpdateTemplate(priceTemplateId, customerId string, tx *gorm.DB)
	UpdateChecked(request pricedomain.Download) error
	FindBuyTemplate(params map[string]interface{}) ([]*pricedomain.PriceTemplate, error)
	FindBuyDetail(params map[string]interface{}) ([]*pricedomain.PriceTemplateDetail, error)
	CreateBuyTemplate(name string) (string, error)
	AddBuyPrice(priceTemplateId string, productId string, price float64) error
	EditBuyPrice(priceTemplateId string, productId string, price float64) error
	DeleteBuyTemplate(templateId string) error
	UpdateBuyTemplate(priceTemplateId, webUserId string, trxId string, createdTime time.Time, tx *gorm.DB)
	UpdateBuyChecked(request pricedomain.Download) error
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT pt.id, pt.name, pt.applied_to FROM price_template pt %s", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*pricedomain.PriceTemplate
	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString
		var AppliedTo sql.NullString

		err = rows.Scan(&ID, &Name, &AppliedTo)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}

		entity := &pricedomain.PriceTemplate{}
		if ID.Valid && ID.String != "" {
			entity.ID = ID.String
		} else {
			return nil, nil
		}

		entity = &pricedomain.PriceTemplate{}
		entity.ID = ID.String
		entity.Name = Name.String
		entity.AppliedTo = AppliedTo.String

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) FindBuyTemplate(params map[string]interface{}) ([]*pricedomain.PriceTemplate, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT pt.id, pt.name FROM buy_price_template pt %s", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*pricedomain.PriceTemplate
	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString

		err = rows.Scan(&ID, &Name)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}

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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT ptd.id, ptd.product_id, ptd.price, ptd.checked FROM public.price_template_detail ptd %s", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*pricedomain.PriceTemplateDetail
	for rows.Next() {
		var ID sql.NullString
		var ProductID sql.NullString
		var Price sql.NullFloat64
		var Checked sql.NullBool

		rows.Scan(&ID, &ProductID, &Price, &Checked)

		entity := &pricedomain.PriceTemplateDetail{}
		if ProductID.Valid && ProductID.String != "" {
			entity.ProductID = ProductID.String
		} else {
			return nil, nil
		}

		entity = &pricedomain.PriceTemplateDetail{}
		entity.ID = ID.String
		entity.ProductID = ProductID.String
		entity.Price = Price.Float64
		entity.Checked = Checked.Bool

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) FindBuyDetail(params map[string]interface{}) ([]*pricedomain.PriceTemplateDetail, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT ptd.id, ptd.product_id, ptd.price, ptd.checked FROM public.buy_price_template_detail ptd %s", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*pricedomain.PriceTemplateDetail
	for rows.Next() {
		var ID sql.NullString
		var ProductID sql.NullString
		var Price sql.NullFloat64
		var Checked sql.NullBool

		rows.Scan(&ID, &ProductID, &Price, &Checked)

		entity := &pricedomain.PriceTemplateDetail{}
		if ProductID.Valid && ProductID.String != "" {
			entity.ProductID = ProductID.String
		} else {
			return nil, nil
		}

		entity = &pricedomain.PriceTemplateDetail{}
		entity.ID = ID.String
		entity.ProductID = ProductID.String
		entity.Price = Price.Float64
		entity.Checked = Checked.Bool

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) Create(name string) (string, error) {
	ID := stringutil.GenerateUUID()

	return ID, global.DBCON.Exec("INSERT INTO public.price_template(id, name) VALUES (?, ?);", ID, name).Error
}

func (r *Repo) CreateBuyTemplate(name string) (string, error) {
	ID := stringutil.GenerateUUID()

	return ID, global.DBCON.Exec("INSERT INTO public.buy_price_template(id, name) VALUES (?, ?);", ID, name).Error
}

func (r *Repo) AddPrice(priceTemplateId string, productId string, price float64) error {
	ID := stringutil.GenerateUUID()

	return global.DBCON.Exec("INSERT INTO public.price_template_detail(id, price_template_id, product_id, price) VALUES (?, ?, ?, ?);", ID, priceTemplateId, productId, price).Error
}

func (r *Repo) AddBuyPrice(priceTemplateId string, productId string, price float64) error {
	ID := stringutil.GenerateUUID()

	return global.DBCON.Exec("INSERT INTO public.buy_price_template_detail(id, price_template_id, product_id, price) VALUES (?, ?, ?, ?);", ID, priceTemplateId, productId, price).Error
}

func (r *Repo) EditPrice(priceTemplateId string, productId string, price float64) error {
	return global.DBCON.Exec("UPDATE public.price_template_detail SET price=? WHERE price_template_id=? AND product_id=?;", price, priceTemplateId, productId).Error
}

func (r *Repo) EditBuyPrice(priceTemplateId string, productId string, price float64) error {
	return global.DBCON.Exec("UPDATE public.buy_price_template_detail SET price=? WHERE price_template_id=? AND product_id=?;", price, priceTemplateId, productId).Error
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

func (r *Repo) DeleteBuyTemplate(templateId string) error {
	tx := global.DBCON.Begin()

	tx.Exec("DELETE FROM public.buy_price_template_detail WHERE price_template_id = ?;", templateId)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("DELETE FROM public.buy_price_template WHERE id = ?;", templateId)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	return tx.Commit().Error
}

func (r *Repo) UpdateTemplate(priceTemplateId, customerId string, tx *gorm.DB) {
	tx.Exec("UPDATE public.price_template SET applied_to=? WHERE id=?;", customerId, priceTemplateId)
}

func (r *Repo) UpdateBuyTemplate(priceTemplateId, webUserId string, txId string, createdTime time.Time, tx *gorm.DB) {
	id := stringutil.GenerateUUID()
	tx.Exec("INSERT INTO public.buy_price_template_transaction(id, buy_price_template_id, transaction_id, created_time, web_user_id) VALUES (?,?,?,?,?);", id, priceTemplateId, txId, createdTime, webUserId)
	if tx.Error != nil {
		return
	}
}

func (r *Repo) UpdateChecked(request pricedomain.Download) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.price_template_detail SET checked=FALSE WHERE price_template_id = ?;", request.TemplateID)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Exec("UPDATE public.price_template_detail SET checked=TRUE WHERE id IN ?;", request.TemplateDetailIDs)
	return tx.Commit().Error
}

func (r *Repo) UpdateBuyChecked(request pricedomain.Download) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.price_template_detail SET checked=FALSE WHERE price_template_id = ?;", request.TemplateID)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Exec("UPDATE public.price_template_detail SET checked=TRUE WHERE id IN ?;", request.TemplateDetailIDs)
	return tx.Commit().Error
}

func (r *Repo) AddBuyPriceTx(entity customerdomain.AddPriceRequest, tx *gorm.DB) {
	tx.Exec("UPDATE public.sell_price SET latest=FALSE "+
		"WHERE customer_id = ? AND product_id = ? AND latest=TRUE;", entity.CustomerId, entity.ProductID)

	if tx.Error != nil {
		return
	}

	tx.Exec("INSERT INTO public.sell_price(id, date, customer_id, product_id, price, web_user_id, latest, transaction_id) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", entity.ID, entity.Date, entity.CustomerId, entity.ProductID, entity.Price, entity.WebUserId, entity.Latest, entity.TransactionId)
}
