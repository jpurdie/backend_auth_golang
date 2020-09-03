package invitation

import (
	"github.com/go-pg/pg/v9"
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo"
)

func (i Invitation) Create(c echo.Context, invite authapi.Invitation) error {
	//op := "Create"


	/*
		Questions:
		1. Would this file/functions be where the business logic goes?
		2. The application is a multi-tenant application. One user can belong to many organizations, and one organization can have many users.
			ER diagram: https://imgur.com/a/eD3woNv

			In order to join an existing organization, you need to be invited (Similar to slack).
			But I want to prevent someone who's already part of the organization from being invited again.

			If I want to check whether a user with the email address in the invitation already exists in the database. What would be the best way to do that?

				Am I able to use the FetchProfile function in pkg/api/user/platform/pgsql/user.go?
				Would I create a new function in /pkg/api/invitation/platform/pgsql/invitation.go to check the users table? (if we do it this way, it seems like DRY rules would be violated)

				In other languages (in this file) I could just do something like this:
				User myUser = new User();
				myUser.email = "foo@bar.com";
				myUser.findUser();

				if(myUser != null) {
					return "User already exists. dont create invitation";
				}
				//continue on

				But with Golang, things seem to be more isolated and it doesn't appear I can do that.

	 */


	//tempInvite, err := i.idb.FindUserByEmail(i.db, invite.Email, invite.Organization.ID)
	//if(err != nil){
	//	return err
	//}
	//if(tempInvite.ID > 0){
	//	return &authapi.Error{
	//		Op:   op,
	//		Code: authapi.ECONFLICT,
	//	}
	//}
	return i.idb.Create(i.db, invite)
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

func (i Invitation) CreateUser(c echo.Context,  cu authapi.Profile, invite authapi.Invitation) error {

	err := i.db.RunInTransaction(func (tx *pg.Tx) error{
		return i.idb.CreateUser(tx, cu, invite)
	})
	return err;
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
