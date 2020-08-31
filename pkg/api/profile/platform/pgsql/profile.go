package pgsql

import (
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/jpurdie/authapi"
	"log"
)

// Custom errors
var (
	ErrCompAlreadyExists  = errors.New("Organization name already exists")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

type Profile struct{}


// Create creates a new user on database
func (p Profile) Create(tx *pg.Tx, profile authapi.Profile) error {
	op := "Create"
	log.Print(op)
	//var organization = new(authapi.Organization)

	//count, err := db.Model(organization).Where("lower(name) = ? and deleted_at is null", strings.ToLower(profile.Organization.Name)).Count()
	//if err != nil {
	//	return &authapi.Error{
	//		Op:   op,
	//		Code: authapi.EINTERNAL,
	//		Err:  err,
	//	}
	//}
	//if count > 0 {
	//	return ErrCompAlreadyExists
	//}
	//var user = new(authapi.User)
	//
	//count, err = db.Model(user).Where("lower(email) = ? and deleted_at is null", strings.ToLower(profile.User.Email)).Count()
	//
	//if err != nil {
	//	return &authapi.Error{
	//		Op:   op,
	//		Code: authapi.EINTERNAL,
	//		Err:  err,
	//	}
	//}
	//if count > 0 {
	//	return ErrEmailAlreadyExists
	//}


	tx, err := tx.Begin()
	trErr := tx.Insert(profile.Organization)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	profile.User.OrganizationID = profile.Organization.ID
	trErr = tx.Insert(profile.User)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	profile.UserID = profile.User.ID
	profile.OrganizationID = profile.Organization.ID
	trErr = tx.Insert(&profile)
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
