package capacityplan

import (
	"time"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/api/capacityplan/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/api/project"
	"github.com/jpurdie/authapi/pkg/api/user"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

type CapacityPlan struct {
	db     *sqlx.DB
	cpdb   CapacityPlanDB
	udb    user.UserDB
	projDB project.ProjectDB
}

// New creates new organization application service
func New(db *sqlx.DB, cpdb CapacityPlanDB, udb user.UserDB, projDB project.ProjectDB) CapacityPlan {
	return CapacityPlan{
		db:     db,
		cpdb:   cpdb,
		udb:    udb,
		projDB: projDB,
	}
}

type Service interface {
	Create(c echo.Context, orgID int, capPlans []model.CapacityPlanEntry) error
	List(c echo.Context, orgID int, resourceToGet uuid.UUID, startDate time.Time, endDate time.Time) ([]model.CapacityPlanEntry, error)
	ListSummary(c echo.Context, orgID int, resourceToGet uuid.UUID, startDate time.Time, endDate time.Time) ([]model.CapacityPlan, error)
	DeleteByID(c echo.Context, capEntryID int) error
	ViewByID(c echo.Context, capEntryID int, orgID int) (model.CapacityPlanEntry, error)
	ViewByUUID(c echo.Context, capEntryUUID uuid.UUID, orgID int) (model.CapacityPlanEntry, error)
}

func Initialize(dbx *sqlx.DB, usrDB user.UserDB, projDB project.ProjectDB) CapacityPlan {
	return New(dbx, pgsql.CapacityPlan{}, usrDB, projDB)
}

type CapacityPlanDB interface {
	Create(dbx sqlx.DB, orgID int, cp []model.CapacityPlanEntry) error
	List(dbx sqlx.DB, orgID int, resourceID int, startDate time.Time, endDate time.Time) ([]model.CapacityPlanEntry, error)
	DeleteByID(dbx sqlx.DB, capEntryID int) error
	ViewByID(dbx sqlx.DB, capEntryID int, orgID int) (model.CapacityPlanEntry, error)
	ViewByUUID(dbx sqlx.DB, capEntryUUID uuid.UUID, orgID int) (model.CapacityPlanEntry, error)
}
