package unitusecase

import (
	unitdomain "dromatech/pos-backend/internal/domain/unit"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

type UnitUsecase interface {
	Find(id, code string, active *bool) ([]*unitdomain.Unit, error)
	Create(code, description string) error
	Edit(id, code, description string, active *bool) error
}

type Usecase struct {
	unitRepo unitRepo
}

type unitRepo interface {
	Find(params map[string]interface{}) ([]*unitdomain.Unit, error)
	Create(product *unitdomain.Unit) error
	Edit(product *unitdomain.Unit) error
}

func New(unitRepo unitRepo) *Usecase {
	uc := &Usecase{
		unitRepo: unitRepo,
	}

	return uc
}

func (uc *Usecase) Find(id, code string, active *bool) ([]*unitdomain.Unit, error) {
	param := make(map[string]interface{})
	if id != "" {
		param["id"] = id
	}
	if code != "" {
		param["code"] = code
	}
	if active != nil {
		param["active"] = active
	}
	return uc.unitRepo.Find(param)
}

func (uc *Usecase) Create(code, description string) error {
	entities, err := uc.unitRepo.Find(map[string]interface{}{"code": code})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan penambahan data satuan")
	}

	if entities != nil || len(entities) > 0 {
		return fmt.Errorf("Satuan dengan kode %s sudah terdaftar", code)
	}

	id := strings.ReplaceAll(uuid.NewString(), "-", "")

	product := &unitdomain.Unit{
		ID:          id,
		Code:        code,
		Description: description,
		Active:      true,
	}

	err = uc.unitRepo.Create(product)

	return nil
}

func (uc *Usecase) Edit(id, code, description string, active *bool) error {
	entities, err := uc.unitRepo.Find(map[string]interface{}{"id": id})
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data satuan")
	}

	if len(entities) != 1 {
		logrus.Errorf("Product with id %s more than 1", id)
		return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data satuan")
	}

	entity := entities[0]

	if code != entity.Code {
		products, err := uc.unitRepo.Find(map[string]interface{}{"code": code})
		if err != nil {
			logrus.Error(err.Error())
			return fmt.Errorf("Ada kesalahan saat melakukan pembaruan data satuan")
		}

		if products != nil || len(products) > 0 {
			return fmt.Errorf("Produk dengan kode %s sudah terdaftar", code)
		}
	}

	if code != "" {
		entity.Code = code
	}
	if description != "" {
		entity.Description = description
	}
	if active != nil {
		entity.Active = *active
	}

	err = uc.unitRepo.Edit(entity)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan pembaruan data satuan")
	}

	return nil
}
