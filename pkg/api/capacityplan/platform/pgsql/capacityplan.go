package pgsql

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/jpurdie/authapi"

	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/utl/model"
)

type CapacityPlan struct {
}

func (cp CapacityPlan) Create(db sqlx.DB, orgID int, cpToSave []model.CapacityPlanEntry) error {
	op := "Create"

	fmt.Println(cpToSave)
	tx := db.MustBegin()
	sql := "DELETE FROM capacity_plan_entries WHERE resource_id=$1 AND work_date=$2;"
	_, err := tx.Exec(sql, cpToSave[0].ResourceID, cpToSave[0].WorkDate)

	if err != nil {
		tx.Rollback()
		log.Println(err)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	sql = "INSERT INTO capacity_plan_entries (uuid, resource_id, project_id, work_date, work_percent, created_at) VALUES ($1,$2,$3,$4,$5,now()); "

	for _, capEntry := range cpToSave {
		_, err := tx.Exec(sql, capEntry.UUID.String(), capEntry.ResourceID, capEntry.ProjectID, capEntry.WorkDate, capEntry.WorkPercent)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (cp CapacityPlan) DeleteByID(dbx sqlx.DB, id int) error {
	op := "DeleteByID"

	sql := "DELETE FROM capacity_plan_entries WHERE id=$1;"
	_, err := dbx.Exec(sql, id)

	if err != nil {
		fmt.Println(op)
		log.Println(err)
		return err
	}
	return nil
}

func (cp CapacityPlan) ViewByUUID(dbx sqlx.DB, capEntryUUID uuid.UUID, orgID int) (model.CapacityPlanEntry, error) {
	op := "ViewByUUID"
	var capEntry model.CapacityPlanEntry
	sql := "SELECT cpe.*, proj.uuid as \"project.uuid\" FROM capacity_plan_entries cpe " +
		"LEFT JOIN projects proj on proj.id = cpe.project_id " +
		"WHERE cpe.uuid=$1 and cpe.deleted_at is null AND proj.organization_id=$2;"
	err := dbx.QueryRowx(sql, capEntryUUID.String(), orgID).StructScan(&capEntry)

	if err != nil {
		fmt.Println(op)
		log.Println(err)
		return model.CapacityPlanEntry{}, err
	}
	return capEntry, nil
}

func (cp CapacityPlan) ViewByID(dbx sqlx.DB, capEntryID int, orgID int) (model.CapacityPlanEntry, error) {
	op := "ViewByID"
	var capEntry model.CapacityPlanEntry
	sql := "SELECT cpe.*, proj.uuid as \"project.uuid\" FROM capacity_plan_entries cpe " +
		"LEFT JOIN projects proj on proj.id = cpe.project_id " +
		"WHERE cpe.id=$1 and cpe.deleted_at is null AND proj.organization_id=$2;"
	err := dbx.QueryRowx(sql, capEntryID, orgID).StructScan(&capEntry)

	if err != nil {
		fmt.Println(op)
		log.Println(err)
		return model.CapacityPlanEntry{}, err
	}
	return capEntry, nil
}

func (cp CapacityPlan) List(dbx sqlx.DB, orgID int, resourceID int, startDate time.Time, endDate time.Time) ([]model.CapacityPlanEntry, error) {
	op := "List"

	sql := "SELECT cpe.*, " +
		"proj.uuid as \"projectUUID\", " +
		"u.uuid as \"resourceUUID\" " +
		"FROM capacity_plan_entries cpe " +
		"LEFT JOIN projects proj on proj.id = cpe.project_id " +
		"LEFT JOIN users u on u.id = cpe.resource_id " +
		"WHERE cpe.resource_id=$1 " +
		"AND cpe.work_date>=$2 " +
		"AND cpe.work_date <=$3 " +
		"AND proj.organization_id =$4 " +
		"AND cpe.deleted_at is null;"
	var returnSl []model.CapacityPlanEntry
	rows, err := dbx.Queryx(sql, resourceID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), orgID)

	if err != nil {
		fmt.Println(op)
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var capEntry model.CapacityPlanEntry
		fmt.Println(rows.StructScan(&capEntry))
		returnSl = append(returnSl, capEntry)
	}
	if err != nil {
		fmt.Println(op)
		log.Println(err)
		return nil, err
	}
	return returnSl, nil
}
