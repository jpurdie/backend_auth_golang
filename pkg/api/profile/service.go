package profile

import (
	"github.com/go-pg/pg/v9"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/profile/platform/pgsql"
	"github.com/labstack/echo"
)

type Service interface {
	Create(echo.Context, authapi.Profile) error
}

// New creates new profile application service
func New(db *pg.DB, pdb ProfileDB) Profile {
	return Profile{
		db:   db,
		pdb:  pdb,
	}
}
// Initialize initalizes profile application service with defaults
func Initialize(db *pg.DB) Profile {
	return New(db, pgsql.Profile{})
}

type Profile struct {
	db   *pg.DB
	pdb  ProfileDB
}

type ProfileDB interface {
	Create(*pg.Tx, authapi.Profile) error
}