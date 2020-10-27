package invitation

import (
	"github.com/go-pg/pg/v9"
	"github.com/jpurdie/authapi"
	authUtil "github.com/jpurdie/authapi/pkg/utl/Auth"
	email "github.com/jpurdie/authapi/pkg/utl/mail"
	"github.com/labstack/echo"
	"os"
	"strconv"
	"time"
)

func (i Invitation) Create(c echo.Context, invite authapi.Invitation) error {
	op := "Create"

	//userID, err := strconv.Atoi(c.Get("userID").(string))
	//if err != nil {
	//	return err
	//}
	//
	//orgID, err := strconv.Atoi(c.Get("orgID").(string))
	//if err != nil {
	//	return err
	//}

	numDaysExp, err := strconv.Atoi(os.Getenv("INVITATION_EXPIRATION_DAYS"))
	if err != nil {
		numDaysExp = 7 //defaulting to 7 days
	}
	now := time.Now()
	expiresAt := now.AddDate(0, 0, numDaysExp)
	tokenStr := authUtil.GenerateInviteToken()
	tokenHash := authUtil.GenerateInviteTokenHash(tokenStr)
	invite.ExpiresAt = &expiresAt
	invite.TokenHash = tokenHash
	invite.TokenStr = tokenStr

	/*
		Checking if the invitation email already exists for that org
	 */
	user, err := i.udb.FetchByEmail(i.db, invite.Email)
	if err != nil {
		return err
	}

	for _, tempProfile := range user.Profile{
		if tempProfile.Organization.ID == invite.OrganizationID {
			return &authapi.Error{
				Op:   op,
				Code: authapi.ECONFLICT,
			}// the user exists for that organization
		}
	}

	err = i.idb.Create(i.db, invite)
	if err != nil {
		return err
	}

	err = email.SendInvitationEmail(&invite)
	if err != nil {
		return err
	}

	return nil
}

func (i Invitation) View(c echo.Context, token string) (authapi.Invitation, error) {
	return i.idb.View(i.db, token)
}

func (i Invitation) Delete(c echo.Context, invite authapi.Invitation) error {
	return i.idb.Delete(i.db, invite)
}

func (i Invitation) List(c echo.Context, o *authapi.Organization, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error) {
	return i.idb.List(i.db, o, includeExpired, includeUsed)
}

func (i Invitation) CreateUser(c echo.Context, cu authapi.Profile, invite authapi.Invitation) error {

	err := i.db.RunInTransaction(func(tx *pg.Tx) error {
		return i.idb.CreateUser(tx, cu, invite)
	})
	return err
}

//
//
//
//type InvitationStore interface {
//	Create(i authapi.Invitation) (authapi.Invitation, error)
//}
//
//// Invitation Resource implements account management handler.
//type InvitationResource struct {
//	Store InvitationStore
//}
//
//func NewInvitationResource(store InvitationStore) *InvitationResource {
//	return &InvitationResource{
//		Store: store,
//	}
//}
//func (rs *InvitationResource) router(r *echo.Group) {
//	r.GET("", rs.list, authMw.CheckAuthorization([]string{"owner", "admin"}))
//	r.POST("", rs.create, authMw.CheckAuthorization([]string{"owner", "admin"}))
//	r.DELETE("/:email", rs.delete, authMw.CheckAuthorization([]string{"owner", "admin"}))
//}
//
//var (
//	CannotFindInvitationErr = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusNotFound, Message: "Cannot find invitation"}}
//)
//
//func (rs *InvitationResource) delete(c echo.Context) error {
//	if len(c.Param("email")) == 0 {
//		return c.JSON(http.StatusNotFound, CannotFindInvitationErr)
//	}
//	orgID, _ := strconv.Atoi(c.Get("orgID").(string))
//
//	i := authapi.Invitation{
//		OrganizationID: orgID,
//		Email:          c.Param("email"),
//	}
//	//delete invite
//	err := rs.Store.Delete(i)
//	if err != nil {
//		log.Println(err)
//		if errCode := authapi.ErrorCode(err); errCode != "" {
//			return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//		}
//		return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//	}
//	return c.JSON(http.StatusOK, "")
//}
//
//type createInvitationReq struct {
//	Email string `json:"email" validate:"required,email"`
//}
//
//func (rs *InvitationResource) create(c echo.Context) error {
//	//op := "create"
//	r := new(createInvitationReq)
//
//	if err := c.Bind(r); err != nil {
//		log.Println(err)
//		return err
//	}
//
//	u := authapi.User{}
//
//	numDaysExp, err := strconv.Atoi(os.Getenv("INVITATION_EXPIRATION_DAYS"))
//	if err != nil {
//		numDaysExp = 7 //defaulting to 7 days
//	}
//	now := time.Now()
//	expiresAt := now.AddDate(0, 0, numDaysExp)
//	tokenStr := authUtil.GenerateInviteToken()
//	tokenHash := authUtil.GenerateInviteTokenHash(tokenStr)
//
//	invitorID, _ := strconv.Atoi(c.Get("userID").(string))
//	orgID, _ := strconv.Atoi(c.Get("orgID").(string))
//
//	invite := authapi.Invitation{
//		Email:          r.Email,
//		Invitor:        &u,
//		TokenStr:       tokenStr,
//		ExpiresAt:      &expiresAt,
//		TokenHash:      tokenHash,
//		InvitorID:      invitorID,
//		OrganizationID: orgID,
//	}
//	log.Println("This should be emailed " + invite.TokenStr)
//	log.Println("This should be saved to db " + invite.TokenHash)
//
//	err = rs.Store.Delete(invite)
//	if err != nil {
//		log.Println(err)
//		if errCode := authapi.ErrorCode(err); errCode != "" {
//			return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//		}
//		return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//	}
//
//	//save invite
//	_, err = rs.Store.Create(invite)
//	if err != nil {
//		log.Println(err)
//		if errCode := authapi.ErrorCode(err); errCode != "" {
//			return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//		}
//		return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//	}
//
//	//send invite after save
//	err = email.SendInvitationEmail(&invite)
//	if err != nil {
//		log.Println(err)
//		if errCode := authapi.ErrorCode(err); errCode != "" {
//			return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//		}
//		return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//	}
//
//	return c.JSON(http.StatusCreated, "")
//}
//
////type listInvitationsResp struct {
////	Invitations []authapi.Invitation `json:"invitations"`
////}
//
//func (rs *InvitationResource) list(c echo.Context) error {
//	log.Println("Inside listTokens(first)")
//	o := authapi.Organization{}
//	o.ID, _ = strconv.Atoi(c.Get("orgID").(string))
//
//	invitations, err := rs.Store.List(&o, false, false)
//
//	if err != nil {
//		log.Println(err)
//		if errCode := authapi.ErrorCode(err); errCode != "" {
//			return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//		}
//		return c.JSON(http.StatusInternalServerError, app.ErrAuth0Unknown)
//	}
//	//resp := listInvitationsResp{
//	//	Invitations: invitations,
//	//}
//	return c.JSON(http.StatusOK, invitations)
//}
//
//func (rs *InvitationResource) ping(c echo.Context) error {
//	return c.JSON(http.StatusOK, "pong")
//}
