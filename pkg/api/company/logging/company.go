package company

import (
	"time"

	"github.com/labstack/echo"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/company"
)

// New creates new Company logging service
func New(svc company.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents Company logging service
type LogService struct {
	company.Service
	logger authapi.Logger
}

const name = "company"

// Create logging
func (ls *LogService) Create(c echo.Context, req authapi.CompanyUser) (resp authapi.CompanyUser, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create company request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}

//// List logging
//func (ls *LogService) List(c echo.Context, req authapi.Pagination) (resp []authapi.Company, err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "List Company request", err,
//			map[string]interface{}{
//				"req":  req,
//				"resp": resp,
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.List(c, req)
//}
//
//// View logging
//func (ls *LogService) View(c echo.Context, req int) (resp authapi.Company, err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "View Company request", err,
//			map[string]interface{}{
//				"req":  req,
//				"resp": resp,
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.View(c, req)
//}
//
//// Delete logging
//func (ls *LogService) Delete(c echo.Context, req int) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Delete Company request", err,
//			map[string]interface{}{
//				"req":  req,
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.Delete(c, req)
//}
//
//// Update logging
//func (ls *LogService) Update(c echo.Context, req company.Update) (resp authapi.Company, err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update Company request", err,
//			map[string]interface{}{
//				"req":  req,
//				"resp": resp,
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.Update(c, req)
//}
