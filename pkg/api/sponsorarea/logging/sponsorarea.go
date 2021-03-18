package logging

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	sponsorarea "github.com/jpurdie/authapi/pkg/api/sponsorarea"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

func New(svc sponsorarea.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	sponsorarea.Service
	logger authapi.Logger
}

const name = "sponsorarea"

func (ls *LogService) List(c echo.Context, oID int) (sa []model.SponsorArea, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Sponsor Area request", err,
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
			name, "Delete Sponsor Area request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, saUUID, orgID)
}

func (ls *LogService) Create(c echo.Context, alignmentName string, orgID int) (sa model.SponsorArea, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Sponsor Area request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, alignmentName, orgID)
}

func (ls *LogService) Update(c echo.Context, mySA model.SponsorArea) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Update Sponsor Area request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, mySA)
}
