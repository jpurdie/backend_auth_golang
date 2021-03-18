package pgsql

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
)

type User struct{}

//

func (u User) Update(db sqlx.DB, userToUpdate authapi.User) error {
	op := "Update"
	sql := "UPDATE users SET first_name=$1, last_name=$2, timezone=$3, updated_at=now() WHERE id=$4"
	_, err := db.Exec(sql, userToUpdate.FirstName, userToUpdate.LastName, userToUpdate.TimeZone, userToUpdate.ID)

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return nil
}

func (u User) FetchUserByUUID(db sqlx.DB, userUUID uuid.UUID, orgID int) (*authapi.User, error) {
	op := "FetchProfileByUUID"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where uuid=$1"
	err := db.QueryRowx(query, userUUID).StructScan(&us)

	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1 and organization_id=$2"
	err = db.Select(&profiles, query, us.ID, orgID)
	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return &us, nil
}

func (u User) FetchByEmail(db sqlx.DB, email string) (*authapi.User, error) {
	op := "FetchByEmail"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where email=$1"
	err := db.QueryRowx(query, email).StructScan(&us)

	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1"
	err = db.Select(&profiles, query, us.ID)
	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return &us, nil
}

func (u User) FetchByExternalID(db sqlx.DB, externalID string) (*authapi.User, error) {
	op := "FetchByExternalID"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where external_id=$1"
	err := db.QueryRowx(query, externalID).StructScan(&us)

	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1"
	err = db.Select(&profiles, query, us.ID)
	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return &us, nil
}

func (u User) FetchUserByID(db sqlx.DB, id int) (*authapi.User, error) {
	op := "FetchUserByID"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where id=$1"
	err := db.QueryRowx(query, id).StructScan(&us)

	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1"
	err = db.Select(&profiles, query, us.ID)
	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return &us, nil
}

func (u User) List(db sqlx.DB, orgID int) ([]authapi.User, error) {
	op := "List"
	var users []authapi.User

	//get user information
	query := "SELECT u.* FROM users u " +
		"LEFT JOIN profiles p on p.user_id = u.id " +
		"where p.organization_id=$1"
	err := db.Select(&users, query, orgID)

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, user := range users {
		tempUser, err := u.FetchUserByID(db, user.ID)
		if err != nil {
			return nil, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		users[i] = *tempUser

	}
	return users, nil
}

func (u User) UpdateRole(db sqlx.DB, level int, profileID int) error {
	op := "Update"

	_, err := db.Exec("UPDATE profiles set role_id=$1 WHERE id=$2", level, profileID)
	if err != nil {
		log.Println("There was a SQL error")
		log.Println(err)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (u User) ListRoles(db sqlx.DB) ([]authapi.Role, error) {
	op := "ListRoles"
	var roles []authapi.Role

	//get user information
	query := "SELECT * from roles"
	err := db.Select(&roles, query)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return roles, nil
}
