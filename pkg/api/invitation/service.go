package invitation

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/invitation/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/api/user"
	"github.com/labstack/echo"
)

type Invitation struct {
	db   *pg.DB
	idb  InvitationDB
	udb user.UserDB
}

// New creates new organization application service
func New(db *pg.DB, idb InvitationDB, udb user.UserDB) Invitation {
	return Invitation{
		db:   db,
		idb:  idb,
		udb:  udb,
	}
}

type Service interface {
	Create(c echo.Context, invite authapi.Invitation) error
	List(c echo.Context, o *authapi.Organization, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error)
	Delete(c echo.Context, invite authapi.Invitation) error
	View(c echo.Context, tokenHash string) (authapi.Invitation, error)
	CreateUser(c echo.Context, cu authapi.Profile, i authapi.Invitation) error
}

// Initialize initalizes profile application service with defaults
func Initialize(db *pg.DB) Invitation {
	return New(db, pgsql.Invitation{},  )
}


type InvitationDB interface {
	Create(db orm.DB, invite authapi.Invitation) error
	Delete(db orm.DB, invite authapi.Invitation) error
	View(db orm.DB, tokenHash string) (authapi.Invitation, error)
	List(db orm.DB, o *authapi.Organization, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error)
	CreateUser(tx *pg.Tx, cu authapi.Profile, i authapi.Invitation) error
	FindUserByEmail(db orm.DB, email string, orgID int) (authapi.Invitation, error)
}
