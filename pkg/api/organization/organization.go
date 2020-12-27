package organization

import (
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo"
	"net/http"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "incorrect old password")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

func (o Organization) Create(c echo.Context, org authapi.Profile) error {
	return o.odb.Create(*o.db, org)
}