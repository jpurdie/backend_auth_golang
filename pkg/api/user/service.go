package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, authapi.User) (authapi.User, error)
	List(echo.Context, authapi.Pagination) ([]authapi.User, error)
	View(echo.Context, int) (authapi.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, Update) (authapi.User, error)
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.User{}, rbac, sec)
}

// User represents user application service
type User struct {
	db   *pg.DB
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(orm.DB, authapi.User) (authapi.User, error)
	View(orm.DB, int) (authapi.User, error)
	List(orm.DB, *authapi.ListQuery, authapi.Pagination) ([]authapi.User, error)
	Update(orm.DB, authapi.User) error
	Delete(orm.DB, authapi.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) authapi.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, authapi.AccessRole, int, int) error
	IsLowerRole(echo.Context, authapi.AccessRole) error
}
