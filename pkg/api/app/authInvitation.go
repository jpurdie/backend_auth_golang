package app

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	authUtil "github.com/jpurdie/authapi/pkg/utl/auth"
	"github.com/jpurdie/authapi/pkg/utl/auth0"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

// Invitation defines database operations for Invitation.
type AuthInvitationStore interface {
	View(invitation authapi.Invitation) (authapi.Invitation, error)
	CreateUser(authapi.Profile, authapi.Invitation) error
}

// Invitation Resource implements account management handler.
type AuthInvitationResource struct {
	Store AuthInvitationStore
}

func NewAuthInvitationResource(store AuthInvitationStore) *AuthInvitationResource {
	return &AuthInvitationResource{
		Store: store,
	}
}
func (rs *AuthInvitationResource) router(r *echo.Group) {
	log.Println("Inside Invitation Router")
	r.GET("/:token", rs.verifyToken)
	r.POST("/users", rs.registerUser)
}

type VerifyTokenResp struct {
	Invitation authapi.Invitation `json:"invitation"`
}
type createUserReq struct {
	FirstName       string `json:"firstName" validate:"required,min=2"`
	LastName        string `json:"lastName" validate:"required,min=2"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
	URLToken        string `json:"urlToken" validate:"required"`
}

func (rs *AuthInvitationResource) registerUser(c echo.Context) error {
	log.Println("Inside registerUser(first)")
	r := new(createUserReq)
	if err := c.Bind(r); err != nil {
		return err
	}
	if r.Password != r.PasswordConfirm {
		return c.JSON(http.StatusBadRequest, ErrPasswordsNotMatching)
	}
	log.Println("Inside registerUser()")

	if !authUtil.VerifyPassword(r.Password) {
		return c.JSON(http.StatusBadRequest, ErrPasswordNotValid)
	}
	log.Println("Inside registerUser()")

	i := authapi.Invitation{TokenStr: r.URLToken}
	log.Println("Inside registerUser()")

	tokenHash := authUtil.GenerateInviteTokenHash(i.TokenStr)
	i.TokenHash = tokenHash
	i, err := rs.Store.View(i)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, UnknownError)
	}
	curTime := time.Now()
	i, err = rs.Store.View(i)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, UnknownError)
	}
	log.Println("Inside registerUser(first)")

	if i.ID == 0 || i.Used || curTime.After(*i.ExpiresAt) { //checking the invitation is not used and not past expiration
		return c.JSON(http.StatusBadRequest, authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Invitation is expired or used"}})
	}
	log.Println("Inside registerUser()")

	/*
		the token is valid at this point
	*/

	organization := *i.Organization
	log.Println("Inside registerUser()")

	u := authapi.User{
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		ExternalID: "",
		//Active:     true,
		UUID: uuid.New(),
	}
	log.Println("Inside registerUser()")

	cu := authapi.Profile{Organization: &organization, User: &u, UUID: uuid.New(), RoleID: 100, Active: true} //create as "user"
	externalID, err := auth0.CreateUser(u)
	log.Println("Inside registerUser()")

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
	log.Println("Inside registerUser()")

	if len(externalID) == 0 { //double checking external ID
		log.Println(err)
		err = auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}
	log.Println("Inside registerUser()")

	u.ExternalID = externalID
	err = rs.Store.CreateUser(cu, i)
	if err != nil {
		log.Println("Inside registerUser()")

		log.Println(err)
		err = auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}
	log.Println("Inside registerUser()")

	err = auth0.SendVerificationEmail(u)
	log.Println("Inside registerUser()")

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, UnknownError)
	}
	log.Println("Inside registerUser()")

	return c.JSON(http.StatusCreated, "")
}
func (rs *AuthInvitationResource) verifyToken(c echo.Context) error {
	log.Println("Inside verifyToken(first)")
	tokenPlainText := c.Param("token")

	if len(tokenPlainText) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusBadRequest, Message: "Missing token"}})
	}

	i := authapi.Invitation{TokenStr: tokenPlainText}

	tokenHash := authUtil.GenerateInviteTokenHash(i.TokenStr)
	i.TokenHash = tokenHash
	i, err := rs.Store.View(i)
	curTime := time.Now()

	if i.ID == 0 || i.Used || curTime.After(*i.ExpiresAt) { //checking the invitation is not used and not past expiration
		return c.JSON(http.StatusBadRequest, authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Invitation is expired or used"}})
	}

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, UnknownError)
	}
	resp := VerifyTokenResp{}
	resp.Invitation.Organization = i.Organization
	resp.Invitation.Email = i.Email
	resp.Invitation.ExpiresAt = i.ExpiresAt

	return c.JSON(http.StatusOK, resp)
}
