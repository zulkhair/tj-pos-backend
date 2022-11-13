package priceusecase

import (
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	pricedomain "dromatech/pos-backend/internal/domain/price"
	customerrepo "dromatech/pos-backend/internal/repo/customer"
	pricerepo "dromatech/pos-backend/internal/repo/price"
	productrepo "dromatech/pos-backend/internal/repo/product"
	dateutil "dromatech/pos-backend/internal/util/date"
	stringutil "dromatech/pos-backend/internal/util/string"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type PriceUsecase interface {
	Find(name string) ([]*pricedomain.PriceTemplate, error)
	FindDetail(templateId string) ([]*pricedomain.PriceTemplateDetail, error)
	Create(templateName string) error
	EditPrice(templateId, productId string, price float64) error
	ApplyToCustomer(templateId string, customerId []string, userId string) error
	DeleteTemplate(templateId string) error
}

type Usecase struct {
	priceRepo    pricerepo.PriceRepo
	productRepo  productrepo.ProductRepo
	customerRepo customerrepo.CustomerRepo
}

func New(priceRepo pricerepo.PriceRepo, productRepo productrepo.ProductRepo, customerRepo customerrepo.CustomerRepo) *Usecase {
	uc := &Usecase{
		priceRepo:    priceRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}

	return uc
}

func (uc *Usecase) Find(name string) ([]*pricedomain.PriceTemplate, error) {
	param := make(map[string]interface{})

	if name != "" {
		param["pt.name"] = name
	}

	entities, err := uc.priceRepo.Find(param)
	if err != nil{
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian")
	}

	return entities, nil
}

func (uc *Usecase) FindDetail(templateId string) ([]*pricedomain.PriceTemplateDetail, error) {
	param := make(map[string]interface{})

	if templateId != "" {
		param["ptd.price_template_id"] = templateId
	}

	entities, err := uc.priceRepo.FindDetail(param)
	if err != nil{
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian")
	}

	return entities, nil
}

func (uc *Usecase) Create(templateName string) error {
	products, err := uc.priceRepo.Find(map[string]interface{}{"pt.name": templateName})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data template harga")
	}

	if products != nil || len(products) > 0 {
		return fmt.Errorf("Template dengan nama %s sudah terdaftar", templateName)
	}

	err = uc.priceRepo.Create(templateName)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data template harga")
	}

	return nil
}

func (uc *Usecase) EditPrice(templateId, productId string, price float64) error {
	priceDetail, err := uc.priceRepo.FindDetail(map[string]interface{}{"ptd.price_template_id": templateId, "ptd.product_id": productId})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}

	if priceDetail != nil {
		err = uc.priceRepo.EditPrice(templateId, productId, price)
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
		}
	} else {
		err = uc.priceRepo.AddPrice(templateId, productId, price)
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
		}
	}

	return nil
}

func (uc *Usecase) ApplyToCustomer(templateId string, customerId []string, userId string) error {
	priceDetail, err := uc.priceRepo.FindDetail(map[string]interface{}{"ptd.price_template_id": templateId})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}

	products, err := uc.productRepo.Find(map[string]interface{}{"p.active": true})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}

	priceDetailMap := make(map[string]float64)
	for _, price := range priceDetail {
		priceDetailMap[price.ProductID] = price.Price
	}

	tx := global.DBCON.Begin()
	for _, cID := range customerId {
		for _, product := range products {
			price := float64(0)
			if detail, ok := priceDetailMap[product.ID]; ok {
				price = detail
			}

			priceRequest := customerdomain.AddPriceRequest{
				ID:            stringutil.GenerateUUID(),
				Date:          time.Now().Format(dateutil.TimeFormat()),
				CustomerId:    cID,
				UnitId:        "",
				ProductID:     product.ID,
				Price:         price,
				WebUserId:     userId,
				Latest:        true,
				TransactionId: nil,
			}

			uc.customerRepo.AddSellPriceTx(priceRequest, tx)
			if tx.Error != nil {
				tx.Rollback()
				logrus.Error(err.Error())
				return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
			}
		}
	}
	tx.Commit()

	return nil
}

func (uc *Usecase) DeleteTemplate(templateId string) error{
	return uc.priceRepo.DeleteTemplate(templateId)
}