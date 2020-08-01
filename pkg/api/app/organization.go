package app

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
)

// Organization defines database operations for Organization.
type OrganizationStore interface {
	//List() error
}

// Organization Resource implements account management handler.
type OrganizationResource struct {
	Store OrganizationStore
}

func NewOrganizationResource(store OrganizationStore) *OrganizationResource {
	return &OrganizationResource{
		Store: store,
	}
}
func (rs *OrganizationResource) router(r *echo.Group) {
	log.Println("Inside Organization Router")
	r.GET("/ping", rs.ping)
	//r.GET("", rs.listAuthorized)
}

func (rs *OrganizationResource) ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}
