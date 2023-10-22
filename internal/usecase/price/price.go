package priceusecase

import (
	"dromatech/pos-backend/global"
	customerdomain "dromatech/pos-backend/internal/domain/customer"
	pricedomain "dromatech/pos-backend/internal/domain/price"
	customerrepo "dromatech/pos-backend/internal/repo/customer"
	pricerepo "dromatech/pos-backend/internal/repo/price"
	productrepo "dromatech/pos-backend/internal/repo/product"
	transactionrepo "dromatech/pos-backend/internal/repo/transaction"
	queryutil "dromatech/pos-backend/internal/util/query"
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
	CopyTemplate(templateId, templateName string) error
	Download(request pricedomain.Download) error
	FindBuyTemplate(name string) ([]*pricedomain.PriceTemplate, error)
	FindBuyDetail(templateId string) ([]*pricedomain.PriceTemplateDetail, error)
	CreateBuyTemplate(templateName string) error
	EditBuyPrice(templateId, productId string, price float64) error
	ApplyToTrx(templateId string, date string, userId string) error
	DeleteBuyTemplate(templateId string) error
	DownloadBuy(request pricedomain.Download) error
	CopyBuyTemplate(templateId, templateName string) error
}

type Usecase struct {
	priceRepo       pricerepo.PriceRepo
	productRepo     productrepo.ProductRepo
	customerRepo    customerrepo.CustomerRepo
	transactionRepo transactionrepo.TransactionRepo
}

func New(priceRepo pricerepo.PriceRepo, productRepo productrepo.ProductRepo, customerRepo customerrepo.CustomerRepo, transactionRepo transactionrepo.TransactionRepo) *Usecase {
	uc := &Usecase{
		priceRepo:       priceRepo,
		productRepo:     productRepo,
		customerRepo:    customerRepo,
		transactionRepo: transactionRepo,
	}

	return uc
}

func (uc *Usecase) Find(name string) ([]*pricedomain.PriceTemplate, error) {
	param := make(map[string]interface{})

	if name != "" {
		param["pt.name"] = name
	}

	entities, err := uc.priceRepo.Find(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian")
	}

	return entities, nil
}

func (uc *Usecase) FindBuyTemplate(name string) ([]*pricedomain.PriceTemplate, error) {
	param := make(map[string]interface{})

	if name != "" {
		param["pt.name"] = name
	}

	entities, err := uc.priceRepo.FindBuyTemplate(param)
	if err != nil {
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
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian")
	}

	return entities, nil
}

func (uc *Usecase) FindBuyDetail(templateId string) ([]*pricedomain.PriceTemplateDetail, error) {
	param := make(map[string]interface{})

	if templateId != "" {
		param["ptd.buy_price_template_id"] = templateId
	}

	entities, err := uc.priceRepo.FindBuyDetail(param)
	if err != nil {
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

	_, err = uc.priceRepo.Create(templateName)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data template harga")
	}

	return nil
}

func (uc *Usecase) CreateBuyTemplate(templateName string) error {
	products, err := uc.priceRepo.FindBuyTemplate(map[string]interface{}{"pt.name": templateName})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data template harga")
	}

	if products != nil || len(products) > 0 {
		return fmt.Errorf("Template dengan nama %s sudah terdaftar", templateName)
	}

	_, err = uc.priceRepo.CreateBuyTemplate(templateName)
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

func (uc *Usecase) EditBuyPrice(templateId, productId string, price float64) error {
	priceDetail, err := uc.priceRepo.FindBuyDetail(map[string]interface{}{"ptd.buy_price_template_id": templateId, "ptd.product_id": productId})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}

	if priceDetail != nil {
		err = uc.priceRepo.EditBuyPrice(templateId, productId, price)
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
		}
	} else {
		err = uc.priceRepo.AddBuyPrice(templateId, productId, price)
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
		}
	}

	return nil
}

func (uc *Usecase) ApplyToCustomer(templateId string, customerId []string, userId string) error {
	date := time.Now().UTC()
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
	appliedCustomer := ""
	for _, cID := range customerId {
		appliedCustomer = appliedCustomer + cID + ";"
		for _, product := range products {
			price := float64(0)
			if detail, ok := priceDetailMap[product.ID]; ok {
				price = detail
			}

			priceRequest := customerdomain.AddPriceRequest{
				ID:            stringutil.GenerateUUID(),
				Date:          date,
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
	uc.priceRepo.UpdateTemplate(templateId, appliedCustomer, tx)
	if tx.Error != nil {
		tx.Rollback()
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}
	tx.Commit()

	return nil
}

func (uc *Usecase) ApplyToTrx(templateId string, date string, userId string) error {
	now := time.Now().UTC()
	priceDetail, err := uc.priceRepo.FindBuyDetail(map[string]interface{}{"ptd.buy_price_template_id": templateId})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}

	productIds := make([]string, 0)
	priceDetailMap := make(map[string]float64)
	for _, price := range priceDetail {
		priceDetailMap[price.ProductID] = price.Price
		productIds = append(productIds, price.ProductID)
	}

	var param []queryutil.Param
	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "td.latest",
		Operator: "=",
		Value:    true,
	})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "t.date",
		Operator: "=",
		Value:    date,
	})

	param = append(param, queryutil.Param{
		Logic:    "AND",
		Field:    "td.product_id",
		Operator: "IN",
		Value:    productIds,
	})

	transactions, err := uc.transactionRepo.FindDetails(param)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
	}
	tx := global.DBCON.Begin()
	for _, transaction := range transactions {
		price := float64(0)
		if detail, ok := priceDetailMap[transaction.ProductID]; ok {
			price = detail
		}

		if price <= 0 {
			continue
		}

		err = uc.transactionRepo.UpdateHargaBeliTx(transaction.ID, price, userId, tx)
		if err != nil {
			tx.Rollback()
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
		}

		uc.priceRepo.UpdateBuyTemplate(templateId, userId, transaction.TransactionID, now, tx)
		if tx.Error != nil {
			tx.Rollback()
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data harga")
		}

	}

	tx.Commit()

	return nil
}

func (uc *Usecase) DeleteTemplate(templateId string) error {
	return uc.priceRepo.DeleteTemplate(templateId)
}

func (uc *Usecase) DeleteBuyTemplate(templateId string) error {
	return uc.priceRepo.DeleteBuyTemplate(templateId)
}

func (uc *Usecase) Download(request pricedomain.Download) error {
	return uc.priceRepo.UpdateChecked(request)
}

func (uc *Usecase) DownloadBuy(request pricedomain.Download) error {
	return uc.priceRepo.UpdateBuyChecked(request)
}

func (uc *Usecase) CopyTemplate(templateId, templateName string) error {
	products, err := uc.priceRepo.Find(map[string]interface{}{"pt.name": templateName})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
	}

	if products != nil || len(products) > 0 {
		return fmt.Errorf("Template dengan nama %s sudah terdaftar", templateName)
	}

	ID, err := uc.priceRepo.Create(templateName)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
	}

	priceDetail, err := uc.priceRepo.FindDetail(map[string]interface{}{"ptd.price_template_id": templateId})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
	}

	for _, price := range priceDetail {
		err = uc.priceRepo.AddPrice(ID, price.ProductID, price.Price)
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
		}
	}

	return nil
}

func (uc *Usecase) CopyBuyTemplate(templateId, templateName string) error {
	products, err := uc.priceRepo.FindBuyTemplate(map[string]interface{}{"pt.name": templateName})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
	}

	if products != nil || len(products) > 0 {
		return fmt.Errorf("Template dengan nama %s sudah terdaftar", templateName)
	}

	ID, err := uc.priceRepo.CreateBuyTemplate(templateName)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
	}

	priceDetail, err := uc.priceRepo.FindBuyDetail(map[string]interface{}{"ptd.buy_price_template_id": templateId})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
	}

	for _, price := range priceDetail {
		err = uc.priceRepo.AddBuyPrice(ID, price.ProductID, price.Price)
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan duplikasi data template harga")
		}
	}

	return nil
}
