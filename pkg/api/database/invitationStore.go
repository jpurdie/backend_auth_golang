package database

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
)

type InvitationStore struct {
	db *pg.DB
}

func NewInvitationStore(db *pg.DB) *InvitationStore {
	return &InvitationStore{
		db: db,
	}
}

func (s *InvitationStore) Create(i authapi.Invitation) (authapi.Invitation, error) {

	op := "Create"
	_, err := s.db.Model(&i).Returning("*").Insert()

	if err != nil {
		return authapi.Invitation{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return authapi.Invitation{}, nil
}

func (s *InvitationStore) List(u *authapi.User, includeExpired bool) ([]authapi.Invitation, error) {
	op := "List"
	invitations := make([]authapi.Invitation, 0)
	inactiveSQL := "invitation.expires_at = NOW()"
	if includeExpired {
		inactiveSQL = "1=1" //will return inactive and active
	}
	err := s.db.Model(&invitations).
		Where("invitation.organization_id = ?", u.OrganizationID).
		Where(inactiveSQL).
		Order("invitation.expires_at").
		Select()

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return invitations, nil
}
