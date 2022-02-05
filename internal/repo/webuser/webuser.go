package userrepo

import webuserdomain "dromatech/pos-backend/internal/domain/webuser"

type UserRepo interface {
}

type Repo struct {
	WebUserCahce webuserdomain.WebUserCache
}

func New() (*Repo, error) {
	repo := &Repo{}

	return repo, nil
}

func reInitCache(repo *Repo) {

}
