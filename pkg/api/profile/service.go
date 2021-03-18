package profile

import (
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/profile/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/api/user"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(echo.Context, authapi.Profile) error
	FetchProfileByExternalID(c echo.Context, externalID string) (authapi.Profile, error)
}

// New creates new profile application service
func New(db *sqlx.DB, pdb ProfileDB, usrDB user.UserDB) Profile {
	return Profile{
		db:     db,
		pdb:    pdb,
		userDB: usrDB,
	}
}

// Initialize initalizes profile application service with defaults
func Initialize(db *sqlx.DB, userDB user.UserDB) Profile {
	return New(db, pgsql.Profile{}, userDB)
}

type Profile struct {
	db     *sqlx.DB
	pdb    ProfileDB
	userDB user.UserDB
}

type ProfileDB interface {
	Create(sqlx.DB, authapi.Profile) error
	FetchProfileByExternalID(db sqlx.DB, externalID string) (authapi.Profile, error)
}
