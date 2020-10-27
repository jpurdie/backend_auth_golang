package project

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi/pkg/api/project/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo"
)

type Project struct {
	db  *pg.DB
	pdb ProjectDB
}

func New(db *pg.DB, pdb ProjectDB) Project {
	return Project{
		db:  db,
		pdb: pdb,
	}
}

type Service interface {
	ListStatuses(c echo.Context) ([]model.ProjectStatus, error)
	ListTypes(c echo.Context) ([]model.ProjectType, error)
	ListComplexities(c echo.Context) ([]model.ProjectComplexity, error)
	ListSizes(c echo.Context) ([]model.ProjectSize, error)
	List(c echo.Context, oID string) ([]model.Project, error)
	View(c echo.Context, pUUID string) (model.Project, error)
	Create(c echo.Context, p model.Project) (model.Project, error)
	UpdateComplexity(c echo.Context, pUUID uuid.UUID, pc model.ProjectComplexity) error
	UpdateType(c echo.Context, pUUID uuid.UUID, pc model.ProjectType) error
	UpdateStatus(c echo.Context, pUUID uuid.UUID, ps model.ProjectStatus) error
	UpdateSize(c echo.Context, pUUID uuid.UUID, ps model.ProjectSize) error
	UpdateRGT(c echo.Context, pUUID uuid.UUID, rgt string) error
	UpdateDescr(c echo.Context, proj model.Project) error
	UpdateName(c echo.Context, proj model.Project) error
}

func Initialize(db *pg.DB) Project {
	return New(db, pgsql.Project{})
}

type ProjectDB interface {
	ListStatuses(db orm.DB) ([]model.ProjectStatus, error)
	ListTypes(db orm.DB) ([]model.ProjectType, error)
	ListComplexities(db orm.DB) ([]model.ProjectComplexity, error)
	ListSizes(db orm.DB) ([]model.ProjectSize, error)
	List(db orm.DB, oID string) ([]model.Project, error)
	View(db orm.DB, pUUID string) (model.Project, error)
	Create(db orm.DB, p model.Project) (model.Project, error)
	UpdateComplexity(db orm.DB, pUUID uuid.UUID, pc model.ProjectComplexity) error
	UpdateType(db orm.DB, pUUID uuid.UUID, pc model.ProjectType) error
	UpdateStatus(db orm.DB, pUUID uuid.UUID, ps model.ProjectStatus) error
	UpdateSize(db orm.DB, pUUID uuid.UUID, ps model.ProjectSize) error
	UpdateRGT(db orm.DB, pUUID uuid.UUID, rgt string) error
	UpdateDescr(db orm.DB, proj model.Project) error
	UpdateName(db orm.DB, proj model.Project) error
}
