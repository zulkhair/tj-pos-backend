package webuserusecase

import (
	"crypto/sha256"
	webuserdomain "dromatech/pos-backend/internal/domain/webuser"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type SessionUsecase interface {
	EditUser(userId, name, username, role, status string) error
	ChangePassword(userId, password1, password2 string) error
	RegisterUser(creatorId, name, username, password, roleId string) error
	FindAllUser() ([]*webuserdomain.WebUser, error)
	ForceChangePassword(userId, password string) error
	ChangeStatus(userId string, status bool)
}

type Usecase struct {
	webuserrepo webUserRepo
}

type webUserRepo interface {
	FindAll() ([]*webuserdomain.WebUser, error)
	Find(id string) *webuserdomain.WebUser
	FindByUsername(username string) *webuserdomain.WebUser
	EditUser(webUser *webuserdomain.WebUser) error
	ChangePassword(userId string, newPassword string)
	RegisterUser(webUser *webuserdomain.WebUser)
	ChangeStatus(userId string, active bool)
}

func New(webuserrepo webUserRepo) *Usecase {
	uc := &Usecase{
		webuserrepo: webuserrepo,
	}

	return uc
}

func (uc *Usecase) EditUser(userId, name, username, role, status string) error {
	webuser := uc.webuserrepo.Find(userId)
	if name != "" {
		webuser.Name = name
	}
	if username != "" {
		u := uc.webuserrepo.FindByUsername(username)
		if u != nil{
			return fmt.Errorf("User dengan username '&s' sudah ada", username)
		}
		webuser.Username = username
	}
	if role != "" {
		webuser.RoleId = role
	}
	if status != "" {
		ac, err := strconv.ParseBool(status)
		if err != nil {
			webuser.Active = ac
		}
	}

	err := uc.webuserrepo.EditUser(webuser)
	if err != nil {
		logrus.Error(err.Error())
		return fmt.Errorf("Terjadi kesalahan saat melakukan perubahan data user")
	}
	return nil
}

func (uc *Usecase) ChangePassword(userId, password1, password2 string) error {
	webuser := uc.webuserrepo.Find(userId)
	hasher := sha256.New()
	hasher.Write([]byte(password1 + webuser.PasswordSalt))
	passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if passwordHash != webuser.PasswordHash {
		return fmt.Errorf("Kata sandi lama tidak sesuai")
	}

	hasher = sha256.New()
	hasher.Write([]byte(password2 + webuser.PasswordSalt))
	passwordHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	uc.webuserrepo.ChangePassword(webuser.ID, passwordHash)
	return nil
}

func (uc *Usecase) ForceChangePassword(userId, password string) error {
	webuser := uc.webuserrepo.Find(userId)
	hasher := sha256.New()
	hasher.Write([]byte(password + webuser.PasswordSalt))
	passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	uc.webuserrepo.ChangePassword(webuser.ID, passwordHash)
	return nil
}

func (uc *Usecase) RegisterUser(creatorId, name, username, password, roleId string) error {
	webuser := uc.webuserrepo.FindByUsername(username)
	if webuser != nil {
		return fmt.Errorf("Pengguna dengan username '%s' sudah ada", username)
	}

	passwordSalt := strings.ReplaceAll(uuid.NewString(), "-", "")
	hasher := sha256.New()
	hasher.Write([]byte(password + passwordSalt))
	passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	webuser = &webuserdomain.WebUser{
		ID:                    strings.ReplaceAll(uuid.NewString(), "-", ""),
		Name:                  name,
		Username:              username,
		PasswordHash:          passwordHash,
		PasswordSalt:          passwordSalt,
		Email:                 "-",
		RoleId:                roleId,
		Active:                true,
		RegistrationTimestamp: time.Now(),
		CreatedBy:             creatorId,
	}

	uc.webuserrepo.RegisterUser(webuser)
	return nil
}

func (uc *Usecase) FindAllUser() ([]*webuserdomain.WebUser, error) {
	return uc.webuserrepo.FindAll()
}

func (uc *Usecase) ChangeStatus(userId string, status bool) {
	uc.webuserrepo.ChangeStatus(userId, status)
}
