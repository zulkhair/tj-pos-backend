package sessionusecase

import (
	"crypto/sha256"
	"dromatech/pos-backend/global"
	configdomain "dromatech/pos-backend/internal/domain/config"
	permissiondomain "dromatech/pos-backend/internal/domain/permission"
	roledomain "dromatech/pos-backend/internal/domain/role"
	sessiondomain "dromatech/pos-backend/internal/domain/session"
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
	Login(username string, password string) (*sessiondomain.Session, error)
	Logout(token string)
	AuthCheck(token string, requestorPath string) (string, int)
	GetSession(token string) *sessiondomain.Session
}

type Usecase struct {
	sessionCache   sessiondomain.SessionCache
	configRepo     configRepo
	webuserrepo    webUserRepo
	permissionRepo permissionRepo
	roleRepo       roleRepo
}

type configRepo interface {
	ReInitCache() error
	GetValue(key string) string
}

type webUserRepo interface {
	ReInitCache() error
	Find(id string) *webuserdomain.WebUser
	FindByUsername(username string) *webuserdomain.WebUser
}

type permissionRepo interface {
	ReInitCache() error
	Find(id string) *permissiondomain.Permission
	FindByRoleId(roleId string) []*permissiondomain.Permission
}

type roleRepo interface {
	ReInitCache() error
	Find(id string) *roledomain.Role
	FindMenu(roleId string) ([]*sessiondomain.Menu, error)
}

func New(configRepo configRepo, webuserrepo webUserRepo, permissionRepo permissionRepo, roleRepo roleRepo) *Usecase {
	sessionCache := sessiondomain.SessionCache{
		DataMap: make(map[string]*sessiondomain.Session),
	}
	uc := &Usecase{
		sessionCache:   sessionCache,
		configRepo:     configRepo,
		webuserrepo:    webuserrepo,
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
	}

	go uc.removeExpiredSession()

	return uc
}

func (uc *Usecase) removeExpiredSession() {
	for {
		now := time.Now()
		count := 0
		for key, value := range uc.sessionCache.DataMap {
			if now.After(value.ExpiredTime) {
				count++
				uc.sessionCache.Lock()
				delete(uc.sessionCache.DataMap, key)
				uc.sessionCache.Unlock()
			}
		}

		logrus.Printf("session sweep finished, %d sessions deleted", count)
		time.Sleep(time.Second * 60)
	}
}

func (uc *Usecase) Login(username string, password string) (*sessiondomain.Session, error) {
	webuser := uc.webuserrepo.FindByUsername(username)
	if webuser == nil {
		return nil, fmt.Errorf("User atau Password yang dimasukan salah")
	}

	hasher := sha256.New()
	hasher.Write([]byte(password + webuser.PasswordSalt))
	passwordHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if passwordHash != webuser.PasswordHash {
		return nil, fmt.Errorf("User atau Password yang dimasukan salah")
	}

	role := uc.roleRepo.Find(webuser.RoleId)
	menus, err := uc.roleRepo.FindMenu(webuser.RoleId)
	if err != nil {
		logrus.Error(err.Error())
	}

	token := strings.ReplaceAll(uuid.NewString(), "-", "")
	minute, _ := strconv.Atoi(uc.configRepo.GetValue(configdomain.SESSION_TIMEOUT_MINUTE))
	expiredTime := time.Now().Add(time.Minute * time.Duration(minute))

	session := &sessiondomain.Session{
		Token:       token,
		ExpiredTime: expiredTime,
		UserID:      webuser.ID,
		UserName:    webuser.Username,
		Name:        webuser.Name,
		RoleName:    role.Name,
		Menu:        menus,
	}

	uc.sessionCache.Lock()
	uc.sessionCache.DataMap[token] = session
	uc.sessionCache.Unlock()

	return session, nil
}

func (uc *Usecase) AuthCheck(token string, requestorPath string) (string, int, *sessiondomain.Session) {
	if token == "" {
		return uc.configRepo.GetValue(configdomain.LOGIN_URL), 301, nil
	}

	session, ok := uc.sessionCache.DataMap[token]
	if !ok {
		return uc.configRepo.GetValue(configdomain.UNAUTHORIZED_URL), 401, nil
	}

	user := uc.webuserrepo.Find(session.UserID)
	if user == nil {
		return uc.configRepo.GetValue(configdomain.UNAUTHORIZED_URL), 401, nil
	}

	permissions := uc.permissionRepo.FindByRoleId(user.RoleId)
	if permissions == nil {
		return uc.configRepo.GetValue(configdomain.FORBIDDEN_URL), 403, nil
	}

	authorized := false
	for _, permission := range permissions {
		paths := strings.Split(permission.Paths, ";")
		for _, path := range paths {
			if requestorPath == path {
				authorized = true
				break
			}
		}
		if authorized {
			break
		}
	}

	if !authorized {
		return uc.configRepo.GetValue(configdomain.FORBIDDEN_URL), 403, nil
	}

	session.ExpiredTime = time.Now().Add(time.Minute * time.Duration(global.SESSION_TIMEOUT_MINUTE))
	return "", 200, session
}

func (uc *Usecase) Logout(token string) {
	uc.sessionCache.Lock()
	delete(uc.sessionCache.DataMap, token)
	uc.sessionCache.Unlock()
}

func (uc *Usecase) GetSession(token string) *sessiondomain.Session {
	if session, ok := uc.sessionCache.DataMap[token]; ok {
		session.ExpiredTime = time.Now().Add(time.Minute * time.Duration(global.SESSION_TIMEOUT_MINUTE))
		return session
	}
	return nil
}
