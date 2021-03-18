package transport

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/strategicalignment"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternal  = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	ErrInvalidID = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid ID"}
)

type HTTP struct {
	svc strategicalignment.Service
}

func NewHTTP(svc strategicalignment.Service, er *echo.Group, db *sqlx.DB) {
	h := HTTP{svc}
	pr := er.Group("/strategicalignments")
	pr.POST("", h.create, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	pr.GET("", h.list, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	pr.PUT("/:id", h.update, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
	pr.DELETE("/:id", h.delete, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin"}))
}

func (h *HTTP) delete(c echo.Context) error {

	id := c.Param("id") //Project UUID
	saUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInvalidID)
	}
	orgID := c.Get("orgID").(int)
	err = h.svc.Delete(c, saUUID, orgID)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.NoContent(http.StatusOK)

}

type updateReq struct {
	Name string `json:"name"  validate:"required" validate:"required,min=10"`
}

func (h *HTTP) update(c echo.Context) error {
	log.Println("Inside updateReq()")
	r := new(updateReq)

	if err := c.Bind(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
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
	}
	orgID := c.Get("orgID").(int)

	id := c.Param("id") //SA UUID
	saUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInvalidID)
	}

	sa := model.StrategicAlignment{}
	sa.UUID = saUUID
	sa.Name = r.Name
	sa.OrganizationID = int(orgID)

	err = h.svc.Update(c, sa)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.NoContent(http.StatusOK)

}

func (h *HTTP) list(c echo.Context) error {
	orgID := c.Get("orgID").(int)
	alignments, err := h.svc.List(c, orgID)
	if err != nil {
		log.Println(err)
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.JSON(http.StatusOK, alignments)

}

type createReq struct {
	Name string `json:"name" validate:"required,min=10"`
}

func (h *HTTP) create(c echo.Context) error {
	log.Println("Inside create()")
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
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
	}
	orgID := c.Get("orgID").(int)
	_, err := h.svc.Create(c, r.Name, orgID)
	if err != nil {
		if errCode := authapi.ErrorType(err); errCode != "" {
			return c.JSON(http.StatusInternalServerError, ErrInternal)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternal)
	}
	return c.NoContent(http.StatusCreated)
}
