package ping

import (
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo/v4"
)

// Custom errors
var ()

// Change changes user's password
func (p Ping) Create(c echo.Context, ping authapi.Ping) error {
	return p.pdb.Create(*p.db, ping)
}
