package strategicalignment

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound = authapi.Error{CodeInt: http.StatusNotFound, Message: "Not found"}
)

func (sa StrategicAlignment) Update(c echo.Context, mySA model.StrategicAlignment) error {
	return sa.sadb.Update(*sa.db, mySA)
}

func (sa StrategicAlignment) List(c echo.Context, oID int) ([]model.StrategicAlignment, error) {
	return sa.sadb.List(*sa.db, oID)
}

func (sa StrategicAlignment) Delete(c echo.Context, saUUID uuid.UUID, orgID int) error {
	return sa.sadb.Delete(*sa.db, saUUID, orgID)
}

func (sa StrategicAlignment) Create(c echo.Context, alignmentName string, orgID int) (model.StrategicAlignment, error) {
	mySA := model.StrategicAlignment{}
	mySA.Name = alignmentName
	mySA.UUID = uuid.New()
	mySA.OrganizationID = int(orgID)
	return sa.sadb.Create(*sa.db, mySA)
}
