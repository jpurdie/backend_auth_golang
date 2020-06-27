package app

import (
	"context"
	"fmt"
	"github.com/jpurdie/authapi"
	authUtil "github.com/jpurdie/authapi/pkg/utl/Auth"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type InvitationStore interface {
	List(u *authapi.User, includeExpired bool) ([]authapi.Invitation, error)
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

	r := new(createInvitationReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	u := authapi.User{
		ExternalID: c.Get("sub").(string),
	}

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
	log.Println(invite.TokenStr)
	log.Println(invite.Token)

	// Create an instance of the Mailgun Client

	// TODO: Move this to a util!
	//mg := mailgun.NewMailgun("thepolyglotdeveloper.com", os.Getenv("MAILGUN_API_KEY"), os.Getenv("MAILGUN_PUB_VAL_KEY"))
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))

	sender := "invitations@vitae.com"
	subject := "You've been invited to Vitae"
	body := "You've been invited to Vitae! \n\nFollow this link: \thttp://localhost.vitae.com/auth/t=" + invite.TokenStr
	recipient := invite.Email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return c.JSON(http.StatusCreated, invite)
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
