package logging

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/project"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

func New(svc project.Service, logger authapi.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	project.Service
	logger authapi.Logger
}

const name = "project"

//func (ls *LogService) Update(c echo.Context, key string, p model.Project) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update Project Request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.Update(c, key, p)
//}

//
//func (ls *LogService) UpdateName(c echo.Context, proj model.Project) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update proj name request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateName(c, proj)
//}
//func (ls *LogService) UpdateDescr(c echo.Context, proj model.Project) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update proj descr request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateDescr(c, proj)
//}
//
//func (ls *LogService) UpdateRGT(c echo.Context, pUUID uuid.UUID, rgt string) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update rgt request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateRGT(c, pUUID, rgt)
//}
//func (ls *LogService) UpdateSize(c echo.Context, pUUID uuid.UUID, ps model.ProjectSize) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update project status request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateSize(c, pUUID, ps)
//}

//func (ls *LogService) UpdateStatus(c echo.Context, pUUID uuid.UUID, ps model.ProjectStatus) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update project status request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateStatus(c, pUUID, ps)
//}
//
//func (ls *LogService) UpdateType(c echo.Context, pUUID uuid.UUID, pt model.ProjectType) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update project type request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateType(c, pUUID, pt)
//}
//
//func (ls *LogService) UpdateComplexity(c echo.Context, pUUID uuid.UUID, pc model.ProjectComplexity) (err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "Update project request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.UpdateComplexity(c, pUUID, pc)
//}

func (ls *LogService) Create(c echo.Context, myProj model.Project) (projUUID uuid.UUID, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create project request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, myProj)
}

func (ls *LogService) ListStatuses(c echo.Context) (ps []model.ProjectStatus, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Statuses request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ListStatuses(c)
}

func (ls *LogService) ListTypes(c echo.Context) (pt []model.ProjectType, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Project Types request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ListTypes(c)
}

func (ls *LogService) ListComplexities(c echo.Context) (pt []model.ProjectComplexity, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Project Complexities request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ListComplexities(c)
}

func (ls *LogService) ListSizes(c echo.Context) (pt []model.ProjectSize, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Project Sizes request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ListSizes(c)
}

func (ls *LogService) List(c echo.Context, oID int, filters map[string]string) (pt []model.Project, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "List Projects request", err,
			map[string]interface{}{
				"took":    time.Since(begin),
				"filters": filters,
			},
		)
	}(time.Now())
	return ls.Service.List(c, oID, filters)
}

//func (ls *LogService) View(c echo.Context, proj model.Project) (pt model.Project, err error) {
//	defer func(begin time.Time) {
//		ls.logger.Log(
//			c,
//			name, "View Project request", err,
//			map[string]interface{}{
//				"took": time.Since(begin),
//			},
//		)
//	}(time.Now())
//	return ls.Service.View(c, proj)
//}
