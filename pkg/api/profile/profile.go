package profile

import (
	"net/http"

	"github.com/jpurdie/authapi"

	"github.com/labstack/echo/v4"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "incorrect old password")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

// Change changes user's password
func (p Profile) Create(c echo.Context, profile authapi.Profile) error {
	return p.pdb.Create(*p.db, profile)
}

func (p Profile) FetchProfileByExternalID(c echo.Context, externalID string) (authapi.Profile, error) {
	return p.pdb.FetchProfileByExternalID(*p.db, externalID)
}
