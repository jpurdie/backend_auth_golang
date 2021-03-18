package pgsql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
)

type SponsorArea struct {
}

func (sa SponsorArea) Delete(db sqlx.DB, saUUID uuid.UUID, orgID int) error {
	op := "Delete"
	sql := "UPDATE sponsor_areas SET deleted_at=now() WHERE uuid=$1 and organization_id=$2;"
	_, err := db.Exec(sql, saUUID, orgID)
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (sa SponsorArea) List(db sqlx.DB, oID int) ([]model.SponsorArea, error) {
	op := "List"
	sponsorAreas := make([]model.SponsorArea, 0)
	sql := "SELECT * FROM sponsor_areas where organization_id=$1 AND deleted_at is null;"
	err := db.Select(&sponsorAreas, sql, oID)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return sponsorAreas, nil
}

func (sa SponsorArea) Update(db sqlx.DB, mySA model.SponsorArea) error {
	op := "Update"
	sql := "UPDATE sponsor_areas set updated_at=now(), name=$1 WHERE uuid=$2 AND organization_id=$3 AND deleted_at is null;"
	_, err := db.Exec(sql, mySA.Name, mySA.UUID, mySA.OrganizationID)
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (sa SponsorArea) Create(db sqlx.DB, newSA model.SponsorArea) (model.SponsorArea, error) {
	op := "Create"
	sql := "INSERT INTO sponsor_areas (created_at, uuid, name, organization_id) VALUES (now(), $1, $2, $3);"
	_, err := db.Exec(sql, newSA.UUID.String(), newSA.Name, newSA.OrganizationID)
	if err != nil {
		return model.SponsorArea{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return newSA, nil
}
