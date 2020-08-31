package user




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

func (u User) Fetch(c echo.Context, us authapi.User) (authapi.User, error) {
	return u.udb.Fetch(u.db, us)
}

func (u User) ListRoles(c echo.Context) ([]authapi.Role, error) {
	return u.udb.ListRoles(u.db)
}

func (u User) List(c echo.Context, orgID int) ([]authapi.User, error) {
	return u.udb.List(u.db, orgID)
}

func (u User) Update(c echo.Context, p authapi.Profile) error {
	return u.udb.Update(u.db, p)
}

func (u User) FetchProfile(c echo.Context, us authapi.User, o authapi.Organization) (authapi.Profile, error) {
	return u.udb.FetchProfile(u.db, us, o)
}

