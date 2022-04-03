package sequencerepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	"gorm.io/gorm"
)

type SequenceRepo interface {
	NextVal(id string) int64
	NextValTx(id string, tx *gorm.DB) int64
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func nextVal(id string) int64 {
	row := global.DBCON.Raw("SELECT id, value FROM public.sequence WHERE id = ?", id).Row()

	var ID sql.NullString
	var NextValue sql.NullInt64

	row.Scan(&ID, &NextValue)

	if !ID.Valid || ID.String == "" {
		return int64(0)
	} else {
		return NextValue.Int64
	}
}

func (r *Repo) NextVal(id string) int64 {
	tx := global.DBCON.Begin()
	nextVal := r.NextValTx(id, tx)
	tx.Commit()

	return nextVal
}

func (r *Repo) NextValTx(id string, tx *gorm.DB) int64 {
	nextVal := nextVal(id)

	if nextVal == 0 {
		tx.Exec("INSERT INTO public.sequence(id, next_value) VALUES (?, ?);", id, 0)
		return int64(0)
	} else {
		tx.Exec("UPDATE public.sequence SET next_value=? WHERE id=?;", nextVal, id)
		return nextVal
	}
}
