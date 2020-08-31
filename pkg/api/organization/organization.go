package organization


import (
	"github.com/go-pg/pg/v9"
	"github.com/jpurdie/authapi"
	"net/http"
	"github.com/labstack/echo"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "incorrect old password")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

func (o Organization) Create(c echo.Context, org authapi.Profile) error {
	err := o.db.RunInTransaction(func (tx *pg.Tx) error{
		return o.odb.Create(tx, org)
	})
	return err;
}