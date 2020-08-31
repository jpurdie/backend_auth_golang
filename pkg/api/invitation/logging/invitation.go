package logging

import (
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/invitation"
	"github.com/labstack/echo"
	"time"
)

func New(svc invitation.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	invitation.Service
	logger authapi.Logger
}

const name = "invitation"
func (ls *LogService) Create(c echo.Context, i authapi.Invitation) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Invitation request", err,
			map[string]interface{}{
				"req":  i,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, i)
}

func (ls *LogService) View(c echo.Context, tokenHash string) (i authapi.Invitation, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Invitation request", err,
			map[string]interface{}{
				"req":  tokenHash,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, tokenHash)
}


func (ls *LogService) CreateUser(c echo.Context, profile authapi.Profile, invite authapi.Invitation) ( err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Invitation request", err,
			map[string]interface{}{
				"profile":  profile,
				"invite":  invite,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.CreateUser(c, profile, invite)
}