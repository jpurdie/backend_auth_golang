package transport

import (
	"github.com/jpurdie/authapi/pkg/api/ping"
	"github.com/labstack/echo"
	"net/http"
)

type HTTP struct {
	svc ping.Service
}

func NewHTTP(svc ping.Service, r *echo.Group) {
	h := HTTP{svc}

	r.GET("/ping", h.ping)
	r.GET("/secureping", h.securePing)
}
func (h HTTP) ping(c echo.Context) error {

	return c.JSON(http.StatusOK, "pong")
}

func (h HTTP) securePing(c echo.Context) error {

	return c.JSON(http.StatusOK, "securepong")
}
