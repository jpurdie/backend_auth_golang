package public

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi/pkg/api/database"
	"github.com/labstack/echo/v4"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
)

// API provides admin application resources and handlers.
type API struct {
	Pings *PingResource
}

// NewAPI configures and returns admin application API.
func NewAPI(db *pg.DB) (*API, error) {

	pingStore := database.NewPingStore(db)
	ping := NewPingResource(pingStore)

	api := &API{
		Pings: ping,
	}
	return api, nil
}
func (a *API) Router(r *echo.Group) {
	a.Pings.router(r)

}
