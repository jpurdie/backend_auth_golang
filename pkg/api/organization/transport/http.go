package transport

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/organization"
	auth0 "github.com/jpurdie/authapi/pkg/utl/Auth0"
	AuthUtil "github.com/jpurdie/authapi/pkg/utl/auth"
	"github.com/labstack/echo/v4"
)

// HTTP represents password http transport service
type HTTP struct {
	svc organization.Service
}

func NewHTTP(svc organization.Service, er *echo.Group) {
	h := HTTP{svc}
	pr := er.Group("/organizations")
	pr.POST("", h.createOrg)
}

var (
	ErrEmailAlreadyExists   = authapi.Error{CodeInt: http.StatusConflict, Message: "The user already exists"}
	ErrPasswordsNotMatching = authapi.Error{CodeInt: http.StatusConflict, Message: "Passwords do not match"}
	ErrPasswordNotValid     = authapi.Error{CodeInt: http.StatusConflict, Message: "Password is not in the required format"}
	UnknownError            = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem registering."}
	ErrAuth0Unknown         = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem registering with provider."}
)

type createOrgUserReq struct {
	OrganizationName string `json:"orgName" validate:"required,min=4"`
	FirstName        string `json:"firstName" validate:"required,min=2"`
	LastName         string `json:"lastName" validate:"required,min=2"`
	Password         string `json:"password" validate:"required"`
	PasswordConfirm  string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email            string `json:"email" validate:"required,email"`
}

func (h *HTTP) createOrg(c echo.Context) error {
	log.Println("Inside CreateAuthProfile(first)")
	r := new(createOrgUserReq)

	if err := c.Bind(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
	}

	if r.Password != r.PasswordConfirm {
		return c.JSON(http.StatusBadRequest, ErrPasswordsNotMatching)
	}

	if !AuthUtil.VerifyPassword(r.Password) {
		return c.JSON(http.StatusBadRequest, ErrPasswordNotValid)
	}
	organization := authapi.Organization{Name: r.OrganizationName, Active: true, UUID: uuid.New()}

	u := authapi.User{
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		ExternalID: "",
		//Active:     true,
		UUID: uuid.New(),
	}
	cu := authapi.Profile{Organization: &organization, User: &u, UUID: uuid.New(), RoleID: 500, Active: true}
	externalID, err := auth0.CreateUser(u)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			if errCode == authapi.ECONFLICT {
				return c.JSON(http.StatusConflict, ErrEmailAlreadyExists)
			} else {
				return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
			}
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	if len(externalID) == 0 { //double checking external ID
		log.Println(err)
		err = auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}

	u.ExternalID = externalID
	err = h.svc.Create(c, cu)

	if err != nil {

		log.Println(err)
		err = auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}

	err = auth0.SendVerificationEmail(u)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}

	return c.JSON(http.StatusCreated, "")

}
