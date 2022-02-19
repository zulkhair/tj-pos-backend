package configrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	configdomain "dromatech/pos-backend/internal/domain/config"
	"github.com/sirupsen/logrus"
)

type ConfigRepo interface {
	ReInitCache() error
	GetValue(key string) string
}

type Repo struct {
	ConfigCache configdomain.ConfigCache
}

func New() (*Repo, error) {
	repo := &Repo{}
	err := repo.ReInitCache()
	return repo, err
}

func (r *Repo) ReInitCache() error {
	r.ConfigCache.Lock()
	r.ConfigCache.DataMap = make(map[string]*configdomain.Config)

	rows, err := global.DBCON.Raw("SELECT id, value FROM config").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ID sql.NullString
		var Value sql.NullString

		rows.Scan(&ID, &Value)

		config := &configdomain.Config{}
		if ID.Valid {
			config.ID = ID.String
		}

		if Value.Valid {
			config.Value = Value.String
		}

		r.ConfigCache.DataMap[config.ID] = config
	}

	r.ConfigCache.Unlock()

	return nil
}

func (r *Repo) GetValue(key string) string {
	if config, ok := r.ConfigCache.DataMap[key]; ok {
		return config.Value
	}
	return ""
}
