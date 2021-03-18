package transport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jpurdie/authapi/pkg/utl/helpers"

	"github.com/jpurdie/authapi"

	"gopkg.in/guregu/null.v4"

	"github.com/jpurdie/authapi/pkg/utl/model"

	"github.com/google/uuid"

	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"

	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/api/capacityplan"
	"github.com/labstack/echo/v4"
)

type HTTP struct {
	svc capacityplan.Service
}

var (
	ErrInternal     = authapi.Error{CodeInt: http.StatusConflict, Message: "There was a problem"}
	ErrInvalidID    = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid ID"}
	ErrInvalidValue = authapi.Error{CodeInt: http.StatusUnprocessableEntity, Message: "Invalid Value"}
)

func NewHTTP(svc capacityplan.Service, r *echo.Group, db *sqlx.DB) {
	h := HTTP{svc}
	ig := r.Group("/capacityplans")
	ig.POST("", h.create, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	ig.DELETE("/:capEntryUUID", h.deleteByID, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	ig.GET("/:resourceID/:startDate/:endDate", h.list, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
	ig.GET("/summaries/:resourceID/:startDate/:endDate", h.listSummary, authMw.Authenticate(), authMw.CheckAuthorization(*db, []string{"owner", "admin", "user"}))
}

type capPlanReq struct {
	UUID        uuid.UUID `json:"id,omitempty"`
	ProjectID   uuid.UUID `json:"projectID" validate:"required,uuid,unique"`
	WorkPercent int       `json:"workPercent" validate:"required,int"`
}
type capPlanArrReq struct {
	CapacityPlans []capPlanReq `json:"capacityPlanEntries,array" validate:"required"`
	ResourceID    string       `json:"resourceID" validate:"required"`
	WorkDate      string       `json:"workDate" validate:"required"`
}

func (h *HTTP) listSummary(c echo.Context) error {
	orgID := c.Get("orgID").(int)
	resourceToGet := c.Param("resourceID")
	resourceUUIDToGet, err := uuid.Parse(resourceToGet)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ErrInvalidValue)
	}

	startDate, _ := time.Parse(helpers.LayoutISO, c.Param("startDate"))
	endDate, _ := time.Parse(helpers.LayoutISO, c.Param("endDate"))

	daysDuration := endDate.Sub(startDate).Hours() / 24
	if daysDuration > 31 {
		return c.JSON(http.StatusUnprocessableEntity, ErrInvalidValue)
	}

	results, _ := h.svc.ListSummary(c, int(orgID), resourceUUIDToGet, startDate, endDate)

	return c.JSON(http.StatusOK, results)

}

func (h *HTTP) list(c echo.Context) error {
	orgID := c.Get("orgID").(int)
	//userID := c.Get("userID").(int)
	resourceToGet := c.Param("resourceID")

	resourceUUIDToGet, err := uuid.Parse(resourceToGet)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ErrInvalidValue)
	}

	startDate, _ := time.Parse(helpers.LayoutISO, c.Param("startDate"))
	endDate, _ := time.Parse(helpers.LayoutISO, c.Param("endDate"))

	daysDuration := endDate.Sub(startDate).Hours() / 24
	if daysDuration > 31 {
		return c.JSON(http.StatusUnprocessableEntity, ErrInvalidValue)
	}

	results, _ := h.svc.List(c, int(orgID), resourceUUIDToGet, startDate, endDate)
	returnSlice := make([]capPlanArrReq, 0)

	fmt.Println(returnSlice)
	return c.JSON(http.StatusOK, results)

}

func (h *HTTP) deleteByID(c echo.Context) error {
	orgID := int(c.Get("orgID").(int))
	//userID := c.Get("userID").(int)

	uuidStrSent := c.Param("capEntryUUID")

	uuidToCheck, err := uuid.Parse(uuidStrSent)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ErrInvalidValue)
	}

	existingCap, err := h.svc.ViewByUUID(c, uuidToCheck, orgID)
	if err != nil || existingCap.ID == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	err = h.svc.DeleteByID(c, existingCap.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInternal)
	}
	return c.NoContent(http.StatusOK)

}

func (h *HTTP) create(c echo.Context) error {
	orgID := c.Get("orgID").(int)
	userID := c.Get("userID").(int)

	capReq := new(capPlanArrReq)
	if err := c.Bind(capReq); err != nil {
		return err
	}

	var allEntries []model.CapacityPlanEntry
	for _, tempCapEntryReq := range capReq.CapacityPlans {
		tempCapEntry := model.CapacityPlanEntry{}
		tempProject := model.Project{}
		tempProject.UUID = tempCapEntryReq.ProjectID
		tempCapEntry.UUID = uuid.New()
		tempCapEntry.ResourceID = userID
		tempWorkDate, err := time.Parse(time.RFC3339, capReq.WorkDate)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ErrInvalidValue)
		}

		tempCapEntry.WorkDate = null.TimeFrom(tempWorkDate)
		tempCapEntry.WorkPercent = tempCapEntryReq.WorkPercent
		tempCapEntry.Project = &tempProject
		allEntries = append(allEntries, tempCapEntry)
	}

	fmt.Println(allEntries)

	created := h.svc.Create(c, orgID, allEntries)
	fmt.Println(created)

	return c.NoContent(http.StatusOK)
}
