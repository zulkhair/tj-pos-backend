package session

import (
	"context"
	"net/http"
)

type SessionUsecase interface {
	AuthCheck(token string) (string, int)
	Logout(ctx context.Context, r *http.Request) error
}
