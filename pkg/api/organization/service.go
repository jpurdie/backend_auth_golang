package organization

import (
	"github.com/go-pg/pg/v9"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/organization/platform/pgsql"
	"github.com/labstack/echo"
)

type Service interface {
	Create(echo.Context, authapi.Profile) error
}

// New creates new organization application service
func New(db *pg.DB, odb OrganizationDB) Organization {
	return Organization{
		db:   db,
		odb:  odb,
	}
}
// Initialize initalizes profile application service with defaults
func Initialize(db *pg.DB) Organization {
	return New(db, pgsql.Organization{})
}

type Organization struct {
	db   *pg.DB
	odb  OrganizationDB
}

type OrganizationDB interface {
	Create(*pg.Tx, authapi.Profile) error
}