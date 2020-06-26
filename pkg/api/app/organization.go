package app

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	AuthUtil "github.com/jpurdie/authapi/pkg/utl/Auth"
	"github.com/jpurdie/authapi/pkg/utl/Auth0"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

// Organization defines database operations for Organization.
type OrganizationStore interface {
	Create(authapi.OrganizationUser) error
	ListAccessible(u *authapi.User, includeInactive bool) ([]authapi.Organization, error)
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
	r.POST("", rs.createOrganization)
	// Everything after here requires authentication
	//authMiddleware := authMw.Authenticate()
	//r.Use(authMiddleware)

	r.GET("/ping", rs.ping)
	r.GET("", rs.listAuthorized)

}

var (
	ErrEmailAlreadyExists   = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "The user already exists"}}
	ErrPasswordsNotMatching = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Passwords do not match"}}
	ErrPasswordNotValid     = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Password is not in the required format"}}
	UnknownError            = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem registering."}}
	ErrAuth0Unknown         = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem registering with provider."}}
)

type createOrgUserReq struct {
	OrganizationName string `json:"orgName" validate:"required,min=4"`
	FirstName        string `json:"firstName" validate:"required,min=2"`
	LastName         string `json:"lastName" validate:"required,min=2"`
	Password         string `json:"password" validate:"required"`
	PasswordConfirm  string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email            string `json:"email" validate:"required,email"`
}

type listAuthorizedResp struct {
	Organizations []authapi.Organization `json:"orgs"`
}

func (rs *OrganizationResource) listAuthorized(c echo.Context) error {
	log.Println("Inside listAuthorized(first)")

	u := authapi.User{
		ExternalID: c.Get("sub").(string),
	}

	organizations, err := rs.Store.ListAccessible(&u, false)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, UnknownError)

	}
	resp := listAuthorizedResp{
		Organizations: organizations,
	}
	return c.JSON(http.StatusOK, resp)
}

func (rs *OrganizationResource) createOrganization(c echo.Context) error {
	log.Println("Inside CreateOrganization(first)")
	r := new(createOrgUserReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	if r.Password != r.PasswordConfirm {
		return c.JSON(http.StatusBadRequest, ErrPasswordsNotMatching)
	}

	if !AuthUtil.VerifyPassword(r.Password) {
		return c.JSON(http.StatusBadRequest, ErrPasswordNotValid)
	}

	organization := authapi.Organization{Name: r.OrganizationName, Active: true}

	u := authapi.User{
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		ExternalID: "",
		Active:     true,
	}
	x := uuid.New()
	cu := authapi.OrganizationUser{Organization: &organization, User: &u, UUID: x, RoleID: 500}
	externalID, err := Auth0.CreateUser(u)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			if errCode == authapi.ECONFLICT {
				return c.JSON(http.StatusConflict, ErrEmailAlreadyExists)
			} else {
				return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
			}
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	u.ExternalID = externalID
	err = rs.Store.Create(cu)
	if err != nil {
		log.Println(err)
		err = Auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}

	err = Auth0.SendVerificationEmail(u)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}

	return c.JSON(http.StatusCreated, "")

}

func (rs *OrganizationResource) ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}
