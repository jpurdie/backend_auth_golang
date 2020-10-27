package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"
	"github.com/labstack/echo"
)

type Service interface {
	FetchByEmail(echo.Context, string) (authapi.User, error)
	FetchByExternalID(c echo.Context, externalID string) (authapi.User, error)
	ListRoles(echo.Context) ([]authapi.Role, error)
	List(c echo.Context, orgID int) ([]authapi.User, error)
	Update(c echo.Context, p authapi.Profile) error
	FetchProfile(c echo.Context, userID int, orgID int) (authapi.Profile, error)
}

// New creates new user application service
func New(db *pg.DB, pdb UserDB) User {
	return User{
		db:  db,
		udb: pdb,
	}
}

// Initialize initalizes user  service with defaults
func Initialize(db *pg.DB) User {
	return New(db, pgsql.User{})
}

type User struct {
	db  *pg.DB
	udb UserDB
}

type UserDB interface {
	List(db orm.DB, orgID int) ([]authapi.User, error)
	ListRoles(orm.DB) ([]authapi.Role, error)
	Update(orm.DB, authapi.Profile) error
	FetchByEmail(orm.DB, string) (authapi.User, error)
	FetchByExternalID(db orm.DB, externalID string) (authapi.User, error)
	//ListAuthorized(db orm.DB,u *authapi.User, includeInactive bool) ([]authapi.Profile, error)
	FetchProfile(db orm.DB, userID int, orgID int) (authapi.Profile, error)
	//Delete(db orm.DB, p authapi.Profile) error
}
