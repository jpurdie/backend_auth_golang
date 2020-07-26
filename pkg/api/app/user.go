package app

import (
	"github.com/jpurdie/authapi"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

type UserStore interface {
	List(o *authapi.Organization) ([]authapi.OrganizationUser, error)
	ListRoles() ([]authapi.Role, error)
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
	log.Println("Inside Organization Router")
	r.GET("", rs.list, authMw.CheckAuthorization([]string{"owner", "admin"}))
	r.GET("/roles", rs.listRoles, authMw.CheckAuthorization([]string{"owner", "admin"}))
}

type listUsersResp struct {
	Users []userResp `json:"users"`
}

type userResp struct {
	authapi.User
	Role authapi.Role `json:"role"`
}

type roleResp struct {
	Roles []authapi.Role `json:"roles"`
}

func (rs *UserResource) listRoles(c echo.Context) error {
	log.Println("Inside list(first)")

	roles, err := rs.Store.ListRoles()

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	rResp := roleResp{}
	rResp.Roles = roles

	return c.JSON(http.StatusOK, rResp)

}

func (rs *UserResource) list(c echo.Context) error {
	log.Println("Inside list(first)")
	orgID, _ := strconv.Atoi(c.Get("orgID").(string))
	o := authapi.Organization{}
	o.ID = orgID

	orgUsers, err := rs.Store.List(&o)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorCode(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
		}
		return c.JSON(http.StatusInternalServerError, ErrAuth0Unknown)
	}

	/*
		loop through all the OrgUsers and convert them into a "flattened" version.
		I'm not seeing a circumstance when the response should be an object like this:
		OrgUser { User {	} }
	*/

	var usersSlice []userResp
	for _, tempOrgUser := range orgUsers { //
		if tempOrgUser.Active {
			var tempUser = userResp{}
			tempUser.UUID = tempOrgUser.UUID
			tempUser.FirstName = tempOrgUser.User.FirstName
			tempUser.LastName = tempOrgUser.User.LastName
			tempUser.Address = tempOrgUser.User.Address
			tempUser.Email = tempOrgUser.User.Email
			tempUser.Mobile = tempOrgUser.User.Mobile
			tempUser.Phone = tempOrgUser.User.Phone
			tempUser.Role = *tempOrgUser.Role
			tempUser.Active = tempOrgUser.Active
			usersSlice = append(usersSlice, tempUser)
		}
	}
	resp := listUsersResp{
		Users: usersSlice,
	}
	return c.JSON(http.StatusOK, resp)

}
