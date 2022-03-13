package customerusecase

import (
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type CustmerUsecase interface {
	Find(id, code, name string) ([]*customerdomain.Customer, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
	GetSellPrice(supplierId, unitId, date string) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
}

type Usecase struct {
	customerRepo customerRepo
}

type customerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
	GetSellPrice(supplierId, unitId, date string) ([]*customerdomain.SellPriceResponse, error)
	UpdateSellPrice(request customerdomain.SellPriceRequest) error
	DeleteSellPrice(supplierId, unitId, date string) error
}

func New(customerRepo customerRepo) *Usecase {
	uc := &Usecase{
		customerRepo: customerRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code, name string) ([]*customerdomain.Customer, error) {
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
	return uc.customerRepo.Find(param)
}

func (uc *Usecase) Create(code, name, description string) error {
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
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Active:      true,
	}

	err = uc.customerRepo.Create(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data customer")
	}

	return nil
}

func (uc *Usecase) Edit(id, code, name, description string, active bool) error {
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

	err = uc.customerRepo.Edit(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data customer")
	}

	return nil
}

func (uc *Usecase) GetSellPrice(customerId, unitId, date string) ([]*customerdomain.SellPriceResponse, error) {
	if customerId == "" {
		return nil, fmt.Errorf("Harap pilih customer terlebih dahulu")
	}
	if unitId == "" {
		return nil, fmt.Errorf("Harap pilih satuan terlebih dahulu")
	}
	if date == "" {
		return nil, fmt.Errorf("Harap pilih tanggal terlebih dahulu")
	}

	entities, err := uc.customerRepo.GetSellPrice(customerId, unitId, date)
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

	err := uc.customerRepo.DeleteSellPrice(request.CustomerId, request.UnitId, request.Date)
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
