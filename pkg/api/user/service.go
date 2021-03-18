package user

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Update(c echo.Context, userUUID uuid.UUID, orgID int, fieldsToUpdate map[string]string) error
	FetchUserByUUID(c echo.Context, userUUID uuid.UUID, orgID int) (*authapi.User, error)
	FetchUserByID(c echo.Context, userID int) (*authapi.User, error)
	FetchByEmail(echo.Context, string) (*authapi.User, error)
	FetchByExternalID(c echo.Context, externalID string) (*authapi.User, error)
	List(c echo.Context, orgID int) ([]authapi.User, error)
	ListRoles(echo.Context) ([]authapi.Role, error)
	UpdateRole(c echo.Context, level int, profileID int) error
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
	Update(db sqlx.DB, user authapi.User) error
	FetchUserByUUID(db sqlx.DB, userUUID uuid.UUID, orgID int) (*authapi.User, error)
	FetchByEmail(sqlx.DB, string) (*authapi.User, error)
	FetchByExternalID(db sqlx.DB, externalID string) (*authapi.User, error)
	FetchUserByID(db sqlx.DB, userID int) (*authapi.User, error)
	List(db sqlx.DB, orgID int) ([]authapi.User, error)
	UpdateRole(sqlx.DB, int, int) error
	ListRoles(sqlx.DB) ([]authapi.Role, error)
}
