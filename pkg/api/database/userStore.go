package database

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
)

type UserStore struct {
	db *pg.DB
}

// NewAdmAccountStore returns an AccountStore.
func NewUserStore(db *pg.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) Update(ou authapi.OrganizationUser) error {
	op := "Update"
	_, err := s.db.Model(&ou).Set("role_id = ?role_id").Where("uuid = ?uuid").Update()

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (s *UserStore) ListRoles() ([]authapi.Role, error) {
	op := "ListRoles"
	var roles []authapi.Role

	err := s.db.Model(&roles).
		Where("active = ?", true).
		Select()

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return roles, nil
}

//Fetch(ou authapi.OrganizationUser) (authapi.OrganizationUser, error)
func (s *UserStore) Fetch(ou authapi.OrganizationUser) (authapi.OrganizationUser, error) {
	op := "Fetch"
	var orgUser authapi.OrganizationUser

	err := s.db.Model(&orgUser).
		Column("organization_user.*").
		Relation("Organization").
		Relation("User").
		Relation("Role").
		Where("organization.uuid = ?", ou.Organization.UUID).
		Where("user.external_id = ?", ou.User.ExternalID).
		Where("organization.active = ?", true).
		Where("organization_user.active = ?", true).
		Order("organization.name asc").
		Select()
	if err != nil {
		return authapi.OrganizationUser{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return orgUser, nil

}

func (s *UserStore) List(o *authapi.Organization) ([]authapi.OrganizationUser, error) {
	op := "List"
	var orgUsers []authapi.OrganizationUser

	err := s.db.Model(&orgUsers).
		Column("organization_user.*").
		Relation("Organization").
		Relation("User").
		Relation("Role").
		Where("organization.id = ?", o.ID).
		Where("organization.active = ?", true).
		Order("organization.name asc").
		Select()
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return orgUsers, nil
}
func (s *UserStore) ListAuthorized(u *authapi.User, includeInactive bool) ([]authapi.OrganizationUser, error) {
	op := "ListAuthorized"
	var OrganizationUser []authapi.OrganizationUser
	inactiveSQL := "organization.active = TRUE"
	if includeInactive {
		inactiveSQL = "1=1" //will return inactive and active
	}
	err := s.db.Model(&OrganizationUser).
		Column("organization_user.*").
		Relation("Organization").
		Relation("User").
		Relation("Role").
		//	Join("JOIN organization_users AS cu ON cu.organization_id = organization.id").
		//	Join("JOIN users AS u ON cu.user_id = u.id").
		Where("external_id = ?", u.ExternalID).
		Where(inactiveSQL).
		Order("organization.name asc").
		Select()

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return OrganizationUser, nil
}
