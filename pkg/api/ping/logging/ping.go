package ping

import (
	"github.com/labstack/echo"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/ping"
	"time"
)

func New(svc ping.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents auth logging service
type LogService struct {
	ping.Service
	logger authapi.Logger
}

const name = "ping"

// Ping Pong logging
func (ls *LogService) Create(c echo.Context, req int) (resp string, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Ping Pong create request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}
