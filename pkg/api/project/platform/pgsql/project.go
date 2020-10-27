package pgsql

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
)

type Project struct {
}

func (p Project) UpdateName(db orm.DB, proj model.Project) error {
	op := "UpdateName"
	_, err := db.Model(&proj).
		Column("name").
		Where("project.uuid = ?", proj.UUID.String()).
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

func (p Project) UpdateDescr(db orm.DB, proj model.Project) error {
	op := "UpdateProject"

	_, err := db.Model(&proj).
		Column("description").
		Where("project.uuid = ?", proj.UUID.String()).
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

func (p Project) UpdateRGT(db orm.DB, pUUID uuid.UUID, rgt string) error {
	op := "UpdateRGT"
	tempProj := &model.Project{}
	tempProj.RGT = rgt
	_, err := db.Model(tempProj).
		Column("rgt").
		Where("project.uuid = ?", pUUID.String()).
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

func (p Project) UpdateSize(db orm.DB, pUUID uuid.UUID, ps model.ProjectSize) error {
	op := "UpdateSize"
	tempProj := &model.Project{}
	tempProj.SizeID = ps.ID
	_, err := db.Model(tempProj).
		Column("size_id").
		Where("project.uuid = ?", pUUID.String()).
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

func (p Project) UpdateStatus(db orm.DB, pUUID uuid.UUID, ps model.ProjectStatus) error {
	op := "UpdateStatus"
	tempProj := &model.Project{}
	tempProj.StatusID = ps.ID
	_, err := db.Model(tempProj).
		Column("status_id").
		Where("project.uuid = ?", pUUID.String()).
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

func (p Project) UpdateType(db orm.DB, pUUID uuid.UUID, pt model.ProjectType) error {
	op := "UpdateType"
	tempProj := &model.Project{}
	tempProj.TypeID = pt.ID
	_, err := db.Model(tempProj).
		Column("type_id").
		Where("project.uuid = ?", pUUID.String()).
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

func (p Project) UpdateComplexity(db orm.DB, pUUID uuid.UUID, pc model.ProjectComplexity) error {
	op := "UpdateComplexity"
	tempProj := &model.Project{}
	tempProj.ComplexityID = pc.ID
	_, err := db.Model(tempProj).
		Column("complexity_id").
		Where("project.uuid = ?", pUUID.String()).
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

func (p Project) Create(db orm.DB, newP model.Project) (model.Project, error) {
	op := "Create"
	_, err := db.Model(&newP).Insert()
	if err != nil {
		return model.Project{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return newP, nil
}

func (p Project) ListStatuses(db orm.DB) ([]model.ProjectStatus, error) {
	op := "ListStatuses"
	projectStatuses := make([]model.ProjectStatus, 0)
	err := db.Model(&projectStatuses).Select()
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return projectStatuses, nil
}

func (p Project) ListTypes(db orm.DB) ([]model.ProjectType, error) {
	op := "ListTypes"
	projectTypes := make([]model.ProjectType, 0)
	err := db.Model(&projectTypes).Select()
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return projectTypes, nil
}

func (p Project) ListComplexities(db orm.DB) ([]model.ProjectComplexity, error) {
	op := "ListComplexities"
	c := make([]model.ProjectComplexity, 0)
	err := db.Model(&c).Select()
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return c, nil
}

func (p Project) ListSizes(db orm.DB) ([]model.ProjectSize, error) {
	op := "ListSizes"
	s := make([]model.ProjectSize, 0)
	err := db.Model(&s).Select()
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return s, nil
}

func (p Project) List(db orm.DB, oID string) ([]model.Project, error) {
	op := "List"
	s := make([]model.Project, 0)
	err := db.Model(&s).
		Relation("Size").
		Relation("Complexity").
		Relation("Status").
		Relation("Type").
		Where("project.organization_id = ?", oID).
		Select()
	if err != nil {
		return nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return s, nil
}

func (p Project) View(db orm.DB, pUUID string) (model.Project, error) {
	op := "View"
	proj := model.Project{}
	err := db.Model(&proj).
		Relation("Size").
		Relation("Complexity").
		Relation("Status").
		Relation("Type").
		Where("project.uuid = ?", pUUID).
		Select()
	if err != nil {
		return model.Project{}, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return proj, nil
}
