package public

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PingStore defines database operations for Ping.
type PingStore interface {
	Ping() error
	//List() error
}

// Ping Resource implements account management handler.
type PingResource struct {
	Store PingStore
}

func NewPingResource(store PingStore) *PingResource {
	return &PingResource{
		Store: store,
	}
}
func (rs *PingResource) router(r *echo.Group) {
	r.GET("/ping", rs.ping)

}

func (rs *PingResource) ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
