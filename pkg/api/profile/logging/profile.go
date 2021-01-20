package logging

import (
	"time"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/profile"

	"github.com/labstack/echo/v4"
)

// New creates new password logging service
func New(svc profile.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	profile.Service
	logger authapi.Logger
}

const name = "profile"

// Change logging
func (ls *LogService) Change(c echo.Context, p authapi.Profile) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Change password request", err,
			map[string]interface{}{
				"req":  p,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, p)
}
