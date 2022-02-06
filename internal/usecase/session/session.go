package session

import (
	"context"
	"net/http"
)

type SessionUsecase interface {
	AuthCheck(token string) (string, int)
	Logout(ctx context.Context, r *http.Request) error
}

type Usecase struct {
}

func New() *Usecase {
	uc := &Usecase{}
	return uc
}

func (uc *Usecase) Login(username string, password string) {

}

func (uc *Usecase) AuthCheck(token string) (string, int) {

}
