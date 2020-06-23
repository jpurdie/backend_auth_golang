package app

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	AuthUtil "github.com/jpurdie/authapi/pkg/utl/Auth"
	"github.com/jpurdie/authapi/pkg/utl/Auth0"
	"github.com/labstack/echo"
	"net/http"
)

// CompanyStore defines database operations for Company.
type CompanyStore interface {
	Create(authapi.CompanyUser) (authapi.CompanyUser, error)
	//List() error
}

// Company Resource implements account management handler.
type CompanyResource struct {
	Store CompanyStore
}

func NewCompanyResource(store CompanyStore) *CompanyResource {
	return &CompanyResource{
		Store: store,
	}
}
func (rs *CompanyResource) router(r *echo.Group) {
	r.POST("", rs.create)
	r.GET("/ping", rs.ping)

}

var (
	ErrCompAlreadyExists   = echo.NewHTTPError(http.StatusConflict, "Company name already exists.")
	ErrEmailAlreadyExists  = echo.NewHTTPError(http.StatusConflict, "Email already exists.")
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	ErrPasswordNotValid    = echo.NewHTTPError(http.StatusBadRequest, "passwords are not in the valid format")
	UnknownError           = echo.NewHTTPError(http.StatusBadRequest, "There was an unknown error")
)

type createOrgUserReq struct {
	CompanyName     string `json:"orgName" validate:"required,min=4"`
	FirstName       string `json:"firstName" validate:"required,min=2"`
	LastName        string `json:"lastName" validate:"required,min=2"`
	Password        string `json:"password" validate:"required`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
}

func (rs *CompanyResource) create(c echo.Context) error {
	r := new(createOrgUserReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	if r.Password != r.PasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if !AuthUtil.VerifyPassword(r.Password) {
		return ErrPasswordNotValid
	}

	company := authapi.Company{Name: r.CompanyName, Active: true}

	u := authapi.User{
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		ExternalID: "",
		Active:     true,
	}
	x := uuid.New()
	cu := authapi.CompanyUser{Company: &company, User: &u, UUID: x, RoleID: 500}
	externalID, err := Auth0.CreateUser(u)
	if err != nil {
		fmt.Println(err)

		return UnknownError
	}

	u.ExternalID = externalID

	companyUser, err := rs.Store.Create(cu)
	if err != nil {
		fmt.Println(err)
		err = Auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return UnknownError
	}

	err = Auth0.SendVerificationEmail(u)

	if err != nil {
		fmt.Println(err)
		return UnknownError
	}

	return c.JSON(http.StatusCreated, companyUser)

}

func (rs *CompanyResource) ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
