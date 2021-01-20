package organization

import (
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/organization/platform/pgsql"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(echo.Context, authapi.Profile) error
}

// New creates new organization application service
func New(db *sqlx.DB, odb OrganizationDB) Organization {
	return Organization{
		db:  db,
		odb: odb,
	}
}

// Initialize initalizes profile application service with defaults
func Initialize(dbx *sqlx.DB) Organization {
	return New(dbx, pgsql.Organization{})
}

type Organization struct {
	db  *sqlx.DB
	odb OrganizationDB
}

type OrganizationDB interface {
	Create(sqlx.DB, authapi.Profile) error
}
