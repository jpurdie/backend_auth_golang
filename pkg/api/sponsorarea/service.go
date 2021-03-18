package sponsorarea

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/api/sponsorarea/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

type SponsorArea struct {
	db   *sqlx.DB
	sadb SponsorAreaDB
}

func New(db *sqlx.DB, sadb SponsorAreaDB) SponsorArea {
	return SponsorArea{
		db:   db,
		sadb: sadb,
	}
}

type Service interface {
	Create(c echo.Context, saName string, orgID int) (model.SponsorArea, error)
	List(c echo.Context, orgID int) ([]model.SponsorArea, error)
	Delete(c echo.Context, saUUID uuid.UUID, oID int) error
	Update(c echo.Context, mySA model.SponsorArea) error
}

func Initialize(db *sqlx.DB) SponsorArea {
	return New(db, pgsql.SponsorArea{})
}

type SponsorAreaDB interface {
	Create(dbx sqlx.DB, sa model.SponsorArea) (model.SponsorArea, error)
	List(dbx sqlx.DB, orgID int) ([]model.SponsorArea, error)
	Update(dbx sqlx.DB, mySA model.SponsorArea) error
	Delete(dbx sqlx.DB, saUUID uuid.UUID, oID int) error
}
