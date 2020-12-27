package user

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"
	"github.com/labstack/echo"
)

type Service interface {
	FetchByEmail(echo.Context, string) (authapi.User, error)
	FetchByExternalID(c echo.Context, externalID string) (authapi.User, error)
	ListRoles(echo.Context) ([]authapi.Role, error)
	List(c echo.Context, orgID uint) ([]authapi.User, error)
	UpdateRole(c echo.Context, level int, profileID uint) error
	FetchProfile(c echo.Context, userID int, orgID int) (authapi.Profile, error)
	FetchUserByUUID(c echo.Context, userUUID uuid.UUID, orgID uint) (authapi.User, error)
}

// New creates new user application service
func New(db *sqlx.DB, pdb UserDB) User {
	return User{
		db:  db,
		udb: pdb,
	}
}

// Initialize initalizes user  service with defaults
func Initialize(db *sqlx.DB) User {
	return New(db, pgsql.User{})
}

type User struct {
	db  *sqlx.DB
	udb UserDB
}

type UserDB interface {
	List(db sqlx.DB, orgID uint) ([]authapi.User, error)
	ListRoles(sqlx.DB) ([]authapi.Role, error)
	UpdateRole(sqlx.DB, int, uint) error
	FetchByEmail(sqlx.DB, string) (authapi.User, error)
	FetchByExternalID(db sqlx.DB, externalID string) (authapi.User, error)
	//ListAuthorized(db orm.DB,u *authapi.User, includeInactive bool) ([]authapi.Profile, error)
	FetchProfile(db sqlx.DB, userID int, orgID int) (authapi.Profile, error)
	FetchUserByUUID(db sqlx.DB, userUUID uuid.UUID, orgID uint) (authapi.User, error)
	//Delete(db orm.DB, p authapi.Profile) error
}
