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
	AuthProfiles    *AuthProfileResource
	AuthInvitations *AuthInvitationResource
	Organizations   *OrganizationResource
	Invitations     *InvitationResource
	Authorizations  *AuthorizationResource
	Users           *UserResource
}

// NewAPI configures and returns admin application API.
func NewAPI(db *pg.DB) (*API, error) {

	authProfileStore := database.NewAuthProfileStore(db)
	authProfile := NewAuthProfileResource(authProfileStore)

	authInvitationStore := database.NewAuthInvitationStore(db)
	authInvitation := NewAuthInvitationResource(authInvitationStore)

	organizationStore := database.NewOrganizationStore(db)
	organization := NewOrganizationResource(organizationStore)

	invitationStore := database.NewInvitationStore(db)
	invitation := NewInvitationResource(invitationStore)

	authorizationStore := database.NewAuthorizationStore(db)
	authorization := NewAuthorizationResource(authorizationStore)

	userStore := database.NewUserStore(db)
	user := NewUserResource(userStore)

	api := &API{
		Organizations:   organization,
		AuthProfiles:    authProfile,
		AuthInvitations: authInvitation,
		Invitations:     invitation,
		Users:           user,
		Authorizations:  authorization,
	}
	return api, nil
}
func (a *API) Router(r *echo.Group) {

	r.GET("/unauthping", func(c echo.Context) error {
		return c.String(http.StatusOK, "unauthpong")
	})

	authProfiles := r.Group("/auth/organizations")
	a.AuthProfiles.router(authProfiles)

	authInvitations := r.Group("/auth/invitations")
	a.AuthInvitations.router(authInvitations)

	// Everything after here requires authentication
	authMiddleware := authMw.Authenticate()
	r.Use(authMiddleware)

	organizations := r.Group("/organizations")
	a.Organizations.router(organizations)

	invitations := r.Group("/invitations")
	a.Invitations.router(invitations)

	users := r.Group("/users")
	a.Users.router(users)

	r.GET("/authping", func(c echo.Context) error {
		return c.String(http.StatusOK, "authpong")
	}, authMw.CheckAuthorization([]string{"owner", "admin", "user"}))

}
