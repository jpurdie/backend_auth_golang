package database

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
	"log"
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

	_, trErr := s.db.Model(&i).Returning("*").Insert()
	if trErr != nil {
		log.Println(trErr)
		return authapi.Invitation{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  trErr,
		}
	}

	return authapi.Invitation{}, nil
}

func (s *InvitationStore) Delete(i authapi.Invitation) error {
	op := "Delete"

	_, err := s.db.Model(&i).Where("email = ?email").Where("organization_id = ?organization_id").Delete()
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (s *InvitationStore) List(o *authapi.Organization, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error) {
	op := "List"
	invitations := make([]authapi.Invitation, 0)
	inactiveSQL := "invitation.expires_at >= NOW()"
	if includeExpired {
		inactiveSQL = "1=1" //will return inactive and active
	}
	usedSQL := "invitation.used = FALSE"
	if includeUsed {
		usedSQL = "1=1" //will return inactive and active
	}

	err := s.db.Model(&invitations).
		Where("invitation.organization_id = ?", o.ID).
		Where(inactiveSQL).
		Where(usedSQL).
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
