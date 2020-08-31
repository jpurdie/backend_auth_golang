package transport

import (
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/invitation"
	auth0 "github.com/jpurdie/authapi/pkg/utl/Auth0"
	authUtil "github.com/jpurdie/authapi/pkg/utl/auth"
	email "github.com/jpurdie/authapi/pkg/utl/mail"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	ErrInternal = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}}
	CannotFindInvitationErr = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusNotFound, Message: "Cannot find invitation"}}
	ErrPasswordsNotMatching = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Passwords do not match"}}
	ErrPasswordNotValid     = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Password is not in the required format"}}
	ErrEmailAlreadyExists   = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "The user already exists"}}
	ErrAuth0Unknown         = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem registering with provider."}}

)

type HTTP struct {
	svc invitation.Service
}

func NewHTTP(svc invitation.Service, er *echo.Group, db *pg.DB) {
	h := HTTP{svc}
	ig := er.Group("/invitations")
	ig.GET("/:token", h.verifyToken)
	ig.POST("/users", h.createUser)

	//everything after here requires auth
	//authMiddleware :=
	//ig.Use(authMiddleware)

	ig.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin"}))
	ig.POST("", h.create, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin"}))
	ig.DELETE("/:email", h.delete, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin"}))
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


func (h *HTTP) createUser(c echo.Context) error {
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
	i, err := h.svc.View(c, tokenHash)
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
	o := authapi.Organization{}
	o.ID, _ = strconv.Atoi(c.Get("orgID").(string))

	invitations, err := h.svc.List(c, &o, false, false)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	//resp := listInvitationsResp{
	//	Invitations: invitations,
	//}
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

	numDaysExp, err := strconv.Atoi(os.Getenv("INVITATION_EXPIRATION_DAYS"))
	if err != nil {
		numDaysExp = 7 //defaulting to 7 days
	}
	now := time.Now()
	expiresAt := now.AddDate(0, 0, numDaysExp)
	tokenStr := authUtil.GenerateInviteToken()
	tokenHash := authUtil.GenerateInviteTokenHash(tokenStr)

	invitorID, _ := strconv.Atoi(c.Get("userID").(string))
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))

	invite := authapi.Invitation{
		Email:          r.Email,
		Invitor:        &u,
		TokenStr:       tokenStr,
		ExpiresAt:      &expiresAt,
		TokenHash:      tokenHash,
		InvitorID:      invitorID,
		OrganizationID: orgID,
	}
	log.Println("This should be emailed " + invite.TokenStr)
	log.Println("This should be saved to db " + invite.TokenHash)

	err = h.svc.Delete(c, invite)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	//save invite
	err = h.svc.Create(c, invite)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	//send invite after save
	err = email.SendInvitationEmail(&invite)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	return c.JSON(http.StatusCreated, "")
}

type listInvitationsResp struct {
	Invitations []authapi.Invitation `json:"invitations"`
}

func (h *HTTP) delete(c echo.Context) error {
	if len(c.Param("email")) == 0 {
		return c.JSON(http.StatusNotFound, CannotFindInvitationErr)
	}
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))

	i := authapi.Invitation{
		OrganizationID: orgID,
		Email:          c.Param("email"),
	}
	//delete invite
	err := h.svc.Delete(c, i)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, "")
}

func (h *HTTP) verifyToken(c echo.Context) error {
	log.Println("Inside verifyToken(first)")
	tokenPlainText := c.Param("token")

	if len(tokenPlainText) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusBadRequest, Message: "Missing token"}})
	}

	i := authapi.Invitation{TokenStr: tokenPlainText}

	tokenHash := authUtil.GenerateInviteTokenHash(i.TokenStr)
	i.TokenHash = tokenHash
	i, err := h.svc.View(c, tokenHash)
	curTime := time.Now()

	if i.ID == 0 || i.Used || curTime.After(*i.ExpiresAt) { //checking the invitation is not used and not past expiration
		return c.JSON(http.StatusBadRequest, authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "Invitation is expired or used"}})
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