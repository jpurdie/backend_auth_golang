package pgsql

import (
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
)

type Project struct {
}

func (p Project) UpdateSponsorAreas(db sqlx.DB, orgID int, projID int, sponsorAreaIDs []int) error {
	op := "UpdateSponsorAreas"

	tx, err := db.Beginx()
	sql := "DELETE FROM project_sponsor_areas where project_id =$1;"
	_, err = tx.Exec(sql, projID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for _, tempSA := range sponsorAreaIDs {
		sql = "INSERT INTO project_sponsor_areas (project_id, sponsor_area_id) VALUES ($1, $2);"
		_, err = tx.Exec(sql, projID, tempSA)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("There was a transaction error")
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

func (p Project) UpdateStrategicAlignments(db sqlx.DB, orgID int, projID int, strategicAlignmentIDs []int) error {
	op := "UpdateStrategicAlignments"

	tx, err := db.Beginx()
	sql := "DELETE FROM project_strategic_alignments where project_id =$1;"
	_, err = tx.Exec(sql, projID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for _, tempSA := range strategicAlignmentIDs {
		sql = "INSERT INTO project_strategic_alignments (project_id, strategic_alignment_id) VALUES ($1, $2);"
		_, err = tx.Exec(sql, projID, tempSA)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("There was a transaction error")
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

func (p Project) Update(db sqlx.DB, orgID int, projID uuid.UUID, key string, val interface{}) error {
	op := "Update"
	sql := ""
	switch key {
	case "openForTimeEntry":
		sql = "UPDATE projects SET open_for_time_entry=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "timeConstrained":
		sql = "UPDATE projects SET time_constrained=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "compliance":
		sql = "UPDATE projects SET compliance=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "rgt":
		sql = "UPDATE projects SET rgt=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "status":
		sql = "UPDATE projects SET status_id=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "type":
		sql = "UPDATE projects SET type_id=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "complexity":
		sql = "UPDATE projects SET complexity_id=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "size":
		sql = "UPDATE projects SET size_id=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "name":
		sql = "UPDATE projects SET name=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	case "description":
		sql = "UPDATE projects SET description=$1, updated_at=now() WHERE uuid=$2 AND organization_id=$3 and deleted_at is null;"
	default:
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINVALID,
		}
	}
	_, err := db.Exec(sql, val, projID, orgID)
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return nil
}

func (p Project) Create(db sqlx.DB, newP model.Project) error {
	op := "Create"

	tx, err := db.Beginx()
	sql := "INSERT INTO projects " +
		"(created_at, " +
		"uuid, " +
		"name, " +
		"description, " +
		"rgt, " +
		"status_id, " +
		"type_id, " +
		"complexity_id, " +
		"size_id, " +
		"organization_id, " +
		"open_for_time_entry, " +
		"compliance, " +
		"time_constrained " +
		") " +
		" VALUES (now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id as xID;"
	projID := 0
	err = tx.QueryRowx(sql, newP.UUID, newP.Name, newP.Description, newP.RGT, newP.StatusID, newP.TypeID, newP.ComplexityID, newP.SizeID, newP.OrganizationID, newP.OpenForTimeEntry, newP.Compliance, newP.TimeConstrained).Scan(&projID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	for _, tempSA := range newP.StrategicAlignments {
		sql = "INSERT INTO project_strategic_alignments (project_id, strategic_alignment_id) VALUES ($1, $2);"
		_, err = tx.Exec(sql, projID, tempSA.ID)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
	}

	for _, tempSA := range newP.SponsorAreas {
		sql = "INSERT INTO project_sponsor_areas (project_id, sponsor_area_id) VALUES ($1, $2);"
		_, err = tx.Exec(sql, projID, tempSA.ID)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return &authapi.Error{
				Op:   op,
				Code: authapi.EINTERNAL,
				Err:  err,
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("There was a transaction error")
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

func (p Project) ListStatuses(db sqlx.DB) ([]model.ProjectStatus, error) {
	op := "ListStatuses"
	var projectStatuses []model.ProjectStatus
	err := db.Select(&projectStatuses, `SELECT id, name, uuid, created_at, updated_at, deleted_at FROM project_statuses WHERE deleted_at is null;`)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return projectStatuses, nil
}

func (p Project) ListTypes(db sqlx.DB) ([]model.ProjectType, error) {
	op := "ListTypes"
	var projectTypes []model.ProjectType
	err := db.Select(&projectTypes, `SELECT id, name, uuid, created_at, updated_at, deleted_at FROM project_types WHERE deleted_at is null`)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return projectTypes, nil
}

func (p Project) ListComplexities(db sqlx.DB) ([]model.ProjectComplexity, error) {
	op := "ListComplexities"
	var c []model.ProjectComplexity
	err := db.Select(&c, `SELECT id, name, uuid, created_at, updated_at, deleted_at FROM project_complexities WHERE deleted_at is null`)
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return c, nil
}

func (p Project) ListSizes(db sqlx.DB) ([]model.ProjectSize, error) {
	op := "ListSizes"
	var projectSizes []model.ProjectSize
	err := db.Select(&projectSizes, "SELECT id, name, uuid, created_at, updated_at, deleted_at FROM project_sizes WHERE deleted_at is null;")
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return projectSizes, nil
}

func (p Project) List(db sqlx.DB, oID int, filters map[string]string) ([]model.Project, error) {
	op := "List"
	var projects []model.Project

	filterQuery := ""
	for key, val := range filters {
		switch strings.ToLower(key) {
		case "openfortimeentry":
			if val == "true" {
				filterQuery += "AND p.open_for_time_entry=TRUE "
			} else {
				filterQuery += "AND p.open_for_time_entry=FALSE "
			}
			break
		}
	}

	query := "SELECT p.*, " +
		"pstat.id AS \"project_status.id\", " +
		"pstat.name AS \"project_status.name\", " +
		"pstat.uuid AS \"project_status.uuid\", " +
		"psize.id AS \"project_size.id\", " +
		"psize.name AS \"project_size.name\", " +
		"psize.uuid AS \"project_size.uuid\", " +
		"pcompl.id AS \"project_complexity.id\", " +
		"pcompl.name AS \"project_complexity.name\", " +
		"pcompl.uuid AS \"project_complexity.uuid\", " +
		"ptype.id AS \"project_type.id\", " +
		"ptype.name AS \"project_type.name\", " +
		"ptype.uuid AS \"project_type.uuid\" " +
		"FROM projects p " +
		"JOIN project_sizes psize on p.size_id = psize.id " +
		"JOIN project_complexities pcompl on pcompl.id = p.complexity_id " +
		"JOIN project_statuses pstat on pstat.id = p.status_id " +
		"JOIN project_types ptype on ptype.id = p.type_id " +
		"WHERE p.organization_id=$1 " +
		filterQuery +
		"AND p.deleted_at is null"

	err := db.Select(&projects, query, oID)

	//for _, tempProject := range projects {
	//	tempProject.Size =
	//}

	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	return projects, nil
}

func (p Project) View(db sqlx.DB, pUUID uuid.UUID) (model.Project, error) {
	op := "View"
	var project model.Project

	query := "SELECT p.*, " +
		"pstat.id AS \"project_status.id\", " +
		"pstat.name AS \"project_status.name\", " +
		"pstat.uuid AS \"project_status.uuid\", " +
		"psize.id AS \"project_size.id\", " +
		"psize.name AS \"project_size.name\", " +
		"psize.uuid AS \"project_size.uuid\", " +
		"pcompl.id AS \"project_complexity.id\", " +
		"pcompl.name AS \"project_complexity.name\", " +
		"pcompl.uuid AS \"project_complexity.uuid\", " +
		"ptype.id AS \"project_type.id\", " +
		"ptype.name AS \"project_type.name\", " +
		"ptype.uuid AS \"project_type.uuid\" " +
		"FROM projects p " +
		"JOIN project_sizes psize on p.size_id = psize.id " +
		"JOIN project_complexities pcompl on pcompl.id = p.complexity_id " +
		"JOIN project_statuses pstat on pstat.id = p.status_id " +
		"JOIN project_types ptype on ptype.id = p.type_id " +
		"WHERE p.uuid=$1 " +
		"AND p.deleted_at is null;"

	err := db.QueryRowx(query, pUUID.String()).StructScan(&project)

	if err != nil {
		return model.Project{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}

	var strategicAlignments []*model.StrategicAlignment
	query = "SELECT " +
		"sa.* " +
		"FROM project_strategic_alignments psa " +
		"JOIN strategic_alignments sa on sa.id = psa.strategic_alignment_id " +
		"where psa.project_id=$1;"
	err = db.Select(&strategicAlignments, query, project.ID)
	if err != nil {
		return model.Project{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	project.StrategicAlignments = strategicAlignments

	var sponsorAreas []*model.SponsorArea
	query = "SELECT " +
		"sa.* " +
		"FROM project_sponsor_areas psa " +
		"JOIN sponsor_areas sa on sa.id = psa.sponsor_area_id " +
		"where psa.project_id=$1;"
	err = db.Select(&sponsorAreas, query, project.ID)
	if err != nil {
		return model.Project{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	project.SponsorAreas = sponsorAreas

	return project, nil
}
