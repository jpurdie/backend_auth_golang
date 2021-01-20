package organization

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

func (o Organization) Create(c echo.Context, org authapi.Profile) error {
	return o.odb.Create(*o.db, org)
}
