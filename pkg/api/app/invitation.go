package app

import (
	"github.com/jpurdie/authapi"
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

type InvitationStore interface {
	List(u *authapi.User, includeExpired bool) ([]authapi.Invitation, error)
	Create(i authapi.Invitation) (authapi.Invitation, error)
}

// Invitation Resource implements account management handler.
type InvitationResource struct {
	Store InvitationStore
}

func NewInvitationResource(store InvitationStore) *InvitationResource {
	return &InvitationResource{
		Store: store,
	}
}
func (rs *InvitationResource) router(r *echo.Group) {
	r.GET("", rs.list, authMw.CheckAuthorization([]string{"owner", "admin"}))
	r.POST("", rs.create, authMw.CheckAuthorization([]string{"owner", "admin"}))
}

var (
//	ErrEmailAlreadyExists   = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "The user already exists"}}

)

type createInvitationResp struct {
	Invitation authapi.Invitation `json:"invitation"`
}
type createInvitationReq struct {
	Email string `json:"email" validate:"required,email"`
}

func (rs *InvitationResource) create(c echo.Context) error {
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
		Token:          tokenHash,
		InvitorID:      invitorID,
		OrganizationID: orgID,
	}
	log.Println("This should be emailed " + invite.TokenStr)
	log.Println("This should be saved to db " + invite.Token)

	//save invite
	_, err = rs.Store.Create(invite)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	//send invite after save
	err = email.SendInvitationEmail(&invite)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	return c.JSON(http.StatusCreated, "")
}

type listInvitationsResp struct {
	Invitations []authapi.Invitation `json:"invitations"`
}

func (rs *InvitationResource) list(c echo.Context) error {
	log.Println("Inside listTokens(first)")

	u := authapi.User{
		ExternalID: c.Get("sub").(string),
	}

	Invitations, err := rs.Store.List(&u, false)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}
	resp := listInvitationsResp{
		Invitations: Invitations,
	}
	return c.JSON(http.StatusOK, resp)
}

func (rs *InvitationResource) ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}
