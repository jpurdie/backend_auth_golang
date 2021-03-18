package project

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/api/project/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/api/sponsorarea"
	"github.com/jpurdie/authapi/pkg/api/strategicalignment"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

type Project struct {
	db          *sqlx.DB
	pdb         ProjectDB
	saDB        strategicalignment.StrategicAlignmentDB
	sponsAreaDB sponsorarea.SponsorAreaDB
}

func New(db *sqlx.DB, pdb ProjectDB, saDB strategicalignment.StrategicAlignmentDB, sponsAreaDB sponsorarea.SponsorAreaDB) Project {
	return Project{
		db:          db,
		pdb:         pdb,
		saDB:        saDB,
		sponsAreaDB: sponsAreaDB,
	}
}

type Service interface {
	ListStatuses(c echo.Context) ([]model.ProjectStatus, error)
	ListTypes(c echo.Context) ([]model.ProjectType, error)
	ListComplexities(c echo.Context) ([]model.ProjectComplexity, error)
	ListSizes(c echo.Context) ([]model.ProjectSize, error)
	List(c echo.Context, oID int, filters map[string]string) ([]model.Project, error)
	View(c echo.Context, pUUID uuid.UUID) (model.Project, error)
	Create(c echo.Context, p model.Project) (uuid.UUID, error)
	Update(c echo.Context, orgID int, pUUID uuid.UUID, key string, val interface{}) error
}

func Initialize(dbx *sqlx.DB, saDB strategicalignment.StrategicAlignmentDB, sponsArea sponsorarea.SponsorAreaDB) Project {
	return New(dbx, pgsql.Project{}, saDB, sponsArea)
}

type ProjectDB interface {
	ListStatuses(dbx sqlx.DB) ([]model.ProjectStatus, error)
	ListTypes(dbx sqlx.DB) ([]model.ProjectType, error)
	ListComplexities(dbx sqlx.DB) ([]model.ProjectComplexity, error)
	ListSizes(dbx sqlx.DB) ([]model.ProjectSize, error)
	List(dbx sqlx.DB, oID int, filters map[string]string) ([]model.Project, error)
	View(dbx sqlx.DB, pUUID uuid.UUID) (model.Project, error)
	Create(dbx sqlx.DB, p model.Project) error
	Update(dbx sqlx.DB, orgID int, pUUID uuid.UUID, key string, val interface{}) error
	UpdateStrategicAlignments(dbx sqlx.DB, orgID int, projID int, strategicAlignmentIDs []int) error
	UpdateSponsorAreas(dbx sqlx.DB, orgID int, projID int, sponsorAreaIDs []int) error
}
