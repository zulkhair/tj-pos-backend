package transactionusecase

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	customerrepo "dromatech/pos-backend/internal/repo/customer"
	kontrabonrepo "dromatech/pos-backend/internal/repo/kontrabon"
	sequencerepo "dromatech/pos-backend/internal/repo/sequence"
	supplierrepo "dromatech/pos-backend/internal/repo/supplier"
	transactionrepo "dromatech/pos-backend/internal/repo/transaction"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	stringutil "dromatech/pos-backend/internal/util/string"
)

type TransactoionUsecase interface {
	CreateTransaction(transaction *transactiondomain.Transaction) (string, error)
	ViewTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.Transaction, error)
	UpdateStatus(transactionID string) error
	UpdateBuyPrice(transactionID, productID string, buyPrice, sellPrice float64, quantity, buyQuantity int64) error
	ViewSellTransaction(startDate, endDate, code, stakeholderID, txType, status, productID, txId string) ([]*transactiondomain.TransactionStatus, error)
	CancelTrx(transactionID string) error
	UpdateTransaction(transaction *transactiondomain.Transaction) error
	FindReport(startDate, endDate, code, stakeholderID, txType, status, productID, txId string) ([]*transactiondomain.ReportDate, error)
	UpdateHargaBeli(request transactiondomain.UpdateHargaBeliRequest) error
	InsertTransactionBuy(request transactiondomain.InsertTransactionBuyRequestBulk) error
	FindCustomerCredit(month time.Time, sell bool) (*transactiondomain.TransactionCredit, error)
	FindCustomerReport(stakeholderId string, month time.Time) (*transactiondomain.LaporanCustomerSumary, error)
}

type Usecase struct {
	transactionRepo transactionrepo.TransactionRepo
	sequenceRepo    sequencerepo.SequenceRepo
	supplierRepo    supplierrepo.SupplierRepo
	customerRepo    customerrepo.CustomerRepo
	kontrabonRepo   kontrabonrepo.KontrabonRepo
}

func New(transactionRepo transactionrepo.TransactionRepo, sequenceRepo sequencerepo.SequenceRepo, supplierRepo supplierrepo.SupplierRepo, customerRepo customerrepo.CustomerRepo, kontrabonRepo kontrabonrepo.KontrabonRepo) *Usecase {
	uc := &Usecase{
		transactionRepo: transactionRepo,
		sequenceRepo:    sequenceRepo,
		supplierRepo:    supplierRepo,
		customerRepo:    customerRepo,
		kontrabonRepo:   kontrabonRepo,
	}

	return uc
}

func (uc *Usecase) CreateTransaction(transaction *transactiondomain.Transaction) (string, error) {
	tx := global.DBCON.Begin()

	timeNow := time.Now().UTC()
	if transaction.Date == "" {
		transaction.Date = timeNow.Format(dateutil.DateFormat())
	}
	stakeHolderCode := ""
	if transaction.TransactionType == transactiondomain.TRANSACTION_TYPE_BUY {
		supplier, err := uc.supplierRepo.Find(map[string]interface{}{"id": transaction.StakeholderID})
		if err != nil {
			logrus.Error(err.Error())
			return "", fmt.Errorf("Supplier dengan kode '%s' tidak ditemukan", transaction.StakeholderID)
		}
		stakeHolderCode = supplier[0].Code
	} else {
		customer, err := uc.customerRepo.Find(map[string]interface{}{"id": transaction.StakeholderID})
		if err != nil {
			logrus.Error(err.Error())
			return "", fmt.Errorf("Customer dengan kode '%s' tidak ditemukan", transaction.StakeholderID)
		}
		stakeHolderCode = customer[0].Code
	}
	transactionID := strings.ReplaceAll(uuid.NewString(), "-", "")
	dateCode, err := time.Parse(dateutil.DateFormat(), transaction.Date)
	if err != nil {
		logrus.Error(err.Error())
		return "", fmt.Errorf("Terjadi kesalahan saat melakukan transaksi")
	}

	seqcode := stakeHolderCode + "/" + dateCode.Format("2006")
	seq := uc.sequenceRepo.NextValTx(seqcode, tx)
	transactionCode := stakeHolderCode + "/" + stringutil.ToRoman(int(dateCode.Month())) + "/" + dateCode.Format("2006")
	transactionCode = strconv.Itoa(int(seq)) + "/" + transactionCode

	transaction.ID = transactionID
	transaction.Code = transactionCode
	transaction.Status = transactiondomain.TRANSACTION_PEMBUATAN
	transaction.CreatedTime = timeNow

	uc.transactionRepo.Create(transaction, tx)

	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return "", fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
	}

	for _, detail := range transaction.TransactionDetail {
		// set buy price never be filled from creating transaction
		detail.BuyPrice = 0

		// cek dan update harga beli
		buyPriceOlds, err := uc.supplierRepo.FindBuyPrice([]queryutil.Param{
			{
				Logic:    "AND",
				Field:    "s.product_id",
				Operator: "=",
				Value:    detail.ProductID,
			},
			{
				Logic:    "AND",
				Field:    "s.latest",
				Operator: "=",
				Value:    strconv.FormatBool(true),
			},
		})

		if err != nil {
			tx.Rollback()
			logrus.Error(err.Error())
			return "", fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
		}

		if buyPriceOlds != nil && len(buyPriceOlds) > 0 {
			if buyPriceOlds[0].Price != detail.BuyPrice {
				var tId *string
				if transaction.ID != "" {
					tId = &transaction.ID
				} else {
					tId = nil
				}
				uc.supplierRepo.AddBuyPriceTx(supplierdomain.AddPriceRequest{
					ID:            strings.ReplaceAll(uuid.NewString(), "-", ""),
					Date:          timeNow,
					UnitId:        detail.UnitID,
					ProductID:     detail.ProductID,
					Price:         detail.BuyPrice,
					WebUserId:     transaction.UserId,
					Latest:        true,
					TransactionId: tId,
				}, tx)
			}
		} else {
			var tId *string
			if transaction.ID != "" {
				tId = &transaction.ID
			} else {
				tId = nil
			}
			uc.supplierRepo.AddBuyPriceTx(supplierdomain.AddPriceRequest{
				ID:            strings.ReplaceAll(uuid.NewString(), "-", ""),
				Date:          timeNow,
				UnitId:        detail.UnitID,
				ProductID:     detail.ProductID,
				Price:         detail.BuyPrice,
				WebUserId:     transaction.UserId,
				Latest:        true,
				TransactionId: tId,
			}, tx)
		}

		// cek dan update harga jual
		sellPriceOlds, err := uc.customerRepo.FindSellPrice([]queryutil.Param{
			{
				Logic:    "AND",
				Field:    "s.customer_id",
				Operator: "=",
				Value:    transaction.StakeholderID,
			},
			{
				Logic:    "AND",
				Field:    "s.product_id",
				Operator: "=",
				Value:    detail.ProductID,
			},
			{
				Logic:    "AND",
				Field:    "s.latest",
				Operator: "=",
				Value:    strconv.FormatBool(true),
			},
		})

		if err != nil {
			tx.Rollback()
			logrus.Error(err.Error())
			return "", fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
		}

		if sellPriceOlds != nil && len(sellPriceOlds) > 0 {
			if sellPriceOlds[0].Price != detail.SellPrice {
				var tId *string
				if transaction.ID != "" {
					tId = &transaction.ID
				} else {
					tId = nil
				}
				uc.customerRepo.AddSellPriceTx(customerdomain.AddPriceRequest{
					ID:            strings.ReplaceAll(uuid.NewString(), "-", ""),
					Date:          timeNow,
					UnitId:        detail.UnitID,
					ProductID:     detail.ProductID,
					Price:         detail.SellPrice,
					WebUserId:     transaction.UserId,
					Latest:        true,
					TransactionId: tId,
					CustomerId:    transaction.StakeholderID,
				}, tx)
			}
		} else {
			var tId *string
			if transaction.ID != "" {
				tId = &transaction.ID
			} else {
				tId = nil
			}
			uc.customerRepo.AddSellPriceTx(customerdomain.AddPriceRequest{
				ID:            strings.ReplaceAll(uuid.NewString(), "-", ""),
				Date:          timeNow,
				UnitId:        detail.UnitID,
				ProductID:     detail.ProductID,
				Price:         detail.SellPrice,
				WebUserId:     transaction.UserId,
				Latest:        true,
				TransactionId: tId,
				CustomerId:    transaction.StakeholderID,
			}, tx)
		}
	}

	tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return "", fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
	}
	return transaction.ID, nil
}

func (uc *Usecase) ViewTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.Transaction, error) {
	var param []queryutil.Param
	if startDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.date",
			Operator: ">=",
			Value:    startDate,
		})
	}
	if endDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.date",
			Operator: "<=",
			Value:    endDate,
		})
	}
	if code != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.code",
			Operator: "=",
			Value:    code,
		})
	}
	if stakeholderID != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.stakeholder_id",
			Operator: "=",
			Value:    stakeholderID,
		})
	}
	if txType != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.type",
			Operator: "=",
			Value:    txType,
		})
	}
	if status != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.status",
			Operator: "=",
			Value:    status,
		})
	}
	if productID != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "td.product_id",
			Operator: "=",
			Value:    productID,
		})
	}

	return uc.transactionRepo.Find(param)
}

func (uc *Usecase) ViewSellTransaction(startDate, endDate, code, stakeholderID, txType, status, productID, txId string) ([]*transactiondomain.TransactionStatus, error) {
	var param []queryutil.Param
	if startDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.date",
			Operator: ">=",
			Value:    startDate,
		})
	}
	if endDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.date",
			Operator: "<=",
			Value:    endDate,
		})
	}
	if code != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.code",
			Operator: "=",
			Value:    code,
		})
	}
	if stakeholderID != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.stakeholder_id",
			Operator: "=",
			Value:    stakeholderID,
		})
	}
	if txType != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.transaction_type",
			Operator: "=",
			Value:    txType,
		})
	}
	if status != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.status",
			Operator: "=",
			Value:    status,
		})
	}
	if productID != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "td.product_id",
			Operator: "=",
			Value:    productID,
		})
	}

	if txId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.id",
			Operator: "=",
			Value:    txId,
		})
	}

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.status",
		Operator: "<>",
		Value:    "BATAL",
	})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    "TRUE",
	})

	return uc.transactionRepo.FindSells(param)
}

func (uc *Usecase) UpdateStatus(transactionID string) error {
	var param []queryutil.Param
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.id",
		Operator: "=",
		Value:    transactionID,
	})

	transactions, err := uc.transactionRepo.FindSells(param)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	if transactions == nil || len(transactions) == 0 {
		logrus.Error("Transactions is nil or zero")
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	status := transactions[0].Status
	if status == transactiondomain.TRANSACTION_PEMBUATAN {
		status = transactiondomain.TRANSACTION_KONTRABON
	} else if transactions[0].Status == transactiondomain.TRANSACTION_KONTRABON {
		status = transactiondomain.TRANSACTION_DIBAYAR
	}
	err = uc.transactionRepo.UpdateStatus(transactionID, status)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}
	return nil
}

func (uc *Usecase) CancelTrx(transactionID string) error {
	var param []queryutil.Param
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.id",
		Operator: "=",
		Value:    transactionID,
	})

	transactions, err := uc.transactionRepo.FindSells(param)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	if transactions == nil || len(transactions) == 0 {
		logrus.Error("Transactions is nil or zero")
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	err = uc.transactionRepo.UpdateStatus(transactionID, transactiondomain.TRANSACTION_BATAL)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}
	return nil
}

func (uc *Usecase) UpdateBuyPrice(transactionID, productID string, buyPrice, sellPrice float64, quantity, buyQuantity int64) error {
	err := uc.transactionRepo.UpdatePrice(transactionID, productID, buyPrice, sellPrice, quantity, buyQuantity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}
	return nil
}

func (uc *Usecase) UpdateTransaction(transaction *transactiondomain.Transaction) error {
	tx := global.DBCON.Begin()

	timeNow := time.Now().UTC()
	transaction.CreatedTime = timeNow

	uc.transactionRepo.UpdateTransaction(transaction, tx)
	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return fmt.Errorf("Terjadi kesalahan saat memperbarui transaksi")
	}

	tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return fmt.Errorf("Terjadi kesalahan saat memperbarui transaksi")
	}
	return nil
}

func (uc *Usecase) FindReport(startDate, endDate, code, stakeholderID, txType, status, productID, txId string) ([]*transactiondomain.ReportDate, error) {
	var param []queryutil.Param
	if startDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.date",
			Operator: ">=",
			Value:    startDate,
		})
	}
	if endDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.date",
			Operator: "<=",
			Value:    endDate,
		})
	}
	//if code != "" {
	//	param = append(param, queryutil.Param{
	//		Logic:    "AND",
	//		Field:    "t.code",
	//		Operator: "=",
	//		Value:    code,
	//	})
	//}
	if stakeholderID != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.stakeholder_id",
			Operator: "=",
			Value:    stakeholderID,
		})
	}
	if productID != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "p.id",
			Operator: "=",
			Value:    productID,
		})
	}
	//if txType != "" {
	//	param = append(param, queryutil.Param{
	//		Logic:    "AND",
	//		Field:    "t.transaction_type",
	//		Operator: "=",
	//		Value:    txType,
	//	})
	//}
	if status != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.status",
			Operator: "=",
			Value:    status,
		})
	} else {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "t.status",
			Operator: "IN",
			Value:    []string{transactiondomain.TRANSACTION_DICETAK, transactiondomain.TRANSACTION_PEMBUATAN, transactiondomain.TRANSACTION_KONTRABON, transactiondomain.TRANSACTION_DIBAYAR},
		})
	}
	//if productID != "" {
	//	param = append(param, queryutil.Param{
	//		Logic:    "AND",
	//		Field:    "td.product_id",
	//		Operator: "=",
	//		Value:    productID,
	//	})
	//}

	//if txId != "" {
	//	param = append(param, queryutil.Param{
	//		Logic:    "AND",
	//		Field:    "t.id",
	//		Operator: "=",
	//		Value:    txId,
	//	})
	//}

	//param = append(param, queryutil.Param{
	//	Logic:    "AND",
	//	Field:    "t.status",
	//	Operator: "<>",
	//	Value:    "BATAL",
	//})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    "TRUE",
	})

	return uc.transactionRepo.FindReport(param)
}

func (uc *Usecase) UpdateHargaBeli(request transactiondomain.UpdateHargaBeliRequest) error {
	err := uc.transactionRepo.UpdateHargaBeli(request.TransactionDetailID, request.BuyPrice, request.WebUserID)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat memperbarui harga beli")
	}
	return nil
}

func (uc *Usecase) InsertTransactionBuy(request transactiondomain.InsertTransactionBuyRequestBulk) error {
	now := time.Now()
	var entities []transactiondomain.TransactionBuy
	for _, detail := range request.Details {
		entity := transactiondomain.TransactionBuy{
			ID:            strings.ReplaceAll(uuid.NewString(), "-", ""),
			TransactionID: request.TransactionID,
			ProductID:     detail.ProductID,
			Quantity:      detail.Quantity,
			Price:         detail.Price,
			PaymentMethod: detail.PaymentMethod,
			CreatedTime:   now,
			WebUserID:     request.WebUserID,
		}

		entities = append(entities, entity)
	}

	err := uc.transactionRepo.InsertTransactionBuy(request.TransactionID, entities)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat menambahkan harga beli")
	}
	return nil
}

func (uc *Usecase) FindTransactionBuyStatus() ([]transactiondomain.TransactionBuyStatus, error) {
	return uc.transactionRepo.FindTransactionBuyStatus()
}

func (uc *Usecase) FindCustomerCredit(month time.Time, sell bool) (*transactiondomain.TransactionCredit, error) {
	transactionCredit := &transactiondomain.TransactionCredit{}

	status := []string{"BATAL"}
	if !sell {
		status = append(status, "DIBAYAR")
	}

	var param []queryutil.Param
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.date",
		Operator: ">=",
		Value:    "2023-09-01",
	})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.date",
		Operator: "<",
		Value:    time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.status",
		Operator: "NOT IN",
		Value:    status,
	})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    true,
	})

	lastCreditMap, err := uc.transactionRepo.FindLastCredit(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat mengambil data piutang")
	}

	var param2 []queryutil.Param
	param2 = append(param2, queryutil.Param{
		Logic:    "AND",
		Field:    "t.date",
		Operator: ">=",
		Value:    time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	})

	param2 = append(param2, queryutil.Param{
		Logic:    "AND",
		Field:    "t.date",
		Operator: "<=",
		Value:    time.Date(month.Year(), month.Month(), dateutil.DaysIn(month.Month(), month.Year()), 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
	})

	param2 = append(param2, queryutil.Param{
		Logic:    "AND",
		Field:    "t.status",
		Operator: "NOT IN",
		Value:    status,
	})

	param2 = append(param2, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    true,
	})

	lastCreditPerMonthMap, err := uc.transactionRepo.FindLastCreditPerMonth(param2)

	customers, err := uc.customerRepo.Find(map[string]interface{}{"active": true})

	transactions := make([]transactiondomain.TransactionCreditDate, 0)
	for _, v := range customers {
		t := transactiondomain.TransactionCreditDate{
			CustomerCode: v.Code,
			CustomerName: v.Name,
			LastCredit:   lastCreditMap[v.Code] + float64(v.InitialCredit),
			Credits:      lastCreditPerMonthMap[v.Code],
		}

		transactions = append(transactions, t)
	}

	transactionCredit.PreviousMonth = dateutil.MonthToString(int(month.AddDate(0, -1, 0).Month()))
	transactionCredit.Days = dateutil.DaysIn(month.Month(), month.Year())
	transactionCredit.Transactions = transactions

	return transactionCredit, nil
}

func (uc *Usecase) FindCustomerReport(stakeholderId string, month time.Time) (*transactiondomain.LaporanCustomerSumary, error) {
	startDate := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDate := time.Date(month.Year(), month.Month(), dateutil.DaysIn(month.Month(), month.Year()), 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	laporanCustomer, totalOrder, err := uc.transactionRepo.FindCustomerReport(stakeholderId, month)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat mengambil data laporan per customer")
	}

	var param []queryutil.Param
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "k.created_time",
		Operator: ">=",
		Value:    startDate,
	})
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "k.created_time",
		Operator: "<=",
		Value:    endDate,
	})
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "customer_id",
		Operator: "=",
		Value:    stakeholderId,
	})

	trx, err := uc.kontrabonRepo.Find(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian kontrabon")
	}

	result := &transactiondomain.LaporanCustomerSumary{
		TotalOrder:  totalOrder,
		ProductData: laporanCustomer,
		Kontrabon:   trx,
		Days:        dateutil.DaysIn(month.Month(), month.Year()),
	}

	return result, nil
}
