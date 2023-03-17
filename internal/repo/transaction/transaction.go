package transactionrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	stringutil "dromatech/pos-backend/internal/util/string"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type TransactionRepo interface {
	Find(params []queryutil.Param) ([]*transactiondomain.Transaction, error)
	Create(entity *transactiondomain.Transaction, tx *gorm.DB)
	Edit(product *transactiondomain.Transaction) error
	UpdateStatus(transactionID, status string) error
	UpdatePrice(transactionID, productID string, buyPrice, sellPrice float64, quantity, buyQuantity int64) error
	FindSells(params []queryutil.Param) ([]*transactiondomain.TransactionStatus, error)
	UpdateTransaction(entity *transactiondomain.Transaction, tx *gorm.DB)
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params []queryutil.Param) ([]*transactiondomain.Transaction, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT t.id, t.code, t.date, t.stakeholder_id, t.transaction_type, t.status, "+
		"u.unit_id, td.product_id, td.buy_price, td.sell_price, td.quantity, td.buy_quantity "+
		"FROM transaction t "+
		"JOIN transaction_detail td ON (td.transaction_id = t.id) "+
		"JOIN product p ON (p.id = td.product_id) "+
		"JOIN unit u ON (u.id = p.unit_id) "+
		"%s ORDER BY t.date DESC", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*transactiondomain.Transaction
	entityMap := make(map[string]*transactiondomain.Transaction)
	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var Date time.Time
		var StakeholderID sql.NullString
		var TransactionType sql.NullString
		var Status sql.NullString
		var UnitID sql.NullString
		var ProductID sql.NullString
		var BuyPrice sql.NullFloat64
		var SellPrice sql.NullFloat64
		var Quantity sql.NullFloat64
		var BuyQuantity sql.NullFloat64

		rows.Scan(&ID, &Code, &Date, &StakeholderID, &TransactionType, &Status, &UnitID, &ProductID, &BuyPrice, &SellPrice, &Quantity, &BuyQuantity)

		var entity *transactiondomain.Transaction
		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		if value, ok := entityMap[ID.String]; ok {
			entity = value
		} else {
			entity = &transactiondomain.Transaction{}
			entity.ID = ID.String
			entity.Code = Code.String
			entity.Date = Date.Format("2006-0102")
			entity.StakeholderID = StakeholderID.String
			entity.TransactionType = TransactionType.String
			entity.Status = Status.String

			entities = append(entities, entity)
			entityMap[ID.String] = entity
		}

		detail := &transactiondomain.TransactionDetail{}
		detail.TransactionID = entity.ID
		detail.UnitID = UnitID.String
		detail.ProductID = ProductID.String
		detail.BuyPrice = BuyPrice.Float64
		detail.SellPrice = SellPrice.Float64
		detail.Quantity = Quantity.Float64
		detail.BuyQuantity = BuyQuantity.Float64

		entity.TransactionDetail = append(entity.TransactionDetail, detail)

	}

	return entities, nil
}

func (r *Repo) FindSells(params []queryutil.Param) ([]*transactiondomain.TransactionStatus, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT t.id, t.code, t.date, t.stakeholder_id, c.code, c.name, t.transaction_type, t.status, "+
		"t.reference_code, t.web_user_id, w.name, t.created_time, "+
		"u.id, u.code, td.product_id, p.code, p.name, td.buy_price, td.sell_price, td.quantity, td.buy_quantity "+
		"FROM transaction t "+
		"JOIN transaction_detail td ON (td.transaction_id = t.id) "+
		"JOIN customer c ON (c.id = t.stakeholder_id) "+
		"JOIN web_user w ON (w.id = t.web_user_id) "+
		"JOIN product p ON (p.id = td.product_id) "+
		"JOIN unit u ON (u.id = p.unit_id) "+
		"%s ORDER BY t.date DESC, td.sorting_val ASC", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*transactiondomain.TransactionStatus
	entityMap := make(map[string]*transactiondomain.TransactionStatus)
	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var Date time.Time
		var StakeholderID sql.NullString
		var CustomerCode sql.NullString
		var CustomerName sql.NullString
		var TransactionType sql.NullString
		var Status sql.NullString
		var ReferenceCode sql.NullString
		var UserId sql.NullString
		var UserName sql.NullString
		var CreatedTime time.Time
		var UnitID sql.NullString
		var UnitCode sql.NullString
		var ProductID sql.NullString
		var ProductCode sql.NullString
		var ProductName sql.NullString
		var BuyPrice sql.NullFloat64
		var SellPrice sql.NullFloat64
		var Quantity sql.NullFloat64
		var BuyQuantity sql.NullFloat64

		rows.Scan(&ID, &Code, &Date, &StakeholderID, &CustomerCode, &CustomerName, &TransactionType, &Status, &ReferenceCode,
			&UserId, &UserName, &CreatedTime, &UnitID, &UnitCode, &ProductID, &ProductCode, &ProductName, &BuyPrice, &SellPrice, &Quantity, &BuyQuantity)

		var entity *transactiondomain.TransactionStatus
		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		if value, ok := entityMap[ID.String]; ok {
			entity = value
			entity.Total = entity.Total + ((SellPrice.Float64) * Quantity.Float64)
		} else {
			entity = &transactiondomain.TransactionStatus{}
			entity.ID = ID.String
			entity.Code = Code.String
			entity.Date = Date.Format(dateutil.DateFormatResponse())
			entity.StakeholderID = StakeholderID.String
			entity.StakeholderCode = CustomerCode.String
			entity.StakeholderName = CustomerName.String
			entity.TransactionType = TransactionType.String
			entity.Status = Status.String
			entity.ReferenceCode = ReferenceCode.String
			entity.UserId = UserId.String
			entity.UserName = UserName.String
			entity.CreatedTime = CreatedTime.Format(dateutil.TimeFormatResponse())
			entity.Total = (SellPrice.Float64) * Quantity.Float64

			entities = append(entities, entity)
			entityMap[ID.String] = entity
		}

		detail := &transactiondomain.TransactionStatusDetail{}
		detail.TransactionID = entity.ID
		detail.UnitID = UnitID.String
		detail.UnitCode = UnitCode.String
		detail.ProductID = ProductID.String
		detail.ProductCode = ProductCode.String
		detail.ProductName = ProductName.String
		detail.BuyPrice = BuyPrice.Float64
		detail.SellPrice = SellPrice.Float64
		detail.Quantity = Quantity.Float64
		detail.BuyQuantity = BuyQuantity.Float64

		entity.TransactionDetail = append(entity.TransactionDetail, detail)

	}

	return entities, nil
}

func (r *Repo) Create(entity *transactiondomain.Transaction, tx *gorm.DB) {
	tx.Exec("INSERT INTO public.transaction(id, code, date, stakeholder_id, transaction_type, status, reference_code, web_user_id, created_time) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
		entity.ID, entity.Code, entity.Date, entity.StakeholderID, entity.TransactionType, entity.Status, entity.ReferenceCode, entity.UserId, entity.CreatedTime)

	if tx.Error != nil {
		return
	}

	for i, detail := range entity.TransactionDetail {
		txDetailId := stringutil.GenerateUUID()
		detail.ID = txDetailId
		tx.Exec("INSERT INTO public.transaction_detail(id, transaction_id, product_id, buy_price, sell_price, quantity, created_time, web_user_id, latest, buy_quantity, sorting_val) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			txDetailId, entity.ID, detail.ProductID, detail.BuyPrice, detail.SellPrice, detail.Quantity, entity.CreatedTime, entity.UserId, true, detail.Quantity, i)

		if tx.Error != nil {
			return
		}
	}

}

func (r *Repo) Edit(entity *transactiondomain.Transaction) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.transaction "+
		"SET status=? WHERE id=?;", entity.Status, entity.ID)

	if tx.Error != nil {
		return tx.Error
	}

	for _, detail := range entity.TransactionDetail {
		tx.Exec("UPDATE public.transaction_detail "+
			"SET buy_price=?, sell_price=?, quantity=?, buy_quantity=? WHERE transaction_id=?, product_id=?;",
			detail.BuyPrice, detail.SellPrice, detail.Quantity, detail.BuyQuantity, detail.TransactionID, detail.ProductID)

		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	tx.Commit()
	return tx.Error
}

func (r *Repo) UpdateStatus(transactionID, status string) error {
	return global.DBCON.Exec("UPDATE public.transaction "+
		"SET status=? WHERE id=?;", status, transactionID).Error
}

func (r *Repo) UpdatePrice(transactionID, productID string, buyPrice, sellPrice float64, quantity, buyQuantity int64) error {
	return global.DBCON.Exec("UPDATE public.transaction_detail "+
		"SET buy_price=?, sell_price=?, quantity=?, buy_quantity=? WHERE transaction_id=?, product_id=?;",
		buyPrice, sellPrice, quantity, buyQuantity, transactionID, productID).Error
}

func (r *Repo) UpdateTransaction(entity *transactiondomain.Transaction, tx *gorm.DB) {
	tx.Exec("UPDATE public.transaction_detail SET latest=? WHERE transaction_id=?;", false, entity.ID)

	if tx.Error != nil {
		return
	}

	for i, detail := range entity.TransactionDetail {
		txDetailId := stringutil.GenerateUUID()
		detail.ID = txDetailId
		tx.Exec("INSERT INTO public.transaction_detail(id, transaction_id, product_id, buy_price, sell_price, quantity, created_time, web_user_id, latest, buy_quantity, sorting_val) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			txDetailId, entity.ID, detail.ProductID, detail.BuyPrice, detail.SellPrice, detail.Quantity, entity.CreatedTime, entity.UserId, true, detail.Quantity, i)

		if tx.Error != nil {
			return
		}
	}

}
