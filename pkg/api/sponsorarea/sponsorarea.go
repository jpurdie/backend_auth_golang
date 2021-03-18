package sponsorarea

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

func (sa SponsorArea) Update(c echo.Context, mySA model.SponsorArea) error {
	return sa.sadb.Update(*sa.db, mySA)
}

func (sa SponsorArea) List(c echo.Context, oID int) ([]model.SponsorArea, error) {
	return sa.sadb.List(*sa.db, oID)
}

func (sa SponsorArea) Delete(c echo.Context, saUUID uuid.UUID, orgID int) error {
	return sa.sadb.Delete(*sa.db, saUUID, orgID)
}

func (sa SponsorArea) Create(c echo.Context, saName string, orgID int) (model.SponsorArea, error) {
	mySA := model.SponsorArea{}
	mySA.Name = saName
	mySA.UUID = uuid.New()
	mySA.OrganizationID = int(orgID)
	return sa.sadb.Create(*sa.db, mySA)
}
