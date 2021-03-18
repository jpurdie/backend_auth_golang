package transport

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/project"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternal     = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	ErrInvalidID    = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid ID"}
	ErrInvalidValue = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid Value"}
)

type HTTP struct {
	svc project.Service
}

func NewHTTP(svc project.Service, r *echo.Group, db *sqlx.DB) {
	h := HTTP{svc}
	pg := r.Group("/projects")
	pg.GET("/statuses", h.listStatuses, authMw.Authenticate())
	pg.GET("/types", h.listTypes, authMw.Authenticate())
	pg.GET("/complexities", h.listComplexities, authMw.Authenticate())
	pg.GET("/sizes", h.listSizes, authMw.Authenticate())
	pg.POST("", h.create, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	pg.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	pg.GET("/:pID", h.view, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	pg.PATCH("/:pID", h.update, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
}

var validate *validator.Validate

// https://restfulapi.net/http-methods/#patch
//{ “op”: “replace”, “path”: “/email”, “value”: “new.email@example.org” }
type patchReq struct {
	Op    string      `json:"op" validate:"required"`
	Path  string      `json:"path" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

func (h *HTTP) update(c echo.Context) error {

	pID := c.Param("pID") //Project UUID
	oID := c.Get("orgID").(int)

	pUUID, err := uuid.Parse(pID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInvalidID)
	}

	patchReq := new(patchReq)
	if err := c.Bind(patchReq); err != nil {
		return err
	}

	p := model.Project{}
	p.UUID = pUUID
	p.OrganizationID = int(oID)

	switch patchReq.Path {

	case "/openForTimeEntry":
		// silly but need to check if it was a boolean value that was passed
		if patchReq.Value == true || patchReq.Value == false {
			err := h.svc.Update(c, oID, pUUID, "openForTimeEntry", patchReq.Value)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.NoContent(http.StatusOK)
		}
		return c.JSON(http.StatusInternalServerError, ErrInvalidValue)
	case "/timeConstrained":
		// silly but need to check if it was a boolean value that was passed
		if patchReq.Value == true || patchReq.Value == false {
			err := h.svc.Update(c, oID, pUUID, "timeConstrained", patchReq.Value)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.NoContent(http.StatusOK)
		}
		return c.JSON(http.StatusInternalServerError, ErrInvalidValue)
	case "/compliance":
		// silly but need to check if it was a boolean value that was passed
		if patchReq.Value == true || patchReq.Value == false {
			err := h.svc.Update(c, oID, pUUID, "compliance", patchReq.Value)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.NoContent(http.StatusOK)
		}
		return c.JSON(http.StatusInternalServerError, ErrInvalidValue)
	case "/rgt":
		str := strings.ToUpper(fmt.Sprintf("%v", patchReq.Value))
		if str != "R" && str != "G" && str != "T" {
			return c.JSON(http.StatusBadRequest, ErrInvalidValue)
		}
		err := h.svc.Update(c, oID, pUUID, "rgt", str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/type":
		str := strings.ToUpper(fmt.Sprintf("%v", patchReq.Value))
		err := h.svc.Update(c, oID, pUUID, "type", str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/complexity":
		str := strings.ToUpper(fmt.Sprintf("%v", patchReq.Value))
		err := h.svc.Update(c, oID, pUUID, "complexity", str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/size":
		str := strings.ToUpper(fmt.Sprintf("%v", patchReq.Value))
		err := h.svc.Update(c, oID, pUUID, "size", str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/status":
		str := strings.ToUpper(fmt.Sprintf("%v", patchReq.Value))
		err := h.svc.Update(c, oID, pUUID, "status", str)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/strategicAlignments":
		var data = patchReq.Value.([]interface{})
		fmt.Println(data)

		if len(data) == 0 {
			return c.JSON(http.StatusBadRequest, ErrInvalidValue)
		}
		err := h.svc.Update(c, oID, pUUID, "strategicAlignments", data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/sponsorAreas":
		var data = patchReq.Value.([]interface{})
		fmt.Println(data)

		if len(data) == 0 {
			return c.JSON(http.StatusBadRequest, ErrInvalidValue)
		}
		err := h.svc.Update(c, oID, pUUID, "sponsorAreas", data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/name":
		var data = patchReq.Value.(interface{})
		fmt.Println(data)
		err := h.svc.Update(c, oID, pUUID, "name", data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	case "/description":
		var data = patchReq.Value.([]interface{})
		fmt.Println(data)

		if len(data) == 0 {
			return c.JSON(http.StatusBadRequest, ErrInvalidValue)
		}
		err := h.svc.Update(c, oID, pUUID, "description", data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.NoContent(http.StatusOK)
	default:
		fmt.Println("Too far away.")
	}
	return c.NoContent(http.StatusOK)
}
func (h *HTTP) view(c echo.Context) error {
	oID := c.Get("orgID").(int)

	pID := c.Param("pID") //Project UUID
	pUUID, err := uuid.Parse(pID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInvalidID)
	}

	p := model.Project{}
	p.UUID = pUUID
	p.OrganizationID = int(oID)

	project, err := h.svc.View(c, p.UUID)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, project)

}

type CreateResp struct {
	ID uuid.UUID `json:"id"`
}

func (h *HTTP) create(c echo.Context) error {

	projToBeCreated := new(model.Project)
	if err := c.Bind(projToBeCreated); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, ErrInternal)
	}

	oID := c.Get("orgID").(int)
	projToBeCreated.OrganizationID = int(oID)

	validate = validator.New()
	err := validate.Struct(projToBeCreated)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		// from here you can create your own error messages in whatever language you wish
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	createdProjUUID, err := h.svc.Create(c, *projToBeCreated)
	if err != nil {
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	projResp := CreateResp{}
	projResp.ID = createdProjUUID

	return c.JSON(http.StatusCreated, projResp)
}

func (h *HTTP) list(c echo.Context) error {

	oID := c.Get("orgID").(int)
	filters := make(map[string]string)
	openForTimeEntry := c.QueryParam("open_for_time_entry")

	if openForTimeEntry == "true" || openForTimeEntry == "false" {
		filters["openfortimeentry"] = openForTimeEntry
	}

	projects, err := h.svc.List(c, oID, filters)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *HTTP) listSizes(c echo.Context) error {
	sizes, err := h.svc.ListSizes(c)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, sizes)
}

func (h *HTTP) listTypes(c echo.Context) error {
	types, err := h.svc.ListTypes(c)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, types)
}

func (h *HTTP) listStatuses(c echo.Context) error {
	statuses, err := h.svc.ListStatuses(c)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, statuses)
}

func (h *HTTP) listComplexities(c echo.Context) error {
	complexities, err := h.svc.ListComplexities(c)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, complexities)
}
