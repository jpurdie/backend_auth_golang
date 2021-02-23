package pgsql

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
)

type Invitation struct {
}

func (i Invitation) Create(db sqlx.DB, invite authapi.Invitation) error {
	op := "Create"

	sql := "INSERT INTO invitations " +
		"(token_hash, expires_at, invitor_id, organization_id, email, used) " +
		"VALUES " +
		"($1, $2, $3, $4, $5, $6); "
	tx := db.MustBegin()
	tx.MustExec(sql, invite.TokenHash, invite.ExpiresAt, invite.InvitorID, invite.OrganizationID, invite.Email, false)
	err := tx.Commit()
	if err != nil {
		log.Println(err)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return nil
}

func (i Invitation) Delete(db sqlx.DB, email string, orgID int) error {
	op := "Delete"
	sql := "UPDATE invitations set deleted_at=now() WHERE email=$1 AND organization_id=$2"
	tx := db.MustBegin()
	tx.MustExec(sql, email, orgID)
	err := tx.Commit()
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (i Invitation) List(db sqlx.DB, orgID int, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error) {
	op := "List"
	invitations := make([]authapi.Invitation, 0)

	sql := "SELECT * FROM invitations i " +
		"WHERE " +
		"deleted_at is null " +
		"AND i.organization_id=$1 "

	expiredSQL := "AND i.expires_at >= NOW() "
	if includeExpired {
		expiredSQL = ""
	}

	usedSQL := "AND i.used = FALSE "
	if includeUsed {
		usedSQL = ""
	}

	sql += usedSQL
	sql += expiredSQL

	err := db.Select(&invitations, sql, orgID)

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return invitations, nil
}

func (i Invitation) ViewByEmail(db sqlx.DB, email string, orgID int) (authapi.Invitation, error) {
	op := "View"
	invite := authapi.Invitation{}

	sql := "SELECT i.* " +
		"FROM invitations i " +
		"WHERE " +
		"deleted_at is null " +
		"AND i.used = false " +
		"AND i.email=$1 " +
		"AND i.organization_id=$2"

	err := db.QueryRowx(sql, email, orgID).StructScan(&invite)

	if err != nil {
		log.Println(err)
		return authapi.Invitation{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return invite, nil
}

func (i Invitation) View(db sqlx.DB, tokenHash string) (authapi.Invitation, error) {
	op := "View"
	invite := authapi.Invitation{}
	org := authapi.Organization{}

	sql := "SELECT i.* FROM invitations i " +
		"WHERE " +
		"deleted_at is null " +
		"AND i.used = false " +
		"AND i.token_hash=$1;"

	err := db.QueryRowx(sql, tokenHash).StructScan(&invite)

	if err != nil {
		log.Println(err)
		return authapi.Invitation{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	sql = "SELECT * FROM organizations where deleted_at is null AND id=$1;"
	err = db.QueryRowx(sql, invite.OrganizationID).StructScan(&org)

	if err != nil {
		log.Println(err)
		return authapi.Invitation{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	invite.Organization = &org

	return invite, nil
}

func (i Invitation) CreateUser(db sqlx.DB, profile authapi.Profile, invite authapi.Invitation) error {
	op := "CreateUser"

	tx, err := db.Beginx()
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
	err = tx.QueryRowx(queryProfile, profile.UUID, userID, invite.OrganizationID, profile.RoleID, true).Scan(&profileID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	queryInvitation := "UPDATE invitations set updated_at=now(), used=true WHERE id=$1;"
	_ = tx.MustExec(queryInvitation, invite.ID)

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
