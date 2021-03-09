package pgsql

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/lib/pq"
)

type User struct{}

func (u User) FetchUserByID(db sqlx.DB, id int) (authapi.User, error) {
	op := "FetchUserByID"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where id=$1"
	err := db.QueryRowx(query, id).StructScan(&us)

	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1"
	err = db.Select(&profiles, query, us.ID)
	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return us, nil
}

func (u User) FetchUserByUUID(db sqlx.DB, userUUID uuid.UUID, orgID int) (userFromDB *authapi.User, error error) {
	op := "FetchProfileByUUID"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where uuid=$1"
	err := db.QueryRowx(query, userUUID).StructScan(&us)

	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1 and organization_id=$2"
	err = db.Select(&profiles, query, us.ID, orgID)
	if err != nil {
		return &authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return &authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return &us, nil
}

func (u User) FetchByEmail(db sqlx.DB, email string) (authapi.User, error) {
	op := "FetchByEmail"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where email=$1"
	err := db.QueryRowx(query, email).StructScan(&us)

	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1"
	err = db.Select(&profiles, query, us.ID)
	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return us, nil
}

func (u User) FetchByExternalID(db sqlx.DB, externalID string) (authapi.User, error) {
	op := "FetchByExternalID"
	us := authapi.User{}
	//get user information
	query := "SELECT * FROM users where external_id=$1"
	err := db.QueryRowx(query, externalID).StructScan(&us)

	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	profiles := []authapi.Profile{}
	query = "SELECT * FROM profiles where user_id=$1"
	err = db.Select(&profiles, query, us.ID)
	if err != nil {
		return authapi.User{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, profile := range profiles {
		org := []authapi.Organization{}
		query = "SELECT * FROM organizations where id=$1"
		err = db.Select(&org, query, profile.OrganizationID)
		if err != nil {
			return authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Organization = &org[0]
	}

	for i, profile := range profiles {
		role := []authapi.Role{}
		query = "SELECT r.* FROM profiles p JOIN roles r on r.id = p.role_id where p.id=$1"
		err = db.Select(&role, query, profile.ID)
		if err != nil {
			return authapi.User{}, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		profiles[i].Role = &role[0]
	}

	us.Profile = profiles
	return us, nil
}

func (u User) UpdateRole(db sqlx.DB, level int, profileID int) error {
	op := "Update"

	_, err := db.Exec("UPDATE profiles set role_id=$1 WHERE id=$2", level, profileID)
	if err != nil {
		log.Println("There was a SQL error")
		log.Println(err)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

//
func (u User) ListRoles(db sqlx.DB) ([]authapi.Role, error) {
	op := "ListRoles"
	var roles []authapi.Role

	//get user information
	query := "SELECT * from roles"
	err := db.Select(&roles, query)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return roles, nil
}

func (u User) List(db sqlx.DB, orgID int) ([]authapi.User, error) {
	op := "List"
	var users []authapi.User

	//get user information
	query := "SELECT u.* FROM users u " +
		"LEFT JOIN profiles p on p.user_id = u.id " +
		"where p.organization_id=$1"
	err := db.Select(&users, query, orgID)

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for i, user := range users {
		tempUser, err := u.FetchUserByID(db, user.ID)
		if err != nil {
			return nil, &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
		users[i] = tempUser

	}

	return users, nil
	//var users []authapi.User
	//err := db.Model(&users).
	//	Join("JOIN profiles AS p ON p.user_id = \"user\".id").
	//	Join("JOIN organizations AS o ON o.id = p.organization_id").
	//	Where("\"o\".id = ?", orgID).
	//	Select()
	//
	//for _, tempUser := range users {
	//	var profiles []authapi.Profile
	//
	//	err = db.Model(&profiles).
	//		Column("profile.*").
	//		Relation("Organization").
	//		Relation("Role").
	//		Where("\"profile\".user_id = ?", tempUser.ID).
	//		Where("\"organization\".id = ?", orgID).
	//		Select()
	//
	//	tempUser.Profile = profiles
	//	if profiles != nil {
	//		returnUsers = append(returnUsers, tempUser)
	//	}
	//
	//}

}
func (u User) Update(db sqlx.DB, userToUpdate authapi.User) error {
	op := "Update"
	sql := "UPDATE users SET first_name=$1, last_name=$2, timezone=$3, updated_at=now() WHERE id=$4"
	_, err := db.Exec(sql, userToUpdate.FirstName, userToUpdate.LastName, userToUpdate.TimeZone, userToUpdate.ID)

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return nil
}

func (u User) FetchProfile(db sqlx.DB, userID int, orgID int) (authapi.Profile, error) {
	op := "FetchProfile"
	var prof authapi.Profile

	query := "SELECT " +
		"u.id as userID, " +
		"u.uuid as userUUID, " +
		"u.first_name as firstName, " +
		"u.last_name as lastName, " +
		"u.email as email, " +
		"p.active as profileActive, " +
		"p.role_id as roleID, " +
		"p.uuid as profileUUID, " +
		"p.created_at as profileCreatedAt, " +
		"p.id as profileID, " +
		"o.name as orgName, " +
		"o.id as orgID, " +
		"o.uuid as orgUUID " +
		"FROM users u " +
		"LEFT JOIN profiles p on p.user_id = u.id " +
		"LEFT JOIN organizations o on p.organization_id = o.id " +
		"LEFT JOIN roles r on p.role_id = r.id " +
		"where p.user_id=$1 " +
		"AND p.organization_id=$2 "

	rows, err := db.Queryx(query, userID, orgID)

	for rows.Next() {
		var userID int
		var userUUID uuid.UUID
		var firstName string
		var lastName string
		var email string
		var profileActive bool
		var roleID authapi.AccessRole
		var profileUUID uuid.UUID
		var profileCreatedAt pq.NullTime
		var profileID int
		var orgID int
		var orgName string
		var orgUUID uuid.UUID

		err = rows.Scan(&userID, &userUUID, &firstName, &lastName, &email, &profileActive, &roleID, &profileUUID, &profileCreatedAt, &profileID, &orgName, &orgID, &orgUUID)

		us := authapi.User{}
		org := authapi.Organization{}
		us.ID = userID
		us.UUID = userUUID
		us.FirstName = firstName
		us.LastName = lastName
		us.Email = email

		org.UUID = orgUUID
		org.Name = orgName
		org.ID = orgID

		prof.UUID = profileUUID
		prof.ID = profileID
		prof.CreatedAt = profileCreatedAt
		prof.User = &us
		prof.Organization = &org
	}

	if err != nil {
		return authapi.Profile{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return prof, nil
}
