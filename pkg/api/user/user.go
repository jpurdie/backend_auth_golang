package user

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo/v4"
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

func (u User) List(c echo.Context, orgID int) ([]authapi.User, error) {
	return u.udb.List(*u.db, orgID)
}

func (u User) UpdateRole(c echo.Context, level int, profileID int) error {
	return u.udb.UpdateRole(*u.db, level, profileID)
}

func (u User) FetchProfile(c echo.Context, userID int, orgID int) (authapi.Profile, error) {
	return u.udb.FetchProfile(*u.db, userID, orgID)
}
func (u User) FetchUserByUUID(c echo.Context, userUUID uuid.UUID, orgID int) (authapi.User, error) {
	return u.udb.FetchUserByUUID(*u.db, userUUID, orgID)
}

func (u User) FetchUserByID(c echo.Context, userID int) (authapi.User, error) {
	return u.udb.FetchUserByID(*u.db, userID)
}
func (u User) Update(c echo.Context, userUUID uuid.UUID, orgID int, fieldsToUpdate map[string]string) error {
	op := "Update"
	currUser, err := u.udb.FetchUserByUUID(*u.db, userUUID, orgID)
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
		}
	}
	for key, val := range fieldsToUpdate {
		switch key {
		case "/firstName":
			currUser.FirstName = val
		case "/lastName":
			currUser.LastName = val
		case "/timeZone":
			currUser.TimeZone = &val
		}
	}
	log.Print(currUser)

	return u.udb.Update(*u.db, currUser)
}
