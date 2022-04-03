package transactionrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	queryutil "dromatech/pos-backend/internal/util/query"
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
	UpdateBuyPrice(transactionID, unitID, productID string, price float64) error
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
		"td.unit_id, td.product_id, td.price, td.quantity"+
		"FROM transaction t "+
		"JOIN transaction_detail td"+
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
		var Price sql.NullFloat64
		var Quantity sql.NullInt64

		rows.Scan(&ID, &Code, &Date, &StakeholderID, &TransactionType, &Status, &UnitID, &ProductID, &Price, &Quantity)

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
			entity.Date = Date
			entity.StakeholderID = StakeholderID.String
			entity.TransactionType = TransactionType.String
			entity.Status = Status.String

			entities = append(entities, entity)
		}

		detail := &transactiondomain.TransactionDetail{}
		detail.TransactionID = entity.ID
		detail.UnitID = UnitID.String
		detail.ProductID = ProductID.String
		detail.Price = Price.Float64
		detail.Quantity = Quantity.Int64

		entity.TransactionDetail = append(entity.TransactionDetail, detail)

	}

	return entities, nil
}

func (r *Repo) Create(entity *transactiondomain.Transaction, tx *gorm.DB) {
	tx.Exec("INSERT INTO public.transaction(id, code, date, stakeholder_id, transaction_type, status) "+
		"VALUES (?, ?, ?, ?, ?, ?);",
		entity.ID, entity.Code, entity.Date, entity.StakeholderID, entity.TransactionType, entity.Status)

	if tx.Error != nil {
		return
	}

	for _, detail := range entity.TransactionDetail {
		tx.Exec("INSERT INTO public.transaction_detail(transaction_id, unit_id, product_id, price, quantity) "+
			"VALUES (?, ?, ?, ?, ?);",
			entity.ID, detail.UnitID, detail.ProductID, detail.Price, detail.Quantity)

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
			"SET price=? WHERE transaction_id=?, unit_id=?, product_id=?;",
			detail.Price, detail.TransactionID, detail.UnitID, detail.ProductID)

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

func (r *Repo) UpdateBuyPrice(transactionID, unitID, productID string, price float64) error {
	return global.DBCON.Exec("UPDATE public.transaction_detail "+
		"SET price=? WHERE transaction_id=?, unit_id=?, product_id=?;",
		price, transactionID, unitID, productID).Error
}
