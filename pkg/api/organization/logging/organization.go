package logging

import (
	"time"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/organization"

	"github.com/labstack/echo/v4"
)

func New(svc organization.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	organization.Service
	logger authapi.Logger
}

const name = "organization"

// Create logging
func (ls *LogService) Create(c echo.Context, p authapi.Profile) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Organization request", err,
			map[string]interface{}{
				"req":  p,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, p)
}
