package invitation

import (
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/invitation/platform/pgsql"
	"github.com/jpurdie/authapi/pkg/api/user"
	"github.com/labstack/echo/v4"
)

type Invitation struct {
	db  *sqlx.DB
	idb InvitationDB
	udb user.UserDB
}

// New creates new organization application service
func New(db *sqlx.DB, idb InvitationDB, udb user.UserDB) Invitation {
	return Invitation{
		db:  db,
		idb: idb,
		udb: udb,
	}
}

type Service interface {
	Create(c echo.Context, invite authapi.Invitation) error
	List(c echo.Context, orgID int, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error)
	Delete(c echo.Context, email string, orgID int) error
	View(c echo.Context, tokenPlainTextString string) (authapi.Invitation, error)
	CreateUser(c echo.Context, cu authapi.Profile, i authapi.Invitation) error
}

// Initialize initalizes profile application service with defaults
func Initialize(dbx *sqlx.DB, usr user.UserDB) Invitation {
	return New(dbx, pgsql.Invitation{}, usr)
}

type InvitationDB interface {
	Create(dbx sqlx.DB, invite authapi.Invitation) error
	Delete(dbx sqlx.DB, email string, orgID int) error
	ViewByEmail(dbx sqlx.DB, email string, orgID int) (authapi.Invitation, error)
	View(dbx sqlx.DB, tokenPlainTextString string) (authapi.Invitation, error)
	List(dbx sqlx.DB, orgID int, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error)
	CreateUser(dbx sqlx.DB, cu authapi.Profile, i authapi.Invitation) error
}
