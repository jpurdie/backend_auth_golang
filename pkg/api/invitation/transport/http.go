package transport

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/invitation"
	auth0 "github.com/jpurdie/authapi/pkg/utl/Auth0"
	authUtil "github.com/jpurdie/authapi/pkg/utl/auth"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternal             = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	CannotFindInvitationErr = authapi.Error{CodeInt: http.StatusNotFound, Message: "Cannot find invitation"}
	ErrPasswordsNotMatching = authapi.Error{CodeInt: http.StatusConflict, Message: "Passwords do not match"}
	ErrPasswordNotValid     = authapi.Error{CodeInt: http.StatusConflict, Message: "Password is not in the required format"}
	ErrEmailAlreadyExists   = authapi.Error{CodeInt: http.StatusConflict, Message: "The user already exists"}
	ErrAuth0Unknown         = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem registering with provider."}
)

type HTTP struct {
	svc invitation.Service
}

func NewHTTP(svc invitation.Service, r *echo.Group, db *sqlx.DB) {
	h := HTTP{svc}
	ig := r.Group("/invitations")
	ig.GET("/validations/:token", h.verifyToken)
	ig.DELETE("/:email", h.delete, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	ig.POST("/users", h.createUser)
	ig.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	ig.POST("", h.create, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
}

type createUserReq struct {
	FirstName       string `json:"firstName" validate:"required,min=2"`
	LastName        string `json:"lastName" validate:"required,min=2"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
	URLToken        string `json:"urlToken" validate:"required"`
}

func (h *HTTP) createUser(c echo.Context) error {
	log.Println("Inside registerUser(first)")
	r := new(createUserReq)
	if err := c.Bind(r); err != nil {
		return err
	}
	if r.Password != r.PasswordConfirm {
		return c.JSON(http.StatusBadRequest, ErrPasswordsNotMatching)
	}

	if !authUtil.VerifyPassword(r.Password) {
		return c.JSON(http.StatusBadRequest, ErrPasswordNotValid)
	}

	i, err := h.svc.View(c, r.URLToken)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ErrInternal)
	}
	curTime := time.Now()
	//i, err = rs.Store.View(i)
	//
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, ErrInternal)
	//}
	log.Println("Inside registerUser(first)")

	if i.ID == 0 || i.Used || curTime.After(*i.ExpiresAt) { //checking the invitation is not used and not past expiration
		return c.JSON(http.StatusBadRequest, authapi.Error{CodeInt: http.StatusConflict, Message: "Invitation is expired or used"})
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
		if errCode := authapi.ErrorType(err); errCode != "" {
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
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	log.Println("Inside registerUser()")

	u.ExternalID = externalID
	err = h.svc.CreateUser(c, cu, i)
	if err != nil {
		log.Println("Inside registerUser()")

		log.Println(err)
		err = auth0.DeleteUser(u) //need to delete user from auth0 since the database failed
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	log.Println("Inside registerUser()")

	err = auth0.SendVerificationEmail(u)
	log.Println("Inside registerUser()")

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	log.Println("Inside registerUser()")

	return c.JSON(http.StatusCreated, "")
}

func (h *HTTP) list(c echo.Context) error {
	oID := c.Get("orgID").(int)
	invitations, err := h.svc.List(c, oID, false, false)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, invitations)
}

type createInvitationReq struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *HTTP) create(c echo.Context) error {
	//op := "create"
	r := new(createInvitationReq)

	if err := c.Bind(r); err != nil {
		log.Println(err)
		return err
	}

	u := authapi.User{}
	u.Email = r.Email

	invitorID := c.Get("userID").(int)
	orgID := c.Get("orgID").(int)

	invite := authapi.Invitation{
		Email:          r.Email,
		Invitor:        &u,
		InvitorID:      invitorID,
		OrganizationID: orgID,
	}
	//save invite
	err := h.svc.Create(c, invite)
	if err != nil {
		if errType := authapi.ErrorType(err); errType != "" {
			if errType == authapi.ECONFLICT {
				return c.JSON(http.StatusConflict, ErrEmailAlreadyExists)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	return c.JSON(http.StatusCreated, "")
}

func (h *HTTP) delete(c echo.Context) error {
	if len(c.Param("email")) == 0 {
		return c.JSON(http.StatusNotFound, CannotFindInvitationErr)
	}

	oID := c.Get("orgID").(int)

	//delete invite
	err := h.svc.Delete(c, c.Param("email"), oID)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, "")
}

type VerifyTokenResp struct {
	Invitation authapi.Invitation `json:"invitation"`
}

func (h *HTTP) verifyToken(c echo.Context) error {
	log.Println("Inside verifyToken(first)")
	tokenPlainText := c.Param("token")

	if len(tokenPlainText) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, authapi.Error{CodeInt: http.StatusBadRequest, Message: "Missing token"})
	}

	i, err := h.svc.View(c, tokenPlainText)
	curTime := time.Now()

	if i.ID == 0 || i.Used || curTime.After(*i.ExpiresAt) { //checking the invitation is not used and not past expiration
		return c.JSON(http.StatusBadRequest, authapi.Error{CodeInt: http.StatusConflict, Message: "Invitation is expired or used"})
	}

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ErrInternal)
	}
	resp := VerifyTokenResp{}
	resp.Invitation.Organization = i.Organization
	resp.Invitation.Email = i.Email
	resp.Invitation.ExpiresAt = i.ExpiresAt

	return c.JSON(http.StatusOK, resp)
}
