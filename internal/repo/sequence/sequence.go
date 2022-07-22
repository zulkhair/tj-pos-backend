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

func nextVal(id string, tx *gorm.DB) int64 {
	row := tx.Raw("SELECT id, next_value FROM public.sequence WHERE id = ?", id).Row()

	var ID sql.NullString
	var NextValue sql.NullInt64

	row.Scan(&ID, &NextValue)

	if !ID.Valid || ID.String == "" {
		return int64(1)
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
	nextVal := nextVal(id, tx)

	if nextVal == 1 {
		tx.Exec("INSERT INTO public.sequence(id, next_value) VALUES (?, ?);", id, 2)
		return nextVal
	} else {
		tx.Exec("UPDATE public.sequence SET next_value=? WHERE id=?;", nextVal+1, id)
		return nextVal
	}
}
