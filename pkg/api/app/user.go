package app

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserStore interface {
	List(o *authapi.Organization) ([]authapi.User, error)
	ListRoles() ([]authapi.Role, error)
	Update(p authapi.Profile) error
	Fetch(u authapi.User) (authapi.User, error)
	ListAuthorized(u *authapi.User, includeInactive bool) ([]authapi.Profile, error)
	FetchProfile(u authapi.User, o authapi.Organization) (authapi.Profile, error)
	Delete(p authapi.Profile) error
}

type UserResource struct {
	Store UserStore
}

func NewUserResource(store UserStore) *UserResource {
	return &UserResource{
		Store: store,
	}
}
func (rs *UserResource) router(r *echo.Group) {
	r.GET("/me", rs.fetchMe)
	r.GET("", rs.list, authMw.CheckAuthorization([]string{"owner", "admin"})) //lists all the users for an organization
	r.GET("/roles", rs.listRoles, authMw.CheckAuthorization([]string{"owner", "admin", "user"}))
	r.PATCH("/:id", rs.patchUser, authMw.CheckAuthorization([]string{"owner", "admin"}))
	r.DELETE("/:id", rs.delete, authMw.CheckAuthorization([]string{"owner", "admin"}))
}

var (
	ErrRoleNotFound = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid role"}}
	ErrUserNotFound = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusNotFound, Message: "User not found"}}
	ErrInternal     = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}}
	ErrModifySelf   = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusUnauthorized, Message: "You cannot modify yourself"}}
	ErrOneOwner     = authapi.ErrorResp{Error: authapi.Error{CodeInt: http.StatusUnauthorized, Message: "There must be at least one owner"}}
)

//type listUsersResp struct {
//	Users []authapi.User `json:"users"`
//}

type fetchUserRes struct {
	User userResp `json:"user"`
}

type userResp struct {
	User authapi.User `json:"user"`
}

//type roleResp struct {
//	Roles []authapi.Role `json:"roles"`
//}

type patchRequest struct {
	RoleName *string `json:"role,omitempty"`
}

func (rs *UserResource) delete(c echo.Context) error {
	log.Println("Inside delete()")

	userUUID, err := uuid.Parse(c.Param("id")) //Get the url parameter and parse it into UUID
	if err != nil {
		log.Println(err)
		return c.JSON(ErrUserNotFound.Error.CodeInt, ErrUserNotFound)
	}
	//get org ID in request
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))           //getting requesting user's org ID
	rProfileID, _ := strconv.Atoi(c.Get("rProfileID").(string)) //getting requesting user's profile ID

	tempUser := authapi.User{}
	tempUser.UUID = userUUID
	tempOrg := authapi.Organization{}
	tempOrg.ID = orgID

	profileToBeDeleted, err := rs.Store.FetchProfile(tempUser, tempOrg)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if rProfileID == profileToBeDeleted.ID {
		return c.JSON(ErrModifySelf.Error.CodeInt, ErrModifySelf)
	}
	//role found if it made it here
	err = rs.Store.Delete(profileToBeDeleted)
	if err != nil {
		log.Println(err)
		return c.JSON(ErrInternal.Error.CodeInt, ErrInternal)
	}
	return c.NoContent(http.StatusOK)

}

func (rs *UserResource) patchUser(c echo.Context) error {
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
		return c.JSON(ErrUserNotFound.Error.CodeInt, ErrUserNotFound)
	}
	//get org ID in request
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))
	rProfileID, _ := strconv.Atoi(c.Get("rProfileID").(string))

	tempUser := authapi.User{}
	tempUser.UUID = userUUID
	tempOrg := authapi.Organization{}
	tempOrg.ID = orgID

	profileToBeUpdated, err := rs.Store.FetchProfile(tempUser, tempOrg)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if rProfileID == profileToBeUpdated.ID {
		return c.JSON(ErrModifySelf.Error.CodeInt, ErrModifySelf)
	}

	if r.RoleName != nil { //checking if the role is being changed
		//TODO Cache this
		roles, err := rs.Store.ListRoles() // list all the roles in the DB
		if err != nil {                    //check if there was a problem getting roles
			log.Println(err)
			if errCode := authapi.ErrorCode(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
			}
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
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
	err = rs.Store.Update(profileToBeUpdated)
	if err != nil {
		log.Println(err)
		return c.JSON(ErrInternal.Error.CodeInt, ErrInternal)
	}
	return c.NoContent(http.StatusOK)
}

func (rs *UserResource) listRoles(c echo.Context) error {
	roles, err := rs.Store.ListRoles()
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}
	//rResp := roleResp{}
	//rResp.Roles = roles

	return c.JSON(http.StatusOK, roles)
}

func (rs *UserResource) fetchMe(c echo.Context) error {

	externalID := c.Get("sub").(string)
	u := authapi.User{}
	u.ExternalID = externalID

	userFromDB, err := rs.Store.Fetch(u)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	if userFromDB.UUID != uuid.Nil {
		var tempUser = userResp{}
		tempUser.User = userFromDB
		return c.JSON(http.StatusOK, tempUser)
	}

	return c.NoContent(http.StatusNotFound)

}

func (rs *UserResource) list(c echo.Context) error {
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))
	o := authapi.Organization{}
	o.ID = orgID

	users, err := rs.Store.List(&o)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}
	//resp := listUsersResp{
	//	Users: users,
	//}
	return c.JSON(http.StatusOK, users)

}

type listAuthorizedRespInner struct {
	OrgName string       `json:"name"`
	UUID    uuid.UUID    `json:"uuid"`
	Role    authapi.Role `json:"role"`
}

type listAuthorizedResp struct {
	Orgs []listAuthorizedRespInner `json:"orgs"`
}

func (rs *UserResource) listAuthorized(c echo.Context) error {

	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		log.Println(err)
		return c.JSON(ErrUserNotFound.Error.CodeInt, ErrUserNotFound)
	}

	u := authapi.User{}
	u.UUID = userUUID

	organizationUser, err := rs.Store.ListAuthorized(&u, false)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)

	}
	x := listAuthorizedResp{}
	for _, tempOrgUser := range organizationUser {
		temp := listAuthorizedRespInner{
			OrgName: tempOrgUser.Organization.Name,
			UUID:    tempOrgUser.Organization.UUID,
			Role:    *tempOrgUser.Role,
		}
		x.Orgs = append(x.Orgs, temp)
	}
	return c.JSON(http.StatusOK, x)
}
