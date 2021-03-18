package logging

import (
	"time"

	"github.com/google/uuid"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/capacityplan"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

func New(svc capacityplan.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents logging service
type LogService struct {
	capacityplan.Service
	logger authapi.Logger
}

const name = "capacityplan"

func (ls *LogService) ViewByUUID(c echo.Context, capEntryUUID uuid.UUID, orgID int) (capEntry model.CapacityPlanEntry, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "ViewByID Cap Entry request", err,
			map[string]interface{}{
				"capEntryUUID": capEntryUUID.String(),
				"orgID":        orgID,
				"took":         time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ViewByUUID(c, capEntryUUID, orgID)
}

func (ls *LogService) ViewByID(c echo.Context, capEntryID int, orgID int) (capEntry model.CapacityPlanEntry, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "ViewByID Cap Entry request", err,
			map[string]interface{}{
				"capEntryID": capEntryID,
				"orgID":      orgID,
				"took":       time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ViewByID(c, capEntryID, orgID)
}

func (ls *LogService) DeleteByID(c echo.Context, capEntryID int) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "DeleteByID CapEntry request", err,
			map[string]interface{}{
				"capEntryID": capEntryID,
				"took":       time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.DeleteByID(c, capEntryID)
}

func (ls *LogService) Create(c echo.Context, orgID int, cp []model.CapacityPlanEntry) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create CapacityPlanEntry request", err,
			map[string]interface{}{
				"req":   cp,
				"orgID": orgID,
				"took":  time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, orgID, cp)
}

func (ls *LogService) List(c echo.Context, orgID int, resourceToGet uuid.UUID, startDate time.Time, endDate time.Time) (cpes []model.CapacityPlanEntry, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List CapacityPlanEntry request", err,
			map[string]interface{}{
				"resourceToGet": resourceToGet,
				"startDate":     startDate,
				"endDate":       endDate,
				"orgID":         orgID,
				"took":          time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, orgID, resourceToGet, startDate, endDate)
}
func (ls *LogService) ListSummary(c echo.Context, orgID int, resourceToGet uuid.UUID, startDate time.Time, endDate time.Time) (cpes []model.CapacityPlan, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Summary CapacityPlanEntry request", err,
			map[string]interface{}{
				"resourceToGet": resourceToGet,
				"startDate":     startDate,
				"endDate":       endDate,
				"orgID":         orgID,
				"took":          time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ListSummary(c, orgID, resourceToGet, startDate, endDate)
}
