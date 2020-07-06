package database

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
	"log"
	"strings"
	"time"
)

type AuthInvitationStore struct {
	db *pg.DB
}

func NewAuthInvitationStore(db *pg.DB) *AuthInvitationStore {
	return &AuthInvitationStore{
		db: db,
	}
}
func (is *AuthInvitationStore) View(i authapi.Invitation) (authapi.Invitation, error) {
	op := "VerifyToken"
	invite := new(authapi.Invitation)

	err := is.db.Model(invite).
		Relation("Organization").
		Where("token_hash = ?", i.TokenHash).
		Where("organization.active = TRUE").
		//Join("JOIN organizations org ON org.id = \"invitation\".\"organization_id\"").
		First()

	if err != nil {
		log.Println(err)
		return authapi.Invitation{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return *invite, nil
}

// Create creates a new user on database
func (s *AuthInvitationStore) CreateUser(cu authapi.OrganizationUser, i authapi.Invitation) error {
	op := "Create"
	//	var organization = new(authapi.Organization)

	var user = new(authapi.User)

	count, err := s.db.Model(user).Where("lower(email) = ? and deleted_at is null", strings.ToLower(cu.User.Email)).Count()

	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if count > 0 {
		return ErrEmailAlreadyExists
	}

	tx, err := s.db.Begin()
	cu.User.OrganizationID = cu.Organization.ID
	trErr := tx.Insert(cu.User)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	cu.UserID = cu.User.ID
	cu.OrganizationID = cu.Organization.ID
	trErr = tx.Insert(&cu)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	i.Used = true
	i.UpdatedAt = time.Now()
	// res, err := db.Model(book).Set("title = ?title").Where("id = ?id").Update()
	_, trErr = tx.Model(&i).Set("used= ?used").Set("updated_at=now()").Where("id = ?id").Update()
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
