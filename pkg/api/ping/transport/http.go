package transport

import (
	"github.com/labstack/echo"
	"github.com/jpurdie/authapi/pkg/api/ping"
	"net/http"
	"strconv"
)

type HTTP struct {
	svc ping.Service
}

func NewHTTP(svc ping.Service, r *echo.Group) {
	h := HTTP{svc}
	r.GET("/ping", h.create)
}
func (h HTTP) create(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	id = 2

	pong, err := h.svc.Create(c, id)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pong)
}
