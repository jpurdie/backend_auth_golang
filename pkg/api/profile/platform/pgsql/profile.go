package pgsql

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/lib/pq"
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

func (p Profile) FetchProfileByExternalID(db sqlx.DB, externalID string) (authapi.Profile, error) {
	op := "FetchProfile"
	var prof authapi.Profile

	query := "SELECT " +
		"u.id as userID, " +
		"u.uuid as userUUID, " +
		"u.first_name as firstName, " +
		"u.last_name as lastName, " +
		"u.email as email, " +
		"p.active as profileActive, " +
		"p.role_id as roleID, " +
		"p.uuid as profileUUID, " +
		"p.created_at as profileCreatedAt, " +
		"p.id as profileID, " +
		"o.name as orgName, " +
		"o.id as orgID, " +
		"o.uuid as orgUUID " +
		"FROM users u " +
		"LEFT JOIN profiles p on p.user_id = u.id " +
		"LEFT JOIN organizations o on p.organization_id = o.id " +
		"LEFT JOIN roles r on p.role_id = r.id " +
		"where u.external_id=$1 "

	rows, err := db.Queryx(query, externalID)

	for rows.Next() {
		var userID int
		var userUUID uuid.UUID
		var firstName string
		var lastName string
		var email string
		var profileActive bool
		var roleID authapi.AccessRole
		var profileUUID uuid.UUID
		var profileCreatedAt pq.NullTime
		var profileID int
		var orgID int
		var orgName string
		var orgUUID uuid.UUID

		err = rows.Scan(&userID, &userUUID, &firstName, &lastName, &email, &profileActive, &roleID, &profileUUID, &profileCreatedAt, &profileID, &orgName, &orgID, &orgUUID)

		us := authapi.User{}
		org := authapi.Organization{}
		us.ID = userID
		us.UUID = userUUID
		us.FirstName = firstName
		us.LastName = lastName
		us.Email = email

		org.UUID = orgUUID
		org.Name = orgName
		org.ID = orgID

		prof.UUID = profileUUID
		prof.ID = profileID
		prof.CreatedAt = profileCreatedAt
		prof.User = &us
		prof.Organization = &org
	}

	if err != nil {
		return authapi.Profile{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return prof, nil
}

func (p Profile) FetchProfile(db sqlx.DB, userID int, orgID int) (authapi.Profile, error) {
	op := "FetchProfile"
	var prof authapi.Profile

	query := "SELECT " +
		"u.id as userID, " +
		"u.uuid as userUUID, " +
		"u.first_name as firstName, " +
		"u.last_name as lastName, " +
		"u.email as email, " +
		"p.active as profileActive, " +
		"p.role_id as roleID, " +
		"p.uuid as profileUUID, " +
		"p.created_at as profileCreatedAt, " +
		"p.id as profileID, " +
		"o.name as orgName, " +
		"o.id as orgID, " +
		"o.uuid as orgUUID " +
		"FROM users u " +
		"LEFT JOIN profiles p on p.user_id = u.id " +
		"LEFT JOIN organizations o on p.organization_id = o.id " +
		"LEFT JOIN roles r on p.role_id = r.id " +
		"where p.user_id=$1 " +
		"AND p.organization_id=$2 "

	rows, err := db.Queryx(query, userID, orgID)

	for rows.Next() {
		var userID int
		var userUUID uuid.UUID
		var firstName string
		var lastName string
		var email string
		var profileActive bool
		var roleID authapi.AccessRole
		var profileUUID uuid.UUID
		var profileCreatedAt pq.NullTime
		var profileID int
		var orgID int
		var orgName string
		var orgUUID uuid.UUID

		err = rows.Scan(&userID, &userUUID, &firstName, &lastName, &email, &profileActive, &roleID, &profileUUID, &profileCreatedAt, &profileID, &orgName, &orgID, &orgUUID)

		us := authapi.User{}
		org := authapi.Organization{}
		us.ID = userID
		us.UUID = userUUID
		us.FirstName = firstName
		us.LastName = lastName
		us.Email = email

		org.UUID = orgUUID
		org.Name = orgName
		org.ID = orgID

		prof.UUID = profileUUID
		prof.ID = profileID
		prof.CreatedAt = profileCreatedAt
		prof.User = &us
		prof.Organization = &org
	}

	if err != nil {
		return authapi.Profile{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return prof, nil
}
