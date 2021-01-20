package transport

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user"
	auth0 "github.com/jpurdie/authapi/pkg/utl/Auth0"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo/v4"
)

type HTTP struct {
	svc user.Service
}

func NewHTTP(svc user.Service, er *echo.Group, db *sqlx.DB) {
	h := HTTP{svc}
	ur := er.Group("/users")
	ur.GET("/me", h.fetchMe, authMw.Authenticate())
	ur.GET("/roles", h.listRoles, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	ur.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	ur.PATCH("/:id", h.patchUser, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	ur.POST("/validationemails", h.sendValidationEmail, authMw.Authenticate())
}

var (
	ErrRoleNotFound = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid role"}
	ErrUserNotFound = authapi.Error{CodeInt: http.StatusNotFound, Message: "User not found"}
	ErrInternal     = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	ErrModifySelf   = authapi.Error{CodeInt: http.StatusUnauthorized, Message: "You cannot modify yourself"}
	ErrOneOwner     = authapi.Error{CodeInt: http.StatusUnauthorized, Message: "There must be at least one owner"}
	ErrAuth0Unknown = authapi.Error{CodeInt: http.StatusBadRequest, Message: "There was a problem with the auth provider."}
)

type fetchUserRes struct {
	User userResp `json:"user"`
}

type userResp struct {
	User authapi.User `json:"user"`
}

type patchRequest struct {
	RoleName *string `json:"role,omitempty"`
}

func (h *HTTP) sendValidationEmail(c echo.Context) error {
	externalID := c.Get("sub").(string)
	u := authapi.User{}
	u.ExternalID = externalID
	err := auth0.SendVerificationEmail(u)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) fetchMe(c echo.Context) error {
	externalID := c.Get("sub").(string)
	userFromDB, err := h.svc.FetchByExternalID(c, externalID)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	if userFromDB.UUID != uuid.Nil {
		var tempUser = userResp{}
		tempUser.User = userFromDB
		return c.JSON(http.StatusOK, tempUser)
	}
	return c.NoContent(http.StatusNotFound)
}

func (h *HTTP) listRoles(c echo.Context) error {
	roles, err := h.svc.ListRoles(c)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, roles)
}

func (h *HTTP) list(c echo.Context) error {
	orgID := c.Get("orgID").(uint)

	users, err := h.svc.List(c, orgID)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	//resp := listUsersResp{
	//	Users: users,
	//}
	return c.JSON(http.StatusOK, users)

}

func (h *HTTP) patchUser(c echo.Context) error {
	log.Println("Inside patchUser()")

	r := new(patchRequest)

	if err := c.Bind(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}
	}

	//Get the url parameter and parse it into UUID
	userUUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, ErrUserNotFound)
	}
	//get org ID in request
	requestOrgID := c.Get("orgID").(uint)
	requestUserID := int(c.Get("userID").(uint))
	//requestProfileID := int(c.Get("profileID").(uint))

	tempUser := authapi.User{}
	tempUser.UUID = userUUID
	tempOrg := authapi.Organization{}
	tempOrg.ID = int(requestOrgID)

	userToBeUpdated, err := h.svc.FetchUserByUUID(c, userUUID, requestOrgID)
	profileID := uint(0)

	for _, tempProf := range userToBeUpdated.Profile {
		if tempProf.OrganizationID == int(requestOrgID) {
			profileID = uint(tempProf.ID)
			break
		}
	}

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if requestUserID == userToBeUpdated.Profile[0].ID {
		return c.JSON(http.StatusUnprocessableEntity, ErrModifySelf)
	}
	roleLevel := 0
	if r.RoleName != nil { //checking if the role is being changed
		//TODO Cache this
		roles, err := h.svc.ListRoles(c) // list all the roles in the DB
		if err != nil {                  //check if there was a problem getting roles
			log.Println(err)
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		} //no problem getting roles
		roleFound := false // checking the role is a valid type
		for _, role := range roles {
			if strings.ToUpper(role.Name) == strings.ToUpper(*r.RoleName) {
				roleLevel = int(role.AccessLevel)
				roleFound = true
				break
			}
		}
		if !roleFound {
			return c.JSON(http.StatusUnprocessableEntity, ErrRoleNotFound)
		}
	}
	//role found if it made it here
	err = h.svc.UpdateRole(c, roleLevel, profileID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, ErrInternal)
	}
	return c.NoContent(http.StatusOK)
}
