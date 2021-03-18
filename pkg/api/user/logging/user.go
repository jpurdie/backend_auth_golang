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

func (ls *LogService) FetchByEmail(c echo.Context, email string) (u *authapi.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "FetchByEmail request", err,
			map[string]interface{}{
				"email": email,
				"took":  time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.FetchByEmail(c, email)
}

func (ls *LogService) FetchUserByID(c echo.Context, userID int) (u *authapi.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "FetchUserByID request", err,
			map[string]interface{}{
				"userID": userID,
				"took":   time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.FetchUserByID(c, userID)
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

// Change logging
func (ls *LogService) FetchByExternalID(c echo.Context, s string) (us *authapi.User, err error) {
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

const name = "user"
