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

func (s *UserStore) FetchProfile(u authapi.User, o authapi.Organization) (authapi.Profile, error) {
	op := "ListRoles"
	var profile authapi.Profile

	err := s.db.Model(&profile).
		Join("JOIN users AS u ON profile.user_id = \"u\".id").
		Join("JOIN organizations AS o ON o.id = profile.organization_id").
		Where("\"o\".id = ?", o.ID).
		Where("\"u\".uuid = ?", u.UUID).
		Where("profile.active = ?", true).
		Where("o.active = ?", true).
		Select()

	if err != nil {
		return authapi.Profile{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return profile, nil
}

func (s *UserStore) Delete(p authapi.Profile) error {
	op := "Delete"

	_, err := s.db.Model(&p).Where("id = ?id").Delete()

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (s *UserStore) Update(p authapi.Profile) error {
	op := "Update"

	_, err := s.db.Model(&p).Set("role_id = ?role_id").Set("updated_at = NOW()").Where("id = ?id").Update()

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

func (s *UserStore) Fetch(u authapi.User) (authapi.User, error) {
	op := "Fetch"

	//get user information
	err := s.db.Model(&u).
		Where("\"user\".external_id = ?", u.ExternalID).
		Select()

	var profiles []authapi.Profile

	err = s.db.Model(&profiles).
		Column("profile.*").
		Relation("Organization").
		Relation("Role").
		Where("\"profile\".user_id = ?", u.ID).
		Select()

	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	u.Profile = profiles

	return u, nil

}

func (s *UserStore) List(o *authapi.Organization) ([]authapi.User, error) {
	op := "List"
	var returnUsers []authapi.User
	var users []authapi.User
	err := s.db.Model(&users).
		Join("JOIN profiles AS p ON p.user_id = \"user\".id").
		Join("JOIN organizations AS o ON o.id = p.organization_id").
		Where("\"o\".id = ?", o.ID).
		Select()

	for _, tempUser := range users {
		var profiles []authapi.Profile

		err = s.db.Model(&profiles).
			Column("profile.*").
			Relation("Organization").
			Relation("Role").
			Where("\"profile\".user_id = ?", tempUser.ID).
			Where("\"organization\".id = ?", o.ID).
			Select()

		tempUser.Profile = profiles
		if profiles != nil {
			returnUsers = append(returnUsers, tempUser)
		}

	}

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return returnUsers, nil
}
func (s *UserStore) ListAuthorized(u *authapi.User, includeInactive bool) ([]authapi.Profile, error) {
	op := "ListAuthorized"
	var profile []authapi.Profile
	inactiveSQL := "organization.active = TRUE"
	if includeInactive {
		inactiveSQL = "1=1" //will return inactive and active
	}
	err := s.db.Model(&profile).
		Column("profile.*").
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
	return profile, nil
}
