package webuserrepo

import (
	"database/sql"
	"dromatech/pos-backend/global"
	webuserdomain "dromatech/pos-backend/internal/domain/webuser"

	"github.com/sirupsen/logrus"
)

type WebUserRepo interface {
	FindAll() ([]*webuserdomain.WebUser, error)
	Find(id string) *webuserdomain.WebUser
	FindByUsername(username string) *webuserdomain.WebUser
	EditUser(webUser *webuserdomain.WebUser) error
	ChangePassword(userId string, newPassword string)
	RegisterUser(webUser *webuserdomain.WebUser)
	ChangeStatus(userId string, active bool)
}

type Repo struct {
}

func New() *Repo {
	repo := &Repo{}
	return repo
}

func (r *Repo) FindAll() ([]*webuserdomain.WebUser, error) {
	rows, err := global.DBCON.Raw("SELECT id, name, username, password_hash, password_salt, email, role_id, active, registration_timestamp, created_by FROM web_user ORDER BY name").Rows()
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []*webuserdomain.WebUser

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

		users = append(users, user)
	}

	return users, nil
}

func (r *Repo) Find(id string) *webuserdomain.WebUser {
	row := global.DBCON.Raw("SELECT id, name, username, password_hash, password_salt, email, role_id, active, registration_timestamp, created_by FROM web_user WHERE id = ?", id).Row()

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
	if ID.Valid && ID.String != "" {
		user.ID = ID.String
	} else {
		return nil
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

func (r *Repo) EditUser(webUser *webuserdomain.WebUser) error {
	return global.DBCON.Exec("UPDATE public.web_user "+
		"SET name=?, username=?, role_id=?, active=? "+
		"WHERE id=?;", webUser.Name, webUser.Username, webUser.RoleId, webUser.Active, webUser.ID).Error
}

func (r *Repo) ChangePassword(userId string, newPassword string) {
	global.DBCON.Exec("UPDATE public.web_user "+
		"SET password_hash=? "+
		"WHERE id=?;", newPassword, userId)
}

func (r *Repo) RegisterUser(webUser *webuserdomain.WebUser) {
	global.DBCON.Exec("INSERT INTO public.web_user(id, username, password_hash, password_salt, email, role_id, active, registration_timestamp, created_by, name) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		webUser.ID, webUser.Username, webUser.PasswordHash, webUser.PasswordSalt, webUser.Email, webUser.RoleId, webUser.Active, webUser.RegistrationTimestamp, webUser.CreatedBy, webUser.Name)
}

func (r *Repo) ChangeStatus(userId string, active bool) {
	global.DBCON.Exec("UPDATE public.web_user "+
		"SET active=? "+
		"WHERE id=?;", active, userId)
}
