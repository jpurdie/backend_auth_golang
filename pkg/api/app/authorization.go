package app

import (
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type AuthorizationStore interface {
	CheckAuthorization(u *authapi.User) error
}

// Invitation Resource implements account management handler.
type AuthorizationResource struct {
	Store AuthorizationStore
}

func NewAuthorizationResource(store AuthorizationStore) *AuthorizationResource {
	return &AuthorizationResource{
		Store: store,
	}
}

func (rs *AuthorizationResource) checkAuthorization(c echo.Context) error {
	log.Println("Check database here for auths")
	return c.JSON(http.StatusUnauthorized, "")
}
