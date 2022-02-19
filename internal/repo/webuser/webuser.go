package webuserrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	webuserdomain "dromatech/pos-backend/internal/domain/webuser"

	"github.com/sirupsen/logrus"
)

type WebUserRepo interface {
	ReInitCache() error
	Find(id string) *webuserdomain.WebUser
	FindByUsername(username string) *webuserdomain.WebUser
}

type Repo struct {
	WebUserCahce webuserdomain.WebUserCache
}

func New() (*Repo, error) {
	repo := &Repo{}
	err := repo.ReInitCache()
	return repo, err
}

func (r *Repo) ReInitCache() error {
	r.WebUserCahce.Lock()
	r.WebUserCahce.DataMap = make(map[string]*webuserdomain.WebUser)

	rows, err := global.DBCON.Raw("SELECT id, name, username, password_hash, password_salt, email, role_id, active, registration_timestamp, created_by FROM web_user").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ID sql.NullString
		var Name sql.NullString
		var Username sql.NullString
		var PasswordHash sql.NullString
		var PasswordSalt sql.NullString
		var Email sql.NullString
		var RoleId sql.NullString
		var Active sql.NullBool
		var RegistrationTimestamp sql.NullTime
		var CreatedBy sql.NullString

		rows.Scan(&ID, &Name, &Username, &PasswordHash, &PasswordSalt, &Email, &RoleId, &Active, &RegistrationTimestamp, &CreatedBy)

		user := &webuserdomain.WebUser{}
		if ID.Valid {
			user.ID = ID.String
		}

		if Name.Valid {
			user.Name = Name.String
		}

		if Username.Valid {
			user.Username = Username.String
		}

		if PasswordHash.Valid {
			user.PasswordHash = PasswordHash.String
		}

		if PasswordSalt.Valid {
			user.PasswordSalt = PasswordSalt.String
		}

		if Email.Valid {
			user.Email = Email.String
		}

		if RoleId.Valid {
			user.RoleId = RoleId.String
		}

		if Active.Valid {
			user.Active = Active.Bool
		}

		if RegistrationTimestamp.Valid {
			user.RegistrationTimestamp = RegistrationTimestamp.Time
		}

		if CreatedBy.Valid {
			user.CreatedBy = CreatedBy.String
		}

		r.WebUserCahce.DataMap[user.ID] = user
	}

	r.WebUserCahce.Unlock()

	return nil
}

func (r *Repo) Find(id string) *webuserdomain.WebUser {
	if webuser, ok := r.WebUserCahce.DataMap[id]; ok {
		return webuser
	}
	return nil
}

func (r *Repo) FindByUsername(username string) *webuserdomain.WebUser {
	row := global.DBCON.Raw("SELECT id, name, username, password_hash, password_salt, email, role_id, active, registration_timestamp, created_by FROM web_user WHERE username = ?", username).Row()

	var ID sql.NullString
	var Name sql.NullString
	var Username sql.NullString
	var PasswordHash sql.NullString
	var PasswordSalt sql.NullString
	var Email sql.NullString
	var RoleId sql.NullString
	var Active sql.NullBool
	var RegistrationTimestamp sql.NullTime
	var CreatedBy sql.NullString

	row.Scan(&ID, &Name, &Username, &PasswordHash, &PasswordSalt, &Email, &RoleId, &Active, &RegistrationTimestamp, &CreatedBy)

	user := &webuserdomain.WebUser{}
	if ID.Valid {
		user.ID = ID.String
	}

	if Name.Valid {
		user.Name = Name.String
	}

	if Username.Valid {
		user.Username = Username.String
	}

	if PasswordHash.Valid {
		user.PasswordHash = PasswordHash.String
	}

	if PasswordSalt.Valid {
		user.PasswordSalt = PasswordSalt.String
	}

	if Email.Valid {
		user.Email = Email.String
	}

	if RoleId.Valid {
		user.RoleId = RoleId.String
	}

	if Active.Valid {
		user.Active = Active.Bool
	}

	if RegistrationTimestamp.Valid {
		user.RegistrationTimestamp = RegistrationTimestamp.Time
	}

	if CreatedBy.Valid {
		user.CreatedBy = CreatedBy.String
	}

	return user
}
