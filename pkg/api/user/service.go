package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"
	"github.com/labstack/echo"
)

type Service interface {
	Fetch(echo.Context, authapi.User) (authapi.User, error)
	ListRoles(echo.Context) ([]authapi.Role, error)
	List(c echo.Context, orgID int) ([]authapi.User, error)
	Update(c echo.Context, p authapi.Profile) error
	FetchProfile(c echo.Context, u authapi.User, o authapi.Organization) (authapi.Profile, error)
}

// New creates new password application service
func New(db *pg.DB, pdb UserDB) User {
	return User{
		db:  db,
		udb: pdb,
	}
}

// Initialize initalizes profile application service with defaults
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
	Update(orm.DB,authapi.Profile) error
	Fetch(orm.DB, authapi.User) (authapi.User, error)
	//ListAuthorized(db orm.DB,u *authapi.User, includeInactive bool) ([]authapi.Profile, error)
	FetchProfile(db orm.DB, u authapi.User, o authapi.Organization) (authapi.Profile, error)
	//Delete(db orm.DB, p authapi.Profile) error
}
