package kontrabonusecase

import (
	"dromatech/pos-backend/global"
	kontrabondomain "dromatech/pos-backend/internal/domain/kontrabon"
	transactiondomain "dromatech/pos-backend/internal/domain/transaction"
	customerrepo "dromatech/pos-backend/internal/repo/customer"
	kontrabonrepo "dromatech/pos-backend/internal/repo/kontrabon"
	sequencerepo "dromatech/pos-backend/internal/repo/sequence"
	queryutil "dromatech/pos-backend/internal/util/query"
	stringutil "dromatech/pos-backend/internal/util/string"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type KontrabonUsecase interface {
	Find(code, startDate, endDate, customerId string) ([]*kontrabondomain.KontrabonResponse, error)
	FindTransaction(kontrabonId string) ([]*transactiondomain.TransactionStatus, error)
	Create(customerId string, transactionIds []string) error
	Update(kontrabonId string, transactionIds []string, status string) error
	UpdateLunas(kontrabonId string, paymentTime time.Time, totalPayment float64, description, paymentDate string) error
}

type Usecase struct {
	kontrabonRepo kontrabonrepo.KontrabonRepo
	sequenceRepo  sequencerepo.SequenceRepo
	customerRepo  customerrepo.CustomerRepo
}

func New(kontrabonRepo kontrabonrepo.KontrabonRepo, sequenceRepo sequencerepo.SequenceRepo, customerRepo customerrepo.CustomerRepo) *Usecase {
	uc := &Usecase{
		kontrabonRepo: kontrabonRepo,
		sequenceRepo:  sequenceRepo,
		customerRepo:  customerRepo,
	}

	return uc
}

func (uc *Usecase) Find(code, startDate, endDate, customerId string) ([]*kontrabondomain.KontrabonResponse, error) {
	var param []queryutil.Param
	if code != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "k.code",
			Operator: "=",
			Value:    code,
		})
	}
	if startDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "k.created_time",
			Operator: ">=",
			Value:    startDate,
		})
	}
	if endDate != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "k.created_time",
			Operator: "<=",
			Value:    endDate,
		})
	}
	if customerId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "k.customer_id",
			Operator: "=",
			Value:    customerId,
		})
	}

	return uc.kontrabonRepo.Find(param)
}

func (uc *Usecase) FindTransaction(kontrabonId string) ([]*transactiondomain.TransactionStatus, error) {
	var param []queryutil.Param
	if kontrabonId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "kt.kontrabon_id",
			Operator: "=",
			Value:    kontrabonId,
		})
	}

	trx, err := uc.kontrabonRepo.FindTransaction(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian transaksi")
	}

	return trx, nil
}

func (uc *Usecase) Create(customerId string, transactionIds []string) error {
	customer, err := uc.customerRepo.Find(map[string]interface{}{"id": customerId})
	if err != nil || len(customer) == 0 {
		return fmt.Errorf("Customer dengan ID %s tidak ditemukan", customerId)
	}

	createdTime := time.Now().UTC()
	tx := global.DBCON.Begin()
	code := uc.sequenceRepo.NextValTx(customer[0].Code, tx)
	if tx.Error != nil {
		logrus.Error(tx.Error.Error())
		return fmt.Errorf("Terjadi kesalahan saat pembuatan kontrabon")
	}

	kontrabon := kontrabondomain.Kontrabon{
		ID:          stringutil.GenerateUUID(),
		Code:        "KTBN/" + strconv.Itoa(int(code)) + "/" + customer[0].Code + "/" + stringutil.ToRoman(int(createdTime.Month())) + "/" + createdTime.Format("2006"),
		CreatedTime: createdTime,
		Status:      kontrabondomain.STATUS_CREATED,
		CustomerID:  customerId,
	}

	uc.kontrabonRepo.CreateTx(kontrabon, transactionIds, tx)
	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(tx.Error.Error())
		return fmt.Errorf("Terjadi kesalahan saat pembuatan kontrabon")
	}

	tx.Commit()
	return nil
}

func (uc *Usecase) Update(kontrabonId string, transactionIds []string, status string) error {
	err := uc.kontrabonRepo.Update(kontrabonId, transactionIds, status)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data kontrabon")
	}

	return nil
}

func (uc *Usecase) UpdateLunas(kontrabonId string, paymentTime time.Time, totalPayment float64, description, paymentDate string) error {
	err := uc.kontrabonRepo.UpdateLunas(kontrabonId, paymentTime, totalPayment, description, paymentDate)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan status")
	}
	return nil
}
