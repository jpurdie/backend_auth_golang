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
	AuthOrganizations *AuthOrganizationResource
	AuthInvitations   *AuthInvitationResource
	Organizations     *OrganizationResource
	Invitations       *InvitationResource
	Authorizations    *AuthorizationResource
}

// NewAPI configures and returns admin application API.
func NewAPI(db *pg.DB) (*API, error) {

	authOrganizationStore := database.NewAuthOrganizationStore(db)
	authOrganization := NewAuthOrganizationResource(authOrganizationStore)

	authInvitationStore := database.NewAuthInvitationStore(db)
	authInvitation := NewAuthInvitationResource(authInvitationStore)

	organizationStore := database.NewOrganizationStore(db)
	organization := NewOrganizationResource(organizationStore)

	invitationStore := database.NewInvitationStore(db)
	invitation := NewInvitationResource(invitationStore)

	authorizationStore := database.NewAuthorizationStore(db)
	authorization := NewAuthorizationResource(authorizationStore)

	api := &API{
		Organizations:     organization,
		AuthOrganizations: authOrganization,
		AuthInvitations:   authInvitation,
		Invitations:       invitation,
		Authorizations:    authorization,
	}
	return api, nil
}
func (a *API) Router(r *echo.Group) {

	authOrganizations := r.Group("/auth/organizations")
	a.AuthOrganizations.router(authOrganizations)

	authInvitations := r.Group("/auth/invitations")
	a.AuthInvitations.router(authInvitations)

	authMiddleware := authMw.Authenticate()
	r.Use(authMiddleware)

	organizations := r.Group("/organizations")
	a.Organizations.router(organizations)

	invitations := r.Group("/invitations")
	a.Invitations.router(invitations)

	// Everything after here requires authentication

	r.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}, authMw.CheckAuthorization([]string{"owner", "admin", "user"}))

}
