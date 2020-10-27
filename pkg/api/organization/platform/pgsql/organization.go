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

type Organization struct{}


//func (p Organization)FindByOr(db *pg.DB, email string ) (authapi.Profile, error) {
//	user, err := s.repo.FindByEmail(email)
//	if user != nil {
//		return fmt.Errorf("%s already exists", email)
//	}
//	if err != nil {
//		return err
//	}
//	return nil
//}


// Create creates a new user on database

func (p Organization) Create(tx *pg.Tx, profile authapi.Profile) error {
	op := "Create"
	log.Print(op)

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
	//trErr = tx.Commit()
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
