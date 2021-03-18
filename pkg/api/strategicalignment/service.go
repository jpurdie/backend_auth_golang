package strategicalignment

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/api/strategicalignment/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

type StrategicAlignment struct {
	db   *sqlx.DB
	sadb StrategicAlignmentDB
}

func New(db *sqlx.DB, sadb StrategicAlignmentDB) StrategicAlignment {
	return StrategicAlignment{
		db:   db,
		sadb: sadb,
	}
}

type Service interface {
	Create(c echo.Context, alignmentName string, orgID int) (model.StrategicAlignment, error)
	List(c echo.Context, oID int) ([]model.StrategicAlignment, error)
	Delete(c echo.Context, saUUID uuid.UUID, oID int) error
	Update(c echo.Context, mySA model.StrategicAlignment) error
}

func Initialize(db *sqlx.DB) StrategicAlignment {
	return New(db, pgsql.StrategicAlignment{})
}

type StrategicAlignmentDB interface {
	Create(dbx sqlx.DB, sa model.StrategicAlignment) (model.StrategicAlignment, error)
	List(dbx sqlx.DB, oID int) ([]model.StrategicAlignment, error)
	Update(dbx sqlx.DB, mySA model.StrategicAlignment) error
	Delete(dbx sqlx.DB, saUUID uuid.UUID, oID int) error
}
