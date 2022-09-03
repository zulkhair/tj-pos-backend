package transactionusecase

import (
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	customerrepo "dromatech/pos-backend/internal/repo/customer"
	sequencerepo "dromatech/pos-backend/internal/repo/sequence"
	supplierrepo "dromatech/pos-backend/internal/repo/supplier"
	transactionrepo "dromatech/pos-backend/internal/repo/transaction"
	dateutil "dromatech/pos-backend/internal/util/date"
	queryutil "dromatech/pos-backend/internal/util/query"
	stringutil "dromatech/pos-backend/internal/util/string"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type TransactoionUsecase interface {
	CreateTransaction(transaction *transactiondomain.Transaction) error
	ViewTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.Transaction, error)
	UpdateStatus(transactionID string) error
	UpdateBuyPrice(transactionID, productID string, buyPrice, sellPrice float64, quantity, buyQuantity int64) error
	ViewSellTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.TransactionStatus, error)
	CancelTrx(transactionID string) error
}

type Usecase struct {
	transactionRepo transactionrepo.TransactionRepo
	sequenceRepo    sequencerepo.SequenceRepo
	supplierRepo    supplierrepo.SupplierRepo
	customerRepo    customerrepo.CustomerRepo
}

func New(transactionRepo transactionrepo.TransactionRepo, sequenceRepo sequencerepo.SequenceRepo, supplierRepo supplierrepo.SupplierRepo, customerRepo customerrepo.CustomerRepo) *Usecase {
	uc := &Usecase{
		transactionRepo: transactionRepo,
		sequenceRepo:    sequenceRepo,
		supplierRepo:    supplierRepo,
		customerRepo:    customerRepo,
	}

	return uc
}

func (uc *Usecase) CreateTransaction(transaction *transactiondomain.Transaction) error {
	tx := global.DBCON.Begin()

	timeNow := time.Now()
	if transaction.Date == "" {
		transaction.Date = timeNow.Format(dateutil.DateFormat())
	}
	stakeHolderCode := ""
	if transaction.TransactionType == transactiondomain.TRANSACTION_TYPE_BUY {
		supplier, err := uc.supplierRepo.Find(map[string]interface{}{"id": transaction.StakeholderID})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Supplier dengan kode '%s' tidak ditemukan", transaction.StakeholderID)
		}
		stakeHolderCode = supplier[0].Code
	} else {
		customer, err := uc.customerRepo.Find(map[string]interface{}{"id": transaction.StakeholderID})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Customer dengan kode '%s' tidak ditemukan", transaction.StakeholderID)
		}
		stakeHolderCode = customer[0].Code
	}
	transactionID := strings.ReplaceAll(uuid.NewString(), "-", "")
	dateCode, err := time.Parse(dateutil.DateFormat(), transaction.Date)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan transaksi")
	}

	transactionCode := stakeHolderCode + "/" + stringutil.ToRoman(int(dateCode.Month())) + "/" + dateCode.Format("2006")
	seq := uc.sequenceRepo.NextValTx(transactionCode, tx)
	transactionCode = strconv.Itoa(int(seq)) + "/" + transactionCode

	transaction.ID = transactionID
	transaction.Code = transactionCode
	transaction.Status = transactiondomain.TRANSACTION_PEMBUATAN
	transaction.CreatedTime = timeNow.Format(dateutil.TimeFormat())

	uc.transactionRepo.Create(transaction, tx)

	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
	}

	for _, detail := range transaction.TransactionDetail {
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
			return fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
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
					Date:          timeNow.Format(dateutil.TimeFormat()),
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
				Date:          timeNow.Format(dateutil.TimeFormat()),
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
			return fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
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
					Date:          timeNow.Format(dateutil.TimeFormat()),
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
				Date:          timeNow.Format(dateutil.TimeFormat()),
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
		return fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
	}
	return nil
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

func (uc *Usecase) ViewSellTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.TransactionStatus, error) {
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
	if err != nil{
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	if transactions == nil || len(transactions) == 0{
		logrus.Error("Transactions is nil or zero")
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	status := transactions[0].Status
	if status == transactiondomain.TRANSACTION_PEMBUATAN {
		status = transactiondomain.TRANSACTION_KONTRABON
	}else if transactions[0].Status == transactiondomain.TRANSACTION_KONTRABON {
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
	if err != nil{
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}

	if transactions == nil || len(transactions) == 0{
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
