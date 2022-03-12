package productusecase

import (
	productdomain "dromatech/pos-backend/internal/domain/product"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type ProductUsecase interface {
	Find(id, code, name string) ([]*productdomain.Product, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
}

type Usecase struct {
	productRepo productRepo
}

type productRepo interface {
	Find(params map[string]interface{}) ([]*productdomain.Product, error)
	Create(product *productdomain.Product) error
	Edit(product *productdomain.Product) error
}

func New(productRepo productRepo) *Usecase {
	uc := &Usecase{
		productRepo: productRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code, name string) ([]*productdomain.Product, error) {
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
	return uc.productRepo.Find(param)
}

func (uc *Usecase) Create(code, name, description string) error {
	products, err := uc.productRepo.Find(map[string]interface{}{"code": code})
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
	}

	err = uc.productRepo.Create(product)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data produk")
	}

	return nil
}

func (uc *Usecase) Edit(id, code, name, description string, active bool) error {
	products, err := uc.productRepo.Find(map[string]interface{}{"id": id})
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
		products, err := uc.productRepo.Find(map[string]interface{}{"code": code})
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
