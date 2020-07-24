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
}

type listUsersResp struct {
	Users []authapi.User `json:"users"`
}

func (rs *UserResource) list(c echo.Context) error {
	log.Println("Inside list(first)")

	// TODO: dont ignore the error
	// userID, _ := strconv.Atoi(c.Get("userID").(string))
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
	//TODO get role into response
	var userResp []authapi.User
	for _, tempOrgUser := range orgUsers {
		var tempUser = tempOrgUser.User
		tempUser.UUID = tempOrgUser.UUID
		userResp = append(userResp, *tempUser)

	}
	resp := listUsersResp{
		Users: userResp,
	}
	return c.JSON(http.StatusOK, resp)

}
