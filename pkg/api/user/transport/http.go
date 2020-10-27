package transport

import (
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user"
	"github.com/labstack/echo"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"strings"
	"github.com/go-playground/validator"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type HTTP struct {
	svc user.Service
}

func NewHTTP(svc user.Service, er *echo.Group, db *pg.DB) {
	h := HTTP{svc}
	ur := er.Group("/users")
	ur.GET("/me", h.fetchMe, authMw.Authenticate())
	ur.GET("/roles", h.listRoles, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin"}))
	ur.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin"}))
	ur.PATCH("/:id", h.patchUser, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin"}))
}

var (
	ErrRoleNotFound =  authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid role"}
	ErrUserNotFound = authapi.Error{CodeInt: http.StatusNotFound, Message: "User not found"}
	ErrInternal     = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	ErrModifySelf   = authapi.Error{CodeInt: http.StatusUnauthorized, Message: "You cannot modify yourself"}
	ErrOneOwner     =  authapi.Error{CodeInt: http.StatusUnauthorized, Message: "There must be at least one owner"}
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
	roles, err := h.svc.ListRoles(c);
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	//rResp := roleResp{}
	//rResp.Roles = roles

	return c.JSON(http.StatusOK, roles)
}

func (h *HTTP) list(c echo.Context) error {
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))

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

func (h *HTTP)  patchUser(c echo.Context) error {
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
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))
	userID, _ := strconv.Atoi(c.Get("userID").(string))
	profileID, _ := strconv.Atoi(c.Get("profileID").(string))

	tempUser := authapi.User{}
	tempUser.UUID = userUUID
	tempOrg := authapi.Organization{}
	tempOrg.ID = orgID

	profileToBeUpdated, err := h.svc.FetchProfile(c, userID, orgID)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if profileID == profileToBeUpdated.ID {
		return c.JSON(http.StatusUnprocessableEntity, ErrModifySelf)
	}

	if r.RoleName != nil { //checking if the role is being changed
		//TODO Cache this
		roles, err := h.svc.ListRoles(c) // list all the roles in the DB
		if err != nil {                    //check if there was a problem getting roles
			log.Println(err)
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		} //no problem getting roles
		roleFound := false // checking the role is a valid type
		for _, role := range roles {
			if strings.ToUpper(role.Name) == strings.ToUpper(*r.RoleName) {
				profileToBeUpdated.RoleID = int(role.AccessLevel)
				roleFound = true
				break
			}
		}
		if !roleFound {
			return c.JSON(http.StatusUnprocessableEntity, ErrRoleNotFound)
		}
	}
	//role found if it made it here
	err = h.svc.Update(c, profileToBeUpdated)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, ErrInternal)
	}
	return c.NoContent(http.StatusOK)
}