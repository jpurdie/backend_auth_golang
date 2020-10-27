package transport

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/project"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrInternal  = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	ErrInvalidID = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid ID"}
)

type HTTP struct {
	svc project.Service
}

func NewHTTP(svc project.Service, r *echo.Group, db *pg.DB) {
	h := HTTP{svc}
	pg := r.Group("/projects")
	pg.GET("/statuses", h.listStatuses, authMw.Authenticate())
	pg.GET("/types", h.listTypes, authMw.Authenticate())
	pg.GET("/complexities", h.listComplexities, authMw.Authenticate())
	pg.GET("/sizes", h.listSizes, authMw.Authenticate())
	pg.POST("", h.create, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin", "user"}))
	pg.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin", "user"}))
	pg.GET("/:pID", h.view, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin", "user"}))
	pg.PATCH("/:pID", h.patch, authMw.Authenticate(), authMw.CheckAuthorization(db, []string{"owner", "admin", "user"}))
}

var validate *validator.Validate

type patchReq struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (h *HTTP) patch(c echo.Context) error {

	pID := c.Param("pID") //Project UUID
	pUUID, err := uuid.Parse(pID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInvalidID)
	}

	pr := new(patchReq)
	if err := c.Bind(pr); err != nil {
		return err
	}

	p := model.Project{}
	p.UUID = pUUID

	/*
		NEED NEED NEED to come up with a better way to handles these patch requests. reflection?
		https://stackoverflow.com/questions/38206479/golang-rest-patch-and-building-an-update-query
	*/

	switch pr.Key {
	case "complexity":
		pc := model.ProjectComplexity{}
		pc.UUID, err = uuid.Parse(pr.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInvalidID)
		}
		err := h.svc.UpdateComplexity(c, p.UUID, pc)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	case "type":
		pt := model.ProjectType{}
		pt.UUID, err = uuid.Parse(pr.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInvalidID)
		}
		err := h.svc.UpdateType(c, p.UUID, pt)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	case "rgt":
		rgt := pr.Value
		err := h.svc.UpdateRGT(c, p.UUID, rgt)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	case "status":
		ps := model.ProjectStatus{}
		ps.UUID, err = uuid.Parse(pr.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInvalidID)
		}
		err := h.svc.UpdateStatus(c, p.UUID, ps)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	case "size":
		ps := model.ProjectSize{}
		ps.UUID, err = uuid.Parse(pr.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrInvalidID)
		}
		err := h.svc.UpdateSize(c, p.UUID, ps)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	case "description":
		p.Description = pr.Value
		err := h.svc.UpdateDescr(c, p)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	case "name":
		p.Name = pr.Value
		err := h.svc.UpdateName(c, p)
		if err != nil {
			if errCode := authapi.ErrorType(err); errCode != "" {
				return c.JSON(http.StatusInternalServerError, ErrInternal)
			}
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
	default:
		fmt.Println("Too far away.")
	}
	return c.NoContent(http.StatusOK)
}
func (h *HTTP) view(c echo.Context) error {

	pID := c.Param("pID")

	project, err := h.svc.View(c, pID)
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

	p := new(model.Project)
	if err := c.Bind(p); err != nil {
		return err
	}
	oID, _ := strconv.Atoi(c.Get("orgID").(string))
	p.OrganizationID = oID

	validate = validator.New()
	err := validate.Struct(p)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		// from here you can create your own error messages in whatever language you wish
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	createdP, err := h.svc.Create(c, *p)

	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}

	projResp := CreateResp{}
	projResp.ID = createdP.UUID

	return c.JSON(http.StatusCreated, projResp)
}

func (h *HTTP) list(c echo.Context) error {
	oID, _ := strconv.Atoi(c.Get("orgID").(string))

	projects, err := h.svc.List(c, strconv.Itoa(oID))
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
