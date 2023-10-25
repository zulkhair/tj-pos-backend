package customerusecase

import (
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type CustmerUsecase interface {
	Find(id, code, name string, active *bool) ([]*customerdomain.Customer, error)
	Create(code, name, description string, initialBalance float64) error
	Edit(id, code, name, description string, active bool, initialBalance float64) error
	GetSellPrice(customerId, unitId, date, productId string) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
	AddSellPrice(entity customerdomain.AddPriceRequest, userId string) error
	FindSellPrice(customerId, unitId, productId string, latest *bool) ([]*customerdomain.PriceResponse, error)
}

type Usecase struct {
	customerRepo customerRepo
}

type customerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
	GetSellPrice(params []queryutil.Param) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
	DeleteSellPrice(supplierId, date string) error
	AddSellPrice(entity customerdomain.AddPriceRequest) error
	FindSellPrice(params []queryutil.Param) ([]*customerdomain.PriceResponse, error)
}

func New(customerRepo customerRepo) *Usecase {
	uc := &Usecase{
		customerRepo: customerRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code, name string, active *bool) ([]*customerdomain.Customer, error) {
	param := make(map[string]interface{})
	if id != "" {
		param["id"] = id
	}
	if code != "" {
		param["code"] = code
	}
	if name != "" {
		param["name"] = name
	}
	if active != nil {
		param["active"] = *active
	}
	return uc.customerRepo.Find(param)
}

func (uc *Usecase) Create(code, name, description string, initialCredit float64) error {
	entities, err := uc.customerRepo.Find(map[string]interface{}{"code": code})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data customer")
	}

	if entities != nil || len(entities) > 0 {
		return fmt.Errorf("Customer dengan kode %s sudah terdaftar", code)
	}

	id := strings.ReplaceAll(uuid.NewString(), "-", "")

	entity := &customerdomain.Customer{
		ID:            id,
		Code:          code,
		Name:          name,
		Description:   description,
		Active:        true,
		InitialCredit: initialCredit,
	}

	err = uc.customerRepo.Create(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data customer")
	}

	return nil
}

func (uc *Usecase) Edit(id, code, name, description string, active bool, initialCredit float64) error {
	entities, err := uc.customerRepo.Find(map[string]interface{}{"id": id})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data customer")
	}

	if len(entities) != 1 {
		logrus.Errorf("Product with id %s more than 1", id)
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data customer")
	}

	entity := entities[0]

	if code != entity.Code {
		products, err := uc.customerRepo.Find(map[string]interface{}{"code": code})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data customer")
		}

		if products != nil || len(products) > 0 {
			return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
		}
	}

	entity.Code = code
	entity.Name = name
	entity.Description = description
	entity.Active = active
	entity.InitialCredit = initialCredit

	err = uc.customerRepo.Edit(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data customer")
	}

	return nil
}

func (uc *Usecase) GetSellPrice(customerId, unitId, date, productId string) ([]*customerdomain.SellPriceResponse, error) {
	var param []queryutil.Param
	if customerId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "bp.customer_id",
			Operator: "=",
			Value:    customerId,
		})
	}
	if unitId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "bp.unit_id",
			Operator: "=",
			Value:    unitId,
		})
	}
	if date != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "bp.date",
			Operator: "=",
			Value:    date,
		})
	}
	if productId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "p.id",
			Operator: "=",
			Value:    productId,
		})
	}

	entities, err := uc.customerRepo.GetSellPrice(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian data harga")
	}

	return entities, nil
}

func (uc *Usecase) UpdateSellPrice(request customerdomain.SellPriceRequest) error {
	if request.CustomerId == "" {
		return fmt.Errorf("Harap pilih customer")
	}
	if request.UnitId == "" {
		return fmt.Errorf("Harap pilih satuan")
	}

	err := uc.customerRepo.DeleteSellPrice(request.CustomerId, request.Date)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data harga")
	}

	err = uc.customerRepo.UpdateSellPrice(request)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data harga")
	}
	return nil
}

func (uc *Usecase) AddSellPrice(entity customerdomain.AddPriceRequest, userId string) error {
	entity.ID = strings.ReplaceAll(uuid.NewString(), "-", "")
	entity.Date = time.Now()
	entity.WebUserId = userId
	entity.Latest = true

	err := uc.customerRepo.AddSellPrice(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat menambahkan data harga")
	}

	return nil
}

func (uc *Usecase) FindSellPrice(customerId, unitId, productId string, latest *bool) ([]*customerdomain.PriceResponse, error) {
	var param []queryutil.Param
	if customerId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "s.customer_id",
			Operator: "=",
			Value:    customerId,
		})
	}
	if unitId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "p.unit_id",
			Operator: "=",
			Value:    unitId,
		})
	}
	if productId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "s.product_id",
			Operator: "=",
			Value:    productId,
		})
	}
	if latest != nil {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "s.latest",
			Operator: "=",
			Value:    strconv.FormatBool(*latest),
		})
	}
	entities, err := uc.customerRepo.FindSellPrice(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarion data harga")
	}
	return entities, nil
}
