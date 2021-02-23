package logging

import (
	"time"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/invitation"
	"github.com/labstack/echo/v4"
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

func (ls *LogService) View(c echo.Context, tokenPlainText string) (i authapi.Invitation, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Invitation request", err,
			map[string]interface{}{
				"req":  tokenPlainText,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, tokenPlainText)
}

func (ls *LogService) CreateUser(c echo.Context, profile authapi.Profile, invite authapi.Invitation) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Invitation request", err,
			map[string]interface{}{
				"profile": profile,
				"invite":  invite,
				"took":    time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.CreateUser(c, profile, invite)
}

func (ls *LogService) List(c echo.Context, orgID int, includeExpired bool, includeUsed bool) (invites []authapi.Invitation, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create Invitation request", err,
			map[string]interface{}{
				"orgID":          orgID,
				"includeExpired": includeExpired,
				"includeUsed":    includeUsed,
				"took":           time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, orgID, includeExpired, includeUsed)
}

func (ls *LogService) Delete(c echo.Context, email string, orgID int) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Delete Invitation request", err,
			map[string]interface{}{
				"email": email,
				"orgID": orgID,
				"took":  time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, email, orgID)
}
