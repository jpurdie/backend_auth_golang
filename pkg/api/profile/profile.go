package profile

import (
	"github.com/jpurdie/authapi"
	"net/http"

	"github.com/labstack/echo"
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