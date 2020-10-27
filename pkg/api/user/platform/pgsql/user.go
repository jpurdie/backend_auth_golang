package pgsql

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/jpurdie/authapi"
)

type User struct{}


func (u User) FetchByExternalID(db orm.DB, externalID string) (authapi.User, error){
	op := "FetchByExternalID"
	us := authapi.User{}
	//get user information
	err := db.Model(&us).
		Where("\"user\".external_id = ?", externalID).
		Select()

	var profiles []authapi.Profile

	err = db.Model(&profiles).
		Column("profile.*").
		Relation("Organization").
		Relation("Role").
		Where("\"profile\".user_id = ?", us.ID).
		Select()

	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	us.Profile = profiles

	return us, nil
}

func (u User) Update(db orm.DB, p authapi.Profile) error {
	op := "Update"

	_, err := db.Model(&p).
		Set("role_id = ?role_id").
		Set("updated_at = NOW()").
		Where("id = ?id").
		Update()

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

//
func (u User) ListRoles(db orm.DB) ([]authapi.Role, error) {
	op := "ListRoles"
	var roles []authapi.Role

	err := db.Model(&roles).
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

func (u User) FetchByEmail(db orm.DB, email string) (authapi.User, error) {
	op := "Fetch"
	myUser := authapi.User{}
	//get user information
	err := db.Model(&myUser).
		Where("\"user\".email = ?", email).
		Select()

	var profiles []authapi.Profile

	err = db.Model(&profiles).
		Column("profile.*").
		Relation("Organization").
		Relation("Role").
		Where("\"profile\".user_id = ?", myUser.ID).
		Select()

	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	myUser.Profile = profiles

	return myUser, nil

}

func (u User) List(db orm.DB, orgID int) ([]authapi.User, error) {
	op := "List"
	var returnUsers []authapi.User
	var users []authapi.User
	err := db.Model(&users).
		Join("JOIN profiles AS p ON p.user_id = \"user\".id").
		Join("JOIN organizations AS o ON o.id = p.organization_id").
		Where("\"o\".id = ?", orgID).
		Select()

	for _, tempUser := range users {
		var profiles []authapi.Profile

		err = db.Model(&profiles).
			Column("profile.*").
			Relation("Organization").
			Relation("Role").
			Where("\"profile\".user_id = ?", tempUser.ID).
			Where("\"organization\".id = ?", orgID).
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
func (u User) ListAuthorized(db orm.DB, us *authapi.User, includeInactive bool) ([]authapi.Profile, error) {
	op := "ListAuthorized"
	var profile []authapi.Profile
	inactiveSQL := "organization.active = TRUE"
	if includeInactive {
		inactiveSQL = "1=1" //will return inactive and active
	}
	err := db.Model(&profile).
		Column("profile.*").
		Relation("Organization").
		Relation("User").
		Relation("Role").
		//	Join("JOIN organization_users AS cu ON cu.organization_id = organization.id").
		//	Join("JOIN users AS u ON cu.user_id = u.id").
		Where("external_id = ?", us.ExternalID).
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

func (u User) FetchProfile(db orm.DB, userID int, orgID int) (authapi.Profile, error) {
	op := "FetchProfile"
	var profile authapi.Profile

	err := db.Model(&profile).
		Join("JOIN users AS u ON profile.user_id = \"u\".id").
		Join("JOIN organizations AS o ON o.id = profile.organization_id").
		Where("\"o\".ID = ?", orgID).
		Where("\"u\".ID = ?", userID).
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
