package database

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
	"log"
	"strings"
)

type OrganizationStore struct {
	db *pg.DB
}

func NewOrganizationStore(db *pg.DB) *OrganizationStore {
	return &OrganizationStore{
		db: db,
	}
}

//
// Custom errors
var (
	ErrCompAlreadyExists  = errors.New("Organization name already exists")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

func (s *OrganizationStore) ListAccessible(u *authapi.User, includeInactive bool) ([]authapi.Organization, error) {
	var companies []authapi.Organization
	inactiveSQL := "organization.active = TRUE"
	if includeInactive {
		inactiveSQL = "1=1" //will return inactive and active
	}
	err := s.db.Model(&companies).
		Join("JOIN organization_users AS cu ON cu.organization_id = organization.id").
		Join("JOIN users AS u ON cu.user_id = u.id").
		Where("u.external_id = ?", u.ExternalID).
		Where(inactiveSQL).
		Order("organization.name asc").
		Select()

	if err != nil {
		return nil, err
	}
	return companies, nil
}

// Create creates a new user on database
func (s *OrganizationStore) Create(cu authapi.OrganizationUser) error {

	var organization = new(authapi.Organization)

	count, err := s.db.Model(organization).Where("lower(name) = ? and deleted_at is null", strings.ToLower(cu.Organization.Name)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCompAlreadyExists
	}
	var user = new(authapi.User)

	count, err = s.db.Model(user).Where("lower(email) = ? and deleted_at is null", strings.ToLower(cu.User.Email)).Count()

	if err != nil {
		return err
	}
	if count > 0 {
		return ErrEmailAlreadyExists
	}

	tx, err := s.db.Begin()
	trErr := tx.Insert(cu.Organization)
	if trErr != nil {
		log.Println(trErr)
	}
	cu.User.OrganizationID = cu.Organization.ID
	trErr = tx.Insert(cu.User)
	if trErr != nil {
		log.Println(trErr)
	}
	cu.UserID = cu.User.ID
	cu.OrganizationID = cu.Organization.ID
	trErr = tx.Insert(&cu)
	if trErr != nil {
		log.Println(trErr)
	}
	trErr = tx.Commit()
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
	}
	return err
}
