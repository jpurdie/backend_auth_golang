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

func (s *UserStore) List(o *authapi.Organization) ([]authapi.OrganizationUser, error) {
	op := "List"
	var orgUsers []authapi.OrganizationUser

	err := s.db.Model(&orgUsers).
		Column("organization_user.*").
		Relation("Organization").
		Relation("User").
		Relation("Role").
		Where("organization.id = ?", o.ID).
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
