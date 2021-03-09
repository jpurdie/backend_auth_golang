package logging

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/user"
	"github.com/labstack/echo/v4"
)

// New creates new password logging service
func New(svc user.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	user.Service
	logger authapi.Logger
}

const name = "user"

// Change logging
func (ls *LogService) FetchByExternalID(c echo.Context, s string) (us authapi.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Fetch user request", err,
			map[string]interface{}{
				"req":  s,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.FetchByExternalID(c, s)
}

func (ls *LogService) ListRoles(c echo.Context) (roles []authapi.Role, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "ListRoles request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ListRoles(c)
}

func (ls *LogService) List(c echo.Context, orgID int) (users []authapi.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List request", err,
			map[string]interface{}{
				"req":  orgID,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, orgID)
}

func (ls *LogService) UpdateRole(c echo.Context, level int, profileID int) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "UpdateRole request", err,
			map[string]interface{}{
				"level":     level,
				"profileID": profileID,
				"took":      time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.UpdateRole(c, level, profileID)
}

func (ls *LogService) FetchProfile(c echo.Context, userID int, orgID int) (p authapi.Profile, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "FetchProfile request", err,
			map[string]interface{}{
				"userID": userID,
				"orgID":  orgID,
				"took":   time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.FetchProfile(c, userID, orgID)
}

func (ls *LogService) FetchUserByUUID(c echo.Context, userUUID uuid.UUID, orgID int) (u *authapi.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "FetchUserByUUID request", err,
			map[string]interface{}{
				"userUUID": userUUID,
				"orgID":    orgID,
				"took":     time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.FetchUserByUUID(c, userUUID, orgID)
}
