package pgsql

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"log"
)

// Custom errors
var (
	ErrCompAlreadyExists  = errors.New("Organization name already exists")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

type Profile struct{}


// Create creates a new user on database
func (p Profile) Create(db sqlx.DB, profile authapi.Profile) error {
	op := "Create"

	tx, err := db.Beginx()
	orgID := 0
	err = tx.QueryRowx("INSERT INTO organizations (\"name\",active,uuid) VALUES ($1, $2, $3) RETURNING id as xID;", profile.Organization.Name, true, profile.Organization.UUID.String()).Scan(&orgID)
	if err != nil {

		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	userID := 0
	queryUser := "INSERT INTO users (\"first_name\", \"last_name\", \"email\", \"external_id\", \"uuid\") VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	err = tx.QueryRowx(queryUser, profile.User.FirstName, profile.User.LastName, profile.User.Email, profile.User.ExternalID, profile.User.UUID.String()).Scan(&userID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	profileID := 0

	queryProfile := "INSERT INTO profiles (\"uuid\",\"user_id\", \"organization_id\", \"role_id\", \"active\") VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = tx.QueryRowx(queryProfile, profile.UUID, userID, orgID, profile.RoleID, true).Scan(&profileID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("There was a transaction error")
		tx.Rollback()
		log.Println(err)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	log.Println("Organization User creation was successful")
	return nil
}
