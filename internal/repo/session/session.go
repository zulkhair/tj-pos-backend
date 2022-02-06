package session

type SessionResource interface {
}

type Repo struct {
}

func New() (*Repo, error) {
	return &Repo{}, nil
}

func (r *Repo) GetMenu() {

}
