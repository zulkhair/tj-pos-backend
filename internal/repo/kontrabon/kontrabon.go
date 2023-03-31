package kontrabonrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	kontrabondomain "dromatech/pos-backend/internal/domain/kontrabon"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type KontrabonRepo interface {
	Find(params []queryutil.Param) ([]*kontrabondomain.Kontrabon, error)
	FindTransaction(params []queryutil.Param) ([]*transactiondomain.TransactionStatus, error)
	Create(entity kontrabondomain.Kontrabon, transactionIds []string) error
	Update(kontrabonId string, transactionIds []string, status string) error
	CreateTx(entity kontrabondomain.Kontrabon, transactionIds []string, tx *gorm.DB)
	UpdateLunas(kontrabonId string) error
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) Find(params []queryutil.Param) ([]*kontrabondomain.Kontrabon, error) {
	params = append(params, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    "TRUE",
	})

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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT k.id, k.code, k.created_time, k.status, SUM(td.sell_price * td.quantity), customer_id "+
		"FROM public.kontrabon k "+
		"JOIN public.kontrabon_transaction kt ON (kt.kontrabon_id = k.id) "+
		"JOIN public.transaction t ON (t.id = kt.transaction_id) "+
		"JOIN public.transaction_detail td ON (td.transaction_id = t.id)"+
		"%s GROUP BY k.id, k.code, k.created_time, k.status ORDER BY k.code DESC, k.created_time DESC", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*kontrabondomain.Kontrabon
	entityMap := make(map[string]*kontrabondomain.Kontrabon)
	for rows.Next() {

		var ID sql.NullString
		var Code sql.NullString
		var CreatedTime time.Time
		var Status sql.NullString
		var Total sql.NullFloat64
		var CustomerID sql.NullString

		rows.Scan(&ID, &Code, &CreatedTime, &Status, &Total, &CustomerID)

		var entity *kontrabondomain.Kontrabon
		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		if value, ok := entityMap[ID.String]; ok {
			entity = value
			entity.Total = entity.Total + Total.Float64
		} else {
			entity = &kontrabondomain.Kontrabon{}
			entity.ID = ID.String
			entity.Code = Code.String
			entity.CreatedTime = CreatedTime.Format(dateutil.DateFormatResponse())
			entity.Status = Status.String
			entity.Total = Total.Float64
			entity.CustomerID = CustomerID.String

			entities = append(entities, entity)
		}

	}

	return entities, nil
}

func (r *Repo) FindTransaction(params []queryutil.Param) ([]*transactiondomain.TransactionStatus, error) {
	params = append(params, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    "TRUE",
	})

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
		"u.id, u.code, td.product_id, p.code, p.name, td.buy_price, td.sell_price, td.quantity "+
		"FROM transaction t "+
		"LEFT JOIN kontrabon_transaction kt ON (kt.transaction_id = t.id) "+
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

		rows.Scan(&ID, &Code, &Date, &StakeholderID, &CustomerCode, &CustomerName, &TransactionType, &Status, &ReferenceCode,
			&UserId, &UserName, &CreatedTime, &UnitID, &UnitCode, &ProductID, &ProductCode, &ProductName, &BuyPrice, &SellPrice, &Quantity)

		var entity *transactiondomain.TransactionStatus
		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		if value, ok := entityMap[ID.String]; ok {
			entity = value
			entity.Total = entity.Total + (SellPrice.Float64)*Quantity.Float64
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

		entity.TransactionDetail = append(entity.TransactionDetail, detail)

	}

	return entities, nil
}

func (r *Repo) Create(entity kontrabondomain.Kontrabon, transactionIds []string) error {
	tx := global.DBCON.Begin()
	r.CreateTx(entity, transactionIds, tx)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *Repo) CreateTx(entity kontrabondomain.Kontrabon, transactionIds []string, tx *gorm.DB) {
	tx.Exec("INSERT INTO public.kontrabon(id, code, created_time, status, customer_id) VALUES (?, ?, ?, ?, ?)", entity.ID, entity.Code, entity.CreatedTime, entity.Status, entity.CustomerID)
	if tx.Error != nil {
		return
	}

	for _, transactionId := range transactionIds {
		tx.Exec("INSERT INTO public.kontrabon_transaction(kontrabon_id, transaction_id) VALUES (?, ?);", entity.ID, transactionId)
		if tx.Error != nil {
			return
		}
		tx.Exec("UPDATE public.transaction SET status=? WHERE id=?;", transactiondomain.TRANSACTION_KONTRABON, transactionId)
		if tx.Error != nil {
			return
		}
	}
}

func (r *Repo) Update(kontrabonId string, transactionIds []string, status string) error {
	tx := global.DBCON.Begin()

	for _, transactionId := range transactionIds {
		if status == transactiondomain.TRANSACTION_KONTRABON {
			tx.Exec("INSERT INTO public.kontrabon_transaction(kontrabon_id, transaction_id) VALUES (?, ?);", kontrabonId, transactionId)
			if tx.Error != nil {
				tx.Rollback()
				return tx.Error
			}
		} else if status == transactiondomain.TRANSACTION_PEMBUATAN {
			tx.Exec("DELETE FROM public.kontrabon_transaction WHERE kontrabon_id=? AND transaction_id=?;", kontrabonId, transactionId)
			if tx.Error != nil {
				tx.Rollback()
				return tx.Error
			}
		}

		tx.Exec("UPDATE public.transaction SET status=? WHERE id=?;", status, transactionId)
		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	tx.Commit()
	return nil
}

func (r *Repo) UpdateLunas(kontrabonId string) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.kontrabon SET status=? WHERE id=?", kontrabondomain.STATUS_LUNAS, kontrabonId)
	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("UPDATE public.transaction SET status=? WHERE id IN (SELECT kt.transaction_id FROM kontrabon_transaction kt WHERE kt.kontrabon_id = ?)", transactiondomain.TRANSACTION_DIBAYAR, kontrabonId)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}
