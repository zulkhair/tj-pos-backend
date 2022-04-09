package transactionusecase

import (
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type TransactoionUsecase interface {
	CreateTransaction(transaction *transactiondomain.Transaction) error
	ViewTransaction(startDate, endDate, code, stakeholderID, txType, status, productID string) ([]*transactiondomain.Transaction, error)
	UpdateStatus(transactionID, status string) error
	UpdateBuyPrice(transactionID, unitID, productID string, price float64) error
}

type Usecase struct {
	transactionRepo transactionRepo
	sequenceRepo    sequenceRepo
	supplierRepo    supplierRepo
	customerRepo    customerRepo
}

type transactionRepo interface {
	Find(params []queryutil.Param) ([]*transactiondomain.Transaction, error)
	Create(entity *transactiondomain.Transaction, tx *gorm.DB)
	Edit(product *transactiondomain.Transaction) error
	UpdateStatus(transactionID, status string) error
	UpdateBuyPrice(transactionID, unitID, productID string, price float64) error
}

type sequenceRepo interface {
	NextVal(id string) int64
	NextValTx(id string, tx *gorm.DB) int64
}

type customerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
	GetSellPrice(params []queryutil.Param) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
	DeleteSellPrice(supplierId, unitId, date string) error
}

type supplierRepo interface {
	Find(params map[string]interface{}) ([]*supplierdomain.Supplier, error)
	Create(product *supplierdomain.Supplier) error
	Edit(product *supplierdomain.Supplier) error
	GetBuyPrice(params []queryutil.Param) ([]*supplierdomain.BuyPriceResponse, error)
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
	DeleteBuyPrice(supplierId, unitId, date string) error
}

func New(transactionRepo transactionRepo, sequenceRepo sequenceRepo, supplierRepo supplierRepo, customerRepo customerRepo) *Usecase {
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
		transaction.Date = timeNow.Format("2006-01-02")
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
	transactionCode := stakeHolderCode + "/" + transaction.Date + "/"
	seq := uc.sequenceRepo.NextValTx(transactionCode, tx)
	transactionCode = transactionCode + strconv.Itoa(int(seq))

	transaction.ID = transactionID
	transaction.Code = transactionCode
	transaction.Status = transactiondomain.TRANSACTION_STATUS_PEMBUATAN

	uc.transactionRepo.Create(transaction, tx)

	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Terjadi kesalahan saat menambahkan transaksi")
	}

	tx.Commit()
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

func (uc *Usecase) UpdateStatus(transactionID, status string) error {
	err := uc.transactionRepo.UpdateStatus(transactionID, status)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}
	return nil
}

func (uc *Usecase) UpdateBuyPrice(transactionID, unitID, productID string, price float64) error {
	err := uc.transactionRepo.UpdateBuyPrice(transactionID, unitID, productID, price)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan update status transaksi")
	}
	return nil
}
