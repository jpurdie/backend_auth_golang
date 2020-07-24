package database

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
	"log"
	"strings"
)

type AuthOrganizationStore struct {
	db *pg.DB
}

func NewAuthOrganizationStore(db *pg.DB) *AuthOrganizationStore {
	return &AuthOrganizationStore{
		db: db,
	}
}

// Custom errors
var (
	ErrCompAlreadyExists  = errors.New("Organization name already exists")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

// Create creates a new user on database
func (s *AuthOrganizationStore) Create(cu authapi.OrganizationUser) error {
	op := "Create"
	var organization = new(authapi.Organization)

	count, err := s.db.Model(organization).Where("lower(name) = ? and deleted_at is null", strings.ToLower(cu.Organization.Name)).Count()
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if count > 0 {
		return ErrCompAlreadyExists
	}
	var user = new(authapi.User)

	count, err = s.db.Model(user).Where("lower(email) = ? and deleted_at is null", strings.ToLower(cu.User.Email)).Count()

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if count > 0 {
		return ErrEmailAlreadyExists
	}

	tx, err := s.db.Begin()
	trErr := tx.Insert(cu.Organization)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	cu.User.OrganizationID = cu.Organization.ID
	trErr = tx.Insert(cu.User)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	cu.UserID = cu.User.ID
	cu.OrganizationID = cu.Organization.ID
	trErr = tx.Insert(&cu)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	trErr = tx.Commit()
	if trErr != nil {
		log.Println("There was a transaction error")
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	log.Println("Organization User creation was successful")
	return nil
}
