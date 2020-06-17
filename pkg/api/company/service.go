package company

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/company/platform/pgsql"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, authapi.CompanyUser) (authapi.CompanyUser, error)
	//List(echo.Context, authapi.Pagination) ([]authapi.Company, error)
	//View(echo.Context, int) (authapi.Company, error)
	//Delete(echo.Context, int) error
	//Update(echo.Context, Update) (authapi.Company, error)
}

// New creates new user application service
func New(db *pg.DB, cdb CDB, rbac RBAC, sec Securer) *CompanyUser {
	return &CompanyUser{db: db, cdb: cdb, rbac: rbac, sec: sec}
}

// Initialize initalizes Company application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *CompanyUser {
	return New(db, pgsql.CompanyUser{}, rbac, sec)
}

// User represents company application service
type CompanyUser struct {
	db   *pg.DB
	cdb  CDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// CDB represents Company repository interface
type CDB interface {
	Create(*pg.DB, authapi.CompanyUser) (authapi.CompanyUser, error)
	//View(orm.DB, int) (authapi.Company, error)
	//List(orm.DB, *authapi.ListQuery, authapi.Pagination) ([]authapi.Company, error)
	//Update(orm.DB, authapi.Company) error
	//Delete(orm.DB, authapi.Company) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) authapi.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, authapi.AccessRole, int, int) error
	IsLowerRole(echo.Context, authapi.AccessRole) error
}
