package profile

import (
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/profile/platform/pgsql"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(echo.Context, authapi.Profile) error
}

// New creates new profile application service
func New(db *sqlx.DB, pdb ProfileDB) Profile {
	return Profile{
		db:  db,
		pdb: pdb,
	}
}

// Initialize initalizes profile application service with defaults
func Initialize(db *sqlx.DB) Profile {
	return New(db, pgsql.Profile{})
}

type Profile struct {
	db  *sqlx.DB
	pdb ProfileDB
}

type ProfileDB interface {
	Create(sqlx.DB, authapi.Profile) error
}
