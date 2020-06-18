package transport

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/company"
	AuthUtil "github.com/jpurdie/authapi/pkg/utl/Auth"
	"github.com/jpurdie/authapi/pkg/utl/Auth0"
	"github.com/labstack/echo"
	"net/http"
)

// HTTP represents company http service
type HTTP struct {
	svc company.Service
}

// NewHTTP creates new company http service
func NewHTTP(svc company.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/companies")
	ur.POST("", h.create)
}

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	ErrPasswordNotValid    = echo.NewHTTPError(http.StatusBadRequest, "passwords are not in the valid format")
)

type createOrgUserReq struct {
	CompanyName     string `json:"orgName" validate:"required,min=4"`
	FirstName       string `json:"firstName" validate:"required,min=2"`
	LastName        string `json:"lastName" validate:"required,min=2"`
	Password        string `json:"password" validate:"required`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
}

func (h HTTP) create(c echo.Context) error {
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
	cu := authapi.CompanyUser{Company: &company, User: &u, UUID: x}

	externalID, err := Auth0.CreateUser(u)
	if err != nil {
		return err
	}

	u.ExternalID = externalID

	companyUser, err := h.svc.Create(c, cu)
	if err != nil {
		//delete auth0 user
		return err
	}

	err = Auth0.SendVerificationEmail(u)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, companyUser)
}
