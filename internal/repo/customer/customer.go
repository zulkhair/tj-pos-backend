package customerrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type CustomerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
	GetSellPrice(params []queryutil.Param) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
	DeleteSellPrice(customerId, date string) error
	AddSellPrice(entity customerdomain.AddPriceRequest) error
	AddSellPriceTx(entity customerdomain.AddPriceRequest, tx *gorm.DB)
	FindSellPrice(params []queryutil.Param) ([]*customerdomain.PriceResponse, error)
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT id, code, name, description, active, initial_credit FROM customer %s ORDER BY code", where), values...).Rows()
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
		var InitialCredit sql.NullFloat64

		rows.Scan(&ID, &Code, &Name, &Description, &Active, &InitialCredit)

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

		entity.InitialCredit = InitialCredit.Float64

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) Create(entity *customerdomain.Customer) error {
	return global.DBCON.Exec("INSERT INTO public.customer(id, code, name, description, active, initial_credit) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		entity.ID, entity.Code, entity.Name, entity.Description, entity.Active, entity.InitialCredit).Error
}

func (r *Repo) Edit(entity *customerdomain.Customer) error {
	return global.DBCON.Exec("UPDATE public.customer "+
		"SET code=?, name=?, description=?, active=?, initial_credit=? "+
		"WHERE id=?;", entity.Code, entity.Name, entity.Description, entity.Active, entity.InitialCredit, entity.ID).Error
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
		tx.Exec("INSERT INTO public.sell_price(date, customer_id, product_id, price) "+
			"VALUES (?, ?, ?, ?, ?)", request.Date, request.CustomerId, detail.ProductID, detail.Price)

		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}
	tx.Commit()
	return tx.Error
}

func (r *Repo) DeleteSellPrice(customerId, date string) error {
	tx := global.DBCON.Exec("DELETE from public.sell_price WHERE customer_id = ? AND date = ? ",
		customerId, date)

	return tx.Error
}

func (r *Repo) AddSellPrice(entity customerdomain.AddPriceRequest) error {

	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.sell_price SET latest=FALSE "+
		"WHERE customer_id = ? AND product_id = ? AND latest=TRUE;", entity.CustomerId, entity.ProductID)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("INSERT INTO public.sell_price(id, date, customer_id, product_id, price, web_user_id, latest, transaction_id) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", entity.ID, entity.Date, entity.CustomerId, entity.ProductID, entity.Price, entity.WebUserId, entity.Latest, entity.TransactionId)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil

}

func (r *Repo) AddSellPriceTx(entity customerdomain.AddPriceRequest, tx *gorm.DB) {
	tx.Exec("UPDATE public.sell_price SET latest=FALSE "+
		"WHERE customer_id = ? AND product_id = ? AND latest=TRUE;", entity.CustomerId, entity.ProductID)

	if tx.Error != nil {
		return
	}

	tx.Exec("INSERT INTO public.sell_price(id, date, customer_id, product_id, price, web_user_id, latest, transaction_id) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", entity.ID, entity.Date, entity.CustomerId, entity.ProductID, entity.Price, entity.WebUserId, entity.Latest, entity.TransactionId)
}

func (r *Repo) FindSellPrice(params []queryutil.Param) ([]*customerdomain.PriceResponse, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT s.id, s.date, s.customer_id, p.unit_id, s.product_id, s.price, w.username, w.name, t.code "+
		"FROM public.sell_price s "+
		"JOIN public.web_user w ON (w.id = s.web_user_id) "+
		"JOIN public.product p ON (p.id = s.product_id) "+
		"JOIN public.unit u ON (u.id = p.unit_id) "+
		"LEFT JOIN public.transaction t ON (t.id = s.transaction_id) "+
		"%s ORDER BY date DESC ", where), values...).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*customerdomain.PriceResponse

	for rows.Next() {
		var ID sql.NullString
		var Date time.Time
		var CustomerId sql.NullString
		var UnitId sql.NullString
		var ProductId sql.NullString
		var Price sql.NullFloat64
		var WebUsername sql.NullString
		var WebUserName sql.NullString
		var TransactionCode sql.NullString

		rows.Scan(&ID, &Date, &CustomerId, &UnitId, &ProductId, &Price, &WebUsername, &WebUserName, &TransactionCode)

		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		entity := &customerdomain.PriceResponse{
			ID:              ID.String,
			Date:            Date.Format(dateutil.TimeFormatResponse()),
			CustomerId:      CustomerId.String,
			UnitId:          UnitId.String,
			ProductID:       ProductId.String,
			Price:           Price.Float64,
			WebUsername:     WebUsername.String,
			WebUserName:     WebUserName.String,
			TransactionCode: TransactionCode.String,
		}

		entities = append(entities, entity)
	}

	return entities, nil
}
