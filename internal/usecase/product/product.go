package productusecase

import (
	productdomain "dromatech/pos-backend/internal/domain/product"
	productrepo "dromatech/pos-backend/internal/repo/product"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type ProductUsecase interface {
	Find(id, code, name string, active *bool) ([]*productdomain.Product, error)
	Create(code, name, description, unitId string) error
	Edit(id, code, name, description string, active bool) error
}

type Usecase struct {
	productRepo productrepo.ProductRepo
}

func New(productRepo productrepo.ProductRepo) *Usecase {
	uc := &Usecase{
		productRepo: productRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code, name string, active *bool) ([]*productdomain.Product, error) {
	param := make(map[string]interface{})
	if id != "" {
		param["p.id"] = id
	}
	if code != "" {
		param["p.code"] = code
	}
	if name != "" {
		param["p.name"] = name
	}
	if active != nil {
		param["p.active"] = *active
	}
	return uc.productRepo.Find(param)
}

func (uc *Usecase) Create(code, name, description, unitId string) error {
	products, err := uc.productRepo.Find(map[string]interface{}{"p.code": code})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data produk")
	}

	if products != nil || len(products) > 0 {
		return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
	}

	id := strings.ReplaceAll(uuid.NewString(), "-", "")

	product := &productdomain.Product{
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Active:      true,
		UnitID:      unitId,
	}

	err = uc.productRepo.Create(product)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data produk")
	}

	return nil
}

func (uc *Usecase) Edit(id, code, name, description string, active bool) error {
	products, err := uc.productRepo.Find(map[string]interface{}{"p.id": id})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data produk")
	}

	if len(products) != 1 {
		logrus.Errorf("Product with id %s more than 1", id)
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data produk")
	}

	product := products[0]

	if code != product.Code {
		products, err := uc.productRepo.Find(map[string]interface{}{"p.code": code})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data produk")
		}

		if products != nil || len(products) > 0 {
			return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
		}
	}

	product.Code = code
	product.Name = name
	product.Description = description
	product.Active = active

	err = uc.productRepo.Edit(product)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data produk")
	}

	return nil
}
