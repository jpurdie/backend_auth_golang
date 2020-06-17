package ping

import (
	"github.com/labstack/echo"
)

func (p Ping) Create(c echo.Context, req int) (string, error) {

	return "pong", nil
}
