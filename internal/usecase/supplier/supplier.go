package supplierusecase

import (
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type SupplierUsecase interface {
	Find(id, code, name string) ([]*supplierdomain.Supplier, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
}

type Usecase struct {
	supplierRepo supplierRepo
}

type supplierRepo interface {
	Find(params map[string]interface{}) ([]*supplierdomain.Supplier, error)
	Create(product *supplierdomain.Supplier) error
	Edit(product *supplierdomain.Supplier) error
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
}

func New(supplierRepo supplierRepo) *Usecase {
	uc := &Usecase{
		supplierRepo: supplierRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code, name string) ([]*supplierdomain.Supplier, error) {
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
	return uc.supplierRepo.Find(param)
}

func (uc *Usecase) Create(code, name, description string) error {
	entity, err := uc.supplierRepo.Find(map[string]interface{}{"code": code})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Ada kesalahan saat melakukan penambahan data supplier")
	}

	if entity != nil || len(entity) > 0 {
		return fmt.Errorf("Supplier dengan kode %s sudah terdaftar", code)
	}

	id := strings.ReplaceAll(uuid.NewString(), "-", "")

	product := &supplierdomain.Supplier{
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Active:      true,
	}

	return uc.supplierRepo.Create(product)
}

func (uc *Usecase) Edit(id, code, name, description string, active bool) error {
	entities, err := uc.supplierRepo.Find(map[string]interface{}{"id": id})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data supplier")
	}

	if len(entities) != 1 {
		logrus.Errorf("Product with id %s more than 1", id)
		return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data supplier")
	}

	entity := entities[0]

	if code != entity.Code {
		products, err := uc.supplierRepo.Find(map[string]interface{}{"code": code})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data supplier")
		}

		if products != nil || len(products) > 0 {
			return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
		}
	}

	entity.Code = code
	entity.Name = name
	entity.Description = description
	entity.Active = active

	return uc.supplierRepo.Edit(entity)
}

func (uc *Usecase) UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error {
	if request.SupplierId == "" {
		return fmt.Errorf("Harap pilih supplier")
	}
	if request.UnitId == "" {
		return fmt.Errorf("Harap pilih satuan")
	}

	err := uc.UpdateBuyPrice(request)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data harga")
	}
	return nil
}
