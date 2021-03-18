package pgsql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
)

type StrategicAlignment struct {
}

func (sa StrategicAlignment) Delete(db sqlx.DB, saUUID uuid.UUID, orgID int) error {
	op := "Delete"
	sql := "UPDATE strategic_alignments SET deleted_at=now() WHERE uuid=$1 and organization_id=$2;"
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

func (sa StrategicAlignment) List(db sqlx.DB, oID int) ([]model.StrategicAlignment, error) {
	op := "List"
	alignments := make([]model.StrategicAlignment, 0)
	sql := "SELECT * FROM strategic_alignments where organization_id=$1 AND deleted_at is null;"
	err := db.Select(&alignments, sql, oID)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return alignments, nil
}

func (sa StrategicAlignment) Update(db sqlx.DB, mySA model.StrategicAlignment) error {
	op := "Update"
	sql := "UPDATE strategic_alignments set updated_at=now(), name=$1 WHERE uuid=$2 AND organization_id=$3 AND deleted_at is null;"
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

func (sa StrategicAlignment) Create(db sqlx.DB, newSA model.StrategicAlignment) (model.StrategicAlignment, error) {
	op := "Create"
	sql := "INSERT INTO strategic_alignments (created_at, uuid, name, organization_id) VALUES (now(), $1, $2, $3);"
	_, err := db.Exec(sql, newSA.UUID.String(), newSA.Name, newSA.OrganizationID)
	if err != nil {
		return model.StrategicAlignment{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return newSA, nil
}
