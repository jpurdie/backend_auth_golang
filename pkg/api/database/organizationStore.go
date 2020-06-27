package database

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
)

type OrganizationStore struct {
	db *pg.DB
}

func NewOrganizationStore(db *pg.DB) *OrganizationStore {
	return &OrganizationStore{
		db: db,
	}
}

func (s *OrganizationStore) ListAccessible(u *authapi.User, includeInactive bool) ([]authapi.Organization, error) {
	op := "ListAccessible"
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
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return companies, nil
}
