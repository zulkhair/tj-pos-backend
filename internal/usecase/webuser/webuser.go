package webuserusecase

import (
	"crypto/sha256"
	webuserdomain "dromatech/pos-backend/internal/domain/webuser"
	"encoding/base64"
	"fmt"
)

type SessionUsecase interface {
	EditUser(userId string, name string)
	ChangePassword(userId string, password1 string, password2 string) error
}

type Usecase struct {
	webuserrepo webUserRepo
}

type webUserRepo interface {
	ReInitCache() error
	Find(id string) *webuserdomain.WebUser
	FindByUsername(username string) *webuserdomain.WebUser
	EditUser(*webuserdomain.WebUser)
	ChangePassword(userId string, newPassword string)
}

func New(webuserrepo webUserRepo) *Usecase {
	uc := &Usecase{
		webuserrepo: webuserrepo,
	}

	return uc
}

func (uc *Usecase) EditUser(userId string, name string) {
	webuser := uc.webuserrepo.Find(userId)
	webuser.Name = name

	uc.webuserrepo.EditUser(webuser)
}

func (uc *Usecase) ChangePassword(userId string, password1 string, password2 string) error {
	webuser := uc.webuserrepo.Find(userId)
	hasher := sha256.New()
	hasher.Write([]byte(password1 + webuser.PasswordSalt))
	passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if passwordHash != webuser.PasswordHash {
		return fmt.Errorf("Password lama tidak sesuai")
	}

	hasher = sha256.New()
	hasher.Write([]byte(password2 + webuser.PasswordSalt))
	passwordHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	uc.webuserrepo.ChangePassword(webuser.ID, passwordHash)
	return nil
}
