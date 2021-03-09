package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	ur.POST("/validationemails", h.sendValidationEmail, authMw.Authenticate())
	ur.PATCH("/:uID", h.update, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	ur.POST("/passwords/:email", h.changePassword, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))

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

func (h *HTTP) changePassword(c echo.Context) error {
	//orgID := c.Get("orgID").(int)
	email := c.Param("email")
	u := authapi.User{}
	u.Email = email
	auth0.InitiatePasswordReset(u)
	return c.NoContent(http.StatusOK)

}

func (h *HTTP) list(c echo.Context) error {
	orgID := c.Get("orgID").(int)

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

type PatchReq struct {
	Op    string      `json:"op" validate:"required"`
	Path  string      `json:"path" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

func (h *HTTP) update(c echo.Context) error {
	log.Println("Inside update()")

	userUUID, err := uuid.Parse(c.Param("uID"))
	orgID := c.Get("orgID").(int)
	if err != nil {
		fmt.Println("ioutil.ReadAll err:", err)
		return err
	}

	var patchItems []PatchReq
	result, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err:", err)
		return err
	}

	err = json.Unmarshal(result, &patchItems)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return err
	}
	fieldsToUpdate := make(map[string]string)

	for _, patchItem := range patchItems {
		switch patchItem.Path {
		case "/firstName":
			str := fmt.Sprintf("%v", patchItem.Value)
			fieldsToUpdate["firstName"] = str
			fmt.Println("First Name being Updated to " + str)
		case "/lastName":
			str := fmt.Sprintf("%v", patchItem.Value)
			fieldsToUpdate["lastName"] = str
			fmt.Println("Lase Name being Updated to " + str)
		case "/timeZone":
			str := strings.ToUpper(fmt.Sprintf("%v", patchItem.Value))
			fmt.Println("TimeZone being Updated to " + str)
			fieldsToUpdate["timeZone"] = str

		}
	}

	err = h.svc.Update(c, userUUID, orgID, fieldsToUpdate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
