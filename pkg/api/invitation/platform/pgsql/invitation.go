package pgsql

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/jpurdie/authapi"
	"log"
	"time"
)

type Invitation struct{
}


func (i Invitation) Create(db orm.DB, invite authapi.Invitation) error {
	op := "Create"

	_, trErr := db.Model(&invite).Returning("*").Insert()
	if trErr != nil {
		log.Println(trErr)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  trErr,
		}
	}

	return nil
}

func (i Invitation) Delete(db orm.DB, invitation authapi.Invitation) error {
	op := "Delete"

	_, err := db.Model(&invitation).Where("email = ?email").Where("organization_id = ?organization_id").Delete()
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (i Invitation) List(db orm.DB, o *authapi.Organization, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error) {
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

	err := db.Model(&invitations).
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

func (i Invitation) View(db orm.DB, tokenHash string) (authapi.Invitation, error) {
	op := "View"
	invite := new(authapi.Invitation)

	err := db.Model(invite).
		Relation("Organization").
		Where("token_hash = ?", tokenHash).
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


func (i Invitation) CreateUser(tx *pg.Tx, cu authapi.Profile, invite authapi.Invitation) error {
	op := "CreateUser"

	cu.User.OrganizationID = cu.Organization.ID
	trErr := tx.Insert(cu.User)
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  trErr,
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
			Err:  trErr,
		}
	}
	invite.Used = true
	invite.UpdatedAt = time.Now()
	// res, err := db.Model(book).Set("title = ?title").Where("id = ?id").Update()
	_, trErr = tx.Model(&invite).Set("used= ?used").Set("updated_at=now()").Where("id = ?id").Update()
	if trErr != nil {
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  trErr,
		}
	}
//	trErr = tx.Commit()
	if trErr != nil {
		log.Println("There was a transaction error")
		log.Println(trErr)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  trErr,
		}
	}
	log.Println("Organization User creation was successful")
	return nil
}