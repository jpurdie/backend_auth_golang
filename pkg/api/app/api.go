package app

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi/pkg/api/database"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo"
	"net/http"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
)

// API provides admin application resources and handlers.
type API struct {
	Organizations *OrganizationResource
}

// NewAPI configures and returns admin application API.
func NewAPI(db *pg.DB) (*API, error) {

	organizationStore := database.NewOrganizationStore(db)
	organization := NewOrganizationResource(organizationStore)

	api := &API{
		Organizations: organization,
	}
	return api, nil
}
func (a *API) Router(r *echo.Group) {

	organizations := r.Group("/organizations")
	a.Organizations.router(organizations)

	// Everything after here requires authentication
	//authMiddleware := authMw.Authenticate()
	//r.Use(authMiddleware)

	r.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}, authMw.Authorize([]string{"owner", "admin", "user"}))

}
