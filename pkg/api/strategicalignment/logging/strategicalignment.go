package logging

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/strategicalignment"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

func New(svc strategicalignment.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	strategicalignment.Service
	logger authapi.Logger
}

const name = "strategicalignment"

func (ls *LogService) List(c echo.Context, oID int) (sa []model.StrategicAlignment, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Strategic Alignment request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, oID)
}

func (ls *LogService) Delete(c echo.Context, saUUID uuid.UUID, orgID int) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Delete Strategic Alignment request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, saUUID, orgID)
}

func (ls *LogService) Create(c echo.Context, alignmentName string, orgID int) (sa model.StrategicAlignment, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Strategic Alignment request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, alignmentName, orgID)
}

func (ls *LogService) Update(c echo.Context, mySA model.StrategicAlignment) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Update Strategic Alignment request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, mySA)
}
