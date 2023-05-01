package supplierusecase

import (
	supplierdomain "dromatech/pos-backend/internal/domain/supplier"
	queryutil "dromatech/pos-backend/internal/util/query"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type SupplierUsecase interface {
	Find(id, code, name string, active *bool) ([]*supplierdomain.Supplier, error)
	Create(code, name, description string) error
	Edit(id, code, name, description string, active bool) error
	GetBuyPrice(supplierId, unitId, date, productId string) ([]*supplierdomain.BuyPriceResponse, error)
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
}

type Usecase struct {
	supplierRepo supplierRepo
}

type supplierRepo interface {
	Find(params map[string]interface{}) ([]*supplierdomain.Supplier, error)
	Create(product *supplierdomain.Supplier) error
	Edit(product *supplierdomain.Supplier) error
	GetBuyPrice(params []queryutil.Param) ([]*supplierdomain.BuyPriceResponse, error)
	UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error
	DeleteBuyPrice(supplierId, unitId, date string) error
	AddBuyPrice(entity supplierdomain.AddPriceRequest) error
	FindBuyPrice(params []queryutil.Param) ([]*supplierdomain.PriceResponse, error)
}

func New(supplierRepo supplierRepo) *Usecase {
	uc := &Usecase{
		supplierRepo: supplierRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code, name string, active *bool) ([]*supplierdomain.Supplier, error) {
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
		param["active"] = active
	}
	return uc.supplierRepo.Find(param)
}

func (uc *Usecase) Create(code, name, description string) error {
	entities, err := uc.supplierRepo.Find(map[string]interface{}{"code": code})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data supplier")
	}

	if entities != nil || len(entities) > 0 {
		return fmt.Errorf("Supplier dengan kode %s sudah terdaftar", code)
	}

	id := strings.ReplaceAll(uuid.NewString(), "-", "")

	entity := &supplierdomain.Supplier{
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Active:      true,
	}

	err = uc.supplierRepo.Create(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data supplier")
	}

	return nil
}

func (uc *Usecase) Edit(id, code, name, description string, active bool) error {
	entities, err := uc.supplierRepo.Find(map[string]interface{}{"id": id})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data supplier")
	}

	if len(entities) != 1 {
		logrus.Errorf("Product with id %s more than 1", id)
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data supplier")
	}

	entity := entities[0]

	if code != entity.Code {
		products, err := uc.supplierRepo.Find(map[string]interface{}{"code": code})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data supplier")
		}

		if products != nil || len(products) > 0 {
			return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
		}
	}

	entity.Code = code
	entity.Name = name
	entity.Description = description
	entity.Active = active

	err = uc.supplierRepo.Edit(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data supplier")
	}

	return nil
}

func (uc *Usecase) GetBuyPrice(supplierId, unitId, date, productId string) ([]*supplierdomain.BuyPriceResponse, error) {
	var param []queryutil.Param
	if supplierId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "bp.supplier_id",
			Operator: "=",
			Value:    supplierId,
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

	entities, err := uc.supplierRepo.GetBuyPrice(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarian data harga")
	}

	return entities, nil
}

func (uc *Usecase) UpdateBuyPrice(request supplierdomain.BuyPriceRequest) error {
	if request.SupplierId == "" {
		return fmt.Errorf("Harap pilih supplier")
	}
	if request.UnitId == "" {
		return fmt.Errorf("Harap pilih satuan")
	}

	err := uc.supplierRepo.DeleteBuyPrice(request.SupplierId, request.UnitId, request.Date)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data harga")
	}

	err = uc.supplierRepo.UpdateBuyPrice(request)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data harga")
	}
	return nil
}

func (uc *Usecase) AddBuyPrice(entity supplierdomain.AddPriceRequest, userId string) error {
	entity.ID = strings.ReplaceAll(uuid.NewString(), "-", "")
	entity.Date = time.Now().UTC()
	entity.WebUserId = userId
	entity.Latest = true

	err := uc.supplierRepo.AddBuyPrice(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat menambahkan data harga")
	}

	return nil
}

func (uc *Usecase) FindBuyPrice(supplierId, unitId, productId string, latest *bool) ([]*supplierdomain.PriceResponse, error) {
	var param []queryutil.Param
	if supplierId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "s.supplier_id",
			Operator: "=",
			Value:    supplierId,
		})
	}
	if unitId != "" {
		param = append(param, queryutil.Param{
			Logic:    "AND",
			Field:    "s.unit_id",
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
	entities, err := uc.supplierRepo.FindBuyPrice(param)
	if err != nil {
		logrus.Error(err.Error())
		return nil, fmt.Errorf("Terjadi kesalahan saat melakukan pencarion data harga")
	}
	return entities, nil
}
