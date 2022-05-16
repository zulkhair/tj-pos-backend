package supplierrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type SupplierRepo interface {
	Find(params map[string]interface{}) ([]*supplierdomain.Supplier, error)
	Create(product *supplierdomain.Supplier) error
	Edit(product *supplierdomain.Supplier) error
	GetBuyPrice(params []queryutil.Param) ([]*supplierdomain.BuyPriceResponse, error)
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
	DeleteBuyPrice(supplierId, unitId, date string) error
	AddBuyPrice(entity supplierdomain.AddPriceRequest) error
	FindBuyPrice(params []queryutil.Param) ([]*supplierdomain.PriceResponse, error)
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

func (r *Repo) GetBuyPrice(params []queryutil.Param) ([]*supplierdomain.BuyPriceResponse, error) {
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
		"JOIN buy_price bp ON (p.id = bp.product_id) "+
		"%s ORDER BY p.code", where), values...).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*supplierdomain.BuyPriceResponse

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

		entity := &supplierdomain.BuyPriceResponse{
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
	tx.Commit()
	return tx.Error
}

func (r *Repo) DeleteBuyPrice(supplierId, unitId, date string) error {
	tx := global.DBCON.Exec("DELETE from public.buy_price WHERE supplier_id = ? AND unit_id = ? AND date = ? ",
		supplierId, unitId, date)

	return tx.Error
}

func (r *Repo) AddBuyPrice(entity supplierdomain.AddPriceRequest) error {

	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.buy_price SET latest=FALSE "+
		"WHERE unit_id = ? AND product_id = ? AND latest=TRUE;", entity.UnitId, entity.ProductID)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("INSERT INTO public.buy_price(id, date, unit_id, product_id, price, web_user_id, latest) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?)", entity.ID, entity.Date, entity.UnitId, entity.ProductID, entity.Price, entity.WebUserId, entity.Latest)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *Repo) FindBuyPrice(params []queryutil.Param) ([]*supplierdomain.PriceResponse, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT s.id, s.date, s.supplier_id, s.unit_id, s.product_id, s.price, w.username, w.name "+
		"FROM public.buy_price s "+
		"JOIN public.web_user w ON (w.id = s.web_user_id) "+
		"%s ORDER BY date DESC ", where), values...).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*supplierdomain.PriceResponse

	for rows.Next() {
		var ID sql.NullString
		var Date time.Time
		var SupplierId sql.NullString
		var UnitId sql.NullString
		var ProductId sql.NullString
		var Price sql.NullFloat64
		var WebUsername sql.NullString
		var WebUserName sql.NullString

		rows.Scan(&ID, &Date, &SupplierId, &UnitId, &ProductId, &Price, &WebUsername, &WebUserName)

		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		entity := &supplierdomain.PriceResponse{
			ID:          ID.String,
			Date:        Date.Format(dateutil.TimeFormat()),
			SupplierId:  SupplierId.String,
			UnitId:      UnitId.String,
			ProductID:   ProductId.String,
			Price:       Price.Float64,
			WebUsername: WebUsername.String,
			WebUserName: WebUserName.String,
		}

		entities = append(entities, entity)
	}

	return entities, nil
}
