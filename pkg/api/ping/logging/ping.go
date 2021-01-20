package logging

import (
	"time"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/ping"

	"github.com/labstack/echo/v4"
)

func New(svc ping.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	ping.Service
	logger authapi.Logger
}

const name = "ping"

// Change logging
func (ls *LogService) Create(c echo.Context, p authapi.Ping) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Ping request", err,
			map[string]interface{}{
				"req":  p,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, p)
}
