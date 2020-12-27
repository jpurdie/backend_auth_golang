package user

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo"
	"net/http"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "incorrect old password")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

func (u User) FetchByEmail(c echo.Context, email string) (authapi.User, error) {
	return u.udb.FetchByEmail(*u.db, email)
}

func (u User) ListRoles(c echo.Context) ([]authapi.Role, error) {
	return u.udb.ListRoles(*u.db)
}
func (u User) FetchByExternalID(c echo.Context, externalID string) (authapi.User, error) {
	return u.udb.FetchByExternalID(*u.db, externalID)
}

func (u User) List(c echo.Context, orgID uint) ([]authapi.User, error) {
	return u.udb.List(*u.db, orgID)
}

func (u User) UpdateRole(c echo.Context, level int, profileID uint) error {
	return u.udb.UpdateRole(*u.db, level, profileID)
}

func (u User) FetchProfile(c echo.Context, userID int, orgID int) (authapi.Profile, error) {
	return u.udb.FetchProfile(*u.db, userID, orgID)
}
func (u User) FetchUserByUUID(c echo.Context, userUUID uuid.UUID, orgID uint) (authapi.User, error) {
	return u.udb.FetchUserByUUID(*u.db, userUUID, orgID)
}
