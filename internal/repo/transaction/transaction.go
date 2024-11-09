package transactionrepo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"dromatech/pos-backend/global"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	stringutil "dromatech/pos-backend/internal/util/string"
)

type TransactionRepo interface {
	Find(params []queryutil.Param) ([]*transactiondomain.Transaction, error)
	Create(entity *transactiondomain.Transaction, tx *gorm.DB)
	Edit(product *transactiondomain.Transaction) error
	UpdateStatus(transactionID, status string) error
	UpdatePrice(transactionID, productID string, buyPrice, sellPrice float64, quantity, buyQuantity int64) error
	FindSells(params []queryutil.Param) ([]*transactiondomain.TransactionStatus, error)
	UpdateTransaction(entity *transactiondomain.Transaction, tx *gorm.DB)
	FindReport(params []queryutil.Param) ([]*transactiondomain.ReportDate, error)
	UpdateHargaBeli(transactionDetailID string, buyPrice int64, webUserID string) error
	InsertTransactionBuy(transactionId string, transactionBuy []transactiondomain.TransactionBuy) error
	FindTransactionBuyStatus() ([]transactiondomain.TransactionBuyStatus, error)
	FindDetails(params []queryutil.Param) ([]*transactiondomain.TransactionDetail, error)
	UpdateHargaBeliTx(transactionDetailID string, buyPrice float64, webUserID string, tx *gorm.DB) error
	FindLastCreditPerMonth(params []queryutil.Param) (map[string]map[int]float64, error)
	FindLastCredit(params []queryutil.Param) (map[string]float64, error)
	FindCustomerReport(stakeHolderID string, month time.Time) ([]*transactiondomain.LaporanCustomer, int, error)
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
		// hard set to 0 when creating new transaction
		detail.BuyPrice = 0
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
			txDetailId, entity.ID, detail.ProductID, detail.BuyPrice, detail.SellPrice, detail.Quantity, entity.CreatedTime, entity.UserId, true, detail.BuyQuantity, i)

		if tx.Error != nil {
			return
		}
	}

}

func (r *Repo) FindReport(params []queryutil.Param) ([]*transactiondomain.ReportDate, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT t.date, t.code, t.reference_code, t.status, p.code, p.name, td.buy_quantity, td.buy_price, td.quantity, td.sell_price, td.id "+
		"FROM transaction t "+
		"JOIN transaction_detail td ON (t.id = td.transaction_id) "+
		"JOIN product p ON (td.product_id = p.id) "+
		"JOIN customer c ON (t.stakeholder_id = c.id) "+
		"%s ORDER BY t.date DESC, t.code ASC, p.id ASC", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var dateEntities []*transactiondomain.ReportDate
	dateMap := make(map[string]*transactiondomain.ReportDate)
	entityMap := make(map[string]*transactiondomain.Report)
	for rows.Next() {
		var Date time.Time
		var Code sql.NullString
		var ReferenceCode sql.NullString
		var Status sql.NullString
		var ProductCode sql.NullString
		var ProductName sql.NullString
		var BuyQuantity sql.NullFloat64
		var BuyPrice sql.NullFloat64
		var Quantity sql.NullFloat64
		var SellPrice sql.NullFloat64
		var ID sql.NullString

		rows.Scan(&Date, &Code, &ReferenceCode, &Status, &ProductCode, &ProductName, &BuyQuantity, &BuyPrice, &Quantity, &SellPrice, &ID)

		var dateEntity *transactiondomain.ReportDate
		var entity *transactiondomain.Report
		if !Code.Valid && Code.String == "" {
			return nil, nil
		}

		dateString := Date.Format(dateutil.DateFormatResponse())
		if value, ok := dateMap[dateString]; ok {
			dateEntity = value
		} else {
			dateEntity = &transactiondomain.ReportDate{}
			dateEntity.Date = dateString

			dateEntities = append(dateEntities, dateEntity)
			dateMap[dateString] = dateEntity
		}

		if value, ok := entityMap[Code.String]; ok {
			entity = value
		} else {
			entity = &transactiondomain.Report{}
			entity.Code = Code.String
			entity.ReferenceCode = ReferenceCode.String
			entity.Status = Status.String

			dateEntity.Reports = append(dateEntity.Reports, entity)
			entityMap[entity.Code] = entity
		}

		detail := &transactiondomain.ReportDetail{}
		detail.ProductCode = ProductCode.String
		detail.ProductName = ProductName.String
		detail.BuyQuantity = BuyQuantity.Float64
		detail.BuyPrice = BuyPrice.Float64
		detail.Quantity = Quantity.Float64
		detail.SellPrice = SellPrice.Float64
		detail.ID = ID.String

		entity.ReportDetails = append(entity.ReportDetails, detail)

	}

	return dateEntities, nil

}
func (r *Repo) UpdateHargaBeli(transactionDetailID string, buyPrice int64, webUserID string) error {
	tx := global.DBCON.Begin()
	tx.Exec("UPDATE public.transaction_detail SET latest=? WHERE id=?;", false, transactionDetailID)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("INSERT INTO public.transaction_detail(id, transaction_id, product_id, buy_price, sell_price, quantity, created_time, web_user_id, latest, buy_quantity, sorting_val) "+
		"SELECT ?, transaction_id, product_id, ?, sell_price, quantity, ?, ?, ?, buy_quantity, sorting_val "+
		"FROM public.transaction_detail "+
		"WHERE id=?;", stringutil.GenerateUUID(), buyPrice, time.Now(), webUserID, true, transactionDetailID)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

func (r *Repo) UpdateHargaBeliTx(transactionDetailID string, buyPrice float64, webUserID string, tx *gorm.DB) error {
	tx.Exec("UPDATE public.transaction_detail SET latest=? WHERE id=?;", false, transactionDetailID)

	if tx.Error != nil {
		return tx.Error
	}

	tx.Exec("INSERT INTO public.transaction_detail(id, transaction_id, product_id, buy_price, sell_price, quantity, created_time, web_user_id, latest, buy_quantity, sorting_val) "+
		"SELECT ?, transaction_id, product_id, ?, sell_price, quantity, ?, ?, ?, buy_quantity, sorting_val "+
		"FROM public.transaction_detail "+
		"WHERE id=?;", stringutil.GenerateUUID(), buyPrice, time.Now(), webUserID, true, transactionDetailID)

	return tx.Error
}

func (r *Repo) InsertTransactionBuy(transactionId string, transactionBuys []transactiondomain.TransactionBuy) error {
	tx := global.DBCON.Begin()

	tx.Exec("UPDATE public.transaction_buy SET latest=? WHERE transaction_id=?;", false, transactionId)

	if tx.Error != nil {
		return tx.Error
	}

	for _, transactionBuy := range transactionBuys {
		tx.Exec("INSERT INTO public.transaction_buy (id, transaction_id, product_id, price, quantity, payment_method, created_time, web_user_id, latest) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);", transactionBuy.ID, transactionBuy.TransactionID, transactionBuy.ProductID, transactionBuy.Price, transactionBuy.Quantity, transactionBuy.PaymentMethod, transactionBuy.CreatedTime, transactionBuy.WebUserID, true)

		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	return nil
}

func (r *Repo) FindTransactionBuyStatus() ([]transactiondomain.TransactionBuyStatus, error) {
	rows, err := global.DBCON.Raw("SELECT t.id, t.code, c.code, c.name, count(tb.id), count(td.id) " +
		"FROM transaction t " +
		"JOIN customer c ON (c.id = t.stakeholder_id) " +
		"LEFT JOIN transaction_buy tb ON (tb.transaction_id = t.id) " +
		"JOIN transaction_detail td ON (td.transaction_id = t.id) " +
		"WHERE (tb.id IS NULL OR tb.latest = true) AND td.latest = true " +
		"GROUP BY t.id, t.code, c.code, c.name").Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []transactiondomain.TransactionBuyStatus

	for rows.Next() {
		var ID sql.NullString
		var Code sql.NullString
		var CustomerCode sql.NullString
		var CustomerName sql.NullString
		var TotalBuy sql.NullInt64
		var TotalSell sql.NullInt64

		rows.Scan(&ID, &Code, &CustomerCode, &CustomerName, &TotalBuy, &TotalSell)

		var entity transactiondomain.TransactionBuyStatus
		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		entity.ID = ID.String
		entity.Code = Code.String
		entity.CustomerCode = CustomerCode.String
		entity.CustomerName = CustomerName.String
		entity.TotalBuy = TotalBuy.Int64
		entity.TotalSell = TotalSell.Int64

		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *Repo) FindDetails(params []queryutil.Param) ([]*transactiondomain.TransactionDetail, error) {
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

	rows, err := global.DBCON.Raw(fmt.Sprintf("SELECT td.id, td.transaction_id, td.product_id, td.buy_price, td.sell_price, td.quantity, "+
		"td.buy_quantity, td.created_time, td.web_user_id, td.latest, td.sorting_val "+
		"FROM transaction_detail td "+
		"JOIN transaction t ON (td.transaction_id = t.id) "+
		"%s ", where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var entities []*transactiondomain.TransactionDetail
	for rows.Next() {
		var ID sql.NullString
		var TransactionID sql.NullString
		var ProductID sql.NullString
		var BuyPrice sql.NullFloat64
		var SellPrice sql.NullFloat64
		var Quantity sql.NullFloat64
		var BuyQuantity sql.NullFloat64
		var CreatedTime sql.NullTime
		var WebUserID sql.NullString
		var Latest sql.NullBool
		var SortingVal sql.NullInt64

		err = rows.Scan(&ID, &TransactionID, &ProductID, &BuyPrice, &SellPrice, &Quantity, &BuyQuantity, &CreatedTime, &WebUserID, &Latest, &SortingVal)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}

		if !ID.Valid && ID.String == "" {
			return nil, nil
		}

		detail := &transactiondomain.TransactionDetail{}
		detail.ID = ID.String
		detail.TransactionID = TransactionID.String
		detail.ProductID = ProductID.String
		detail.BuyPrice = BuyPrice.Float64
		detail.SellPrice = SellPrice.Float64
		detail.Quantity = Quantity.Float64
		detail.BuyQuantity = BuyQuantity.Float64
		detail.CreatedTime = CreatedTime.Time
		detail.WebUserID = WebUserID.String
		detail.Latest = Latest.Bool
		detail.SortingVal = SortingVal.Int64

		entities = append(entities, detail)

	}

	return entities, nil
}

func (r *Repo) FindLastCredit(params []queryutil.Param) (map[string]float64, error) {
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

	query := "SELECT c.code, SUM(td.quantity * td.sell_price) AS \"balance\" FROM transaction t " +
		"JOIN transaction_detail td ON (t.id = td.transaction_id) " +
		"JOIN customer c ON (t.stakeholder_id = c.id) " +
		"%s GROUP BY c.code ORDER BY c.code;"

	rows, err := global.DBCON.Raw(fmt.Sprintf(query, where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	entities := make(map[string]float64)
	for rows.Next() {
		var CustomerCode sql.NullString
		var LastCredit sql.NullFloat64

		err = rows.Scan(&CustomerCode, &LastCredit)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}

		if !CustomerCode.Valid && CustomerCode.String == "" {
			return nil, nil
		}

		entities[CustomerCode.String] = LastCredit.Float64
	}

	return entities, nil
}

func (r *Repo) FindLastCreditPerMonth(params []queryutil.Param) (map[string]map[int]float64, error) {
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

	query := "SELECT c.code, t.date, SUM(td.quantity * td.sell_price) FROM transaction t " +
		"JOIN transaction_detail td ON (t.id = td.transaction_id) " +
		"JOIN customer c ON (t.stakeholder_id = c.id) " +
		"%s GROUP BY c.code, t.date ORDER BY c.code, t.date;"

	rows, err := global.DBCON.Raw(fmt.Sprintf(query, where), values...).Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	entities := make(map[string]map[int]float64)
	for rows.Next() {
		var CustomerCode sql.NullString
		var Date time.Time
		var Credit sql.NullFloat64

		err = rows.Scan(&CustomerCode, &Date, &Credit)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}

		if !CustomerCode.Valid && CustomerCode.String == "" {
			return nil, nil
		}

		if _, ok := entities[CustomerCode.String]; !ok {
			entities[CustomerCode.String] = make(map[int]float64)
			entities[CustomerCode.String] = make(map[int]float64)
		}

		entities[CustomerCode.String][Date.Day()] = Credit.Float64

	}

	return entities, nil
}

func (r *Repo) FindCustomerReport(stakeHolderID string, month time.Time) ([]*transactiondomain.LaporanCustomer, int, error) {
	query := `select p.code, p.name, count(p.*), t.date, t.id 
	from public.transaction_detail td 
	join public.transaction t on t.id = td.transaction_id 
	join public.product p on p.id = td.product_id 
	where t.transaction_type = 'SELL' 
	and (t.status = 'KONTRABON' or t.status = 'DIBAYAR') 
	and t.date >= ? 
	and t.date <= ?
	and t.stakeholder_id = ?
	group by p.code, p.name, t.date, t.stakeholder_id, t.id 
	order by p.name, t.date`

	startDate := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDate := time.Date(month.Year(), month.Month(), dateutil.DaysIn(month.Month(), month.Year()), 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	rows, err := global.DBCON.Raw(query, startDate, endDate, stakeHolderID).Rows()

	if err != nil {
		logrus.Error(err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	entities := make(map[string]*transactiondomain.LaporanCustomer)
	txId := make(map[string]struct{})
	var result []*transactiondomain.LaporanCustomer = make([]*transactiondomain.LaporanCustomer, 0)
	for rows.Next() {
		var ProductCode sql.NullString
		var ProductName sql.NullString
		var Count sql.NullInt32
		var Date time.Time
		var TxID sql.NullString

		err = rows.Scan(&ProductCode, &ProductName, &Count, &Date, &TxID)
		if err != nil {
			logrus.Error(err.Error())
			return nil, 0, err
		}

		if !ProductCode.Valid && ProductCode.String == "" {
			return nil, 0, nil
		}

		if _, ok := entities[ProductCode.String]; !ok {
			entity := &transactiondomain.LaporanCustomer{
				ProductCode: ProductCode.String,
				ProductName: ProductName.String,
				Counts:      make(map[int]int),
			}

			entities[ProductCode.String] = entity
			result = append(result, entity)
		}

		entities[ProductCode.String].Counts[Date.Day()] = int(Count.Int32)
		txId[TxID.String] = struct{}{}
	}

	return result, len(txId), nil

}
