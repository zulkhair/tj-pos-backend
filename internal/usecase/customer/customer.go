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
}

type Usecase struct {
	customerRepo customerRepo
}

type customerRepo interface {
	Find(params map[string]interface{}) ([]*customerdomain.Customer, error)
	Create(product *customerdomain.Customer) error
	Edit(product *customerdomain.Customer) error
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
		return fmt.Errorf("Ada kesalahan saat melakukan penambahan data customer")
	}

	if entities != nil || len(entities) > 0 {
		return fmt.Errorf("Customer dengan kode %s sudah terdaftar", code)
	}

	id := strings.ReplaceAll(uuid.NewString(), "-", "")

	product := &customerdomain.Customer{
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Active:      true,
	}

	return uc.customerRepo.Create(product)
}

func (uc *Usecase) Edit(id, code, name, description string, active bool) error {
	entities, err := uc.customerRepo.Find(map[string]interface{}{"id": id})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data customer")
	}

	if len(entities) != 1 {
		logrus.Errorf("Product with id %s more than 1", id)
		return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data customer")
	}

	entity := entities[0]

	if code != entity.Code {
		products, err := uc.customerRepo.Find(map[string]interface{}{"code": code})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data customer")
		}

		if products != nil || len(products) > 0 {
			return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
		}
	}

	entity.Code = code
	entity.Name = name
	entity.Description = description
	entity.Active = active

	return uc.customerRepo.Edit(entity)
}
