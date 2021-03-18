package capacityplan

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jpurdie/authapi/pkg/utl/helpers"

	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

func (cp CapacityPlan) DeleteByID(c echo.Context, capEntryID int) error {
	op := "DeleteByID"
	fmt.Println(op)

	err := cp.cpdb.DeleteByID(*cp.db, capEntryID)
	if err != nil {
		return err
	}
	return nil
}

func (cp CapacityPlan) ViewByUUID(c echo.Context, capEntryUUID uuid.UUID, orgID int) (model.CapacityPlanEntry, error) {
	op := "ViewByUUID"
	fmt.Println(op)

	capEntry, err := cp.cpdb.ViewByUUID(*cp.db, capEntryUUID, orgID)
	if err != nil {
		return model.CapacityPlanEntry{}, err
	}
	return capEntry, nil
}

func (cp CapacityPlan) ViewByID(c echo.Context, capEntryID int, orgID int) (model.CapacityPlanEntry, error) {
	op := "ViewByID"
	fmt.Println(op)

	capEntry, err := cp.cpdb.ViewByID(*cp.db, capEntryID, orgID)
	if err != nil {
		return model.CapacityPlanEntry{}, err
	}
	return capEntry, nil
}

func (cp CapacityPlan) Create(c echo.Context, orgID int, sentCapacityPlans []model.CapacityPlanEntry) error {
	op := "Create"
	fmt.Println(op)

	var distinctUUIDs []string
	capEntrySumMap := make(map[string]int)
	capEntriesToSave := make([]model.CapacityPlanEntry, 0)

	for i, tempCapEntry := range sentCapacityPlans {
		tempProj, err := cp.projDB.View(*cp.db, tempCapEntry.Project.UUID)
		if err != nil {
			return err
		}
		sentCapacityPlans[i].ProjectID = tempProj.ID
		if !helpers.StringContains(distinctUUIDs, tempCapEntry.Project.UUID.String()) {
			//distinct project entry
			distinctUUIDs = append(distinctUUIDs, tempCapEntry.Project.UUID.String())
		}
	}

	//loop through each distinct project UUID
	for _, projUUID := range distinctUUIDs {
		//sum the total for the project. This is in case the client sends the same project on the same day for the same person.
		for _, tempCapEntry := range sentCapacityPlans {
			if tempCapEntry.Project.UUID.String() == projUUID {
				if _, ok := capEntrySumMap[projUUID]; ok {
					capEntrySumMap[projUUID] = capEntrySumMap[projUUID] + tempCapEntry.WorkPercent
				} else {
					capEntrySumMap[projUUID] = tempCapEntry.WorkPercent
				}
			}
		}
		// we now have the sum of each project-workdate-resource
	}
	fmt.Println(capEntrySumMap)

	for _, projUUID := range distinctUUIDs {
		for _, tempCapEntry := range sentCapacityPlans {
			if tempCapEntry.Project.UUID.String() == projUUID {
				//setting the value to the sum
				tempCapEntry.WorkPercent = capEntrySumMap[projUUID]
				capEntriesToSave = append(capEntriesToSave, tempCapEntry)
				// we only want to add one of each project|resource|workdate
				break
			}

		}
	}
	fmt.Println(capEntriesToSave)

	if len(capEntriesToSave) > 0 {
		err := cp.cpdb.Create(*cp.db, orgID, capEntriesToSave)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp CapacityPlan) List(c echo.Context, orgID int, resourceToGet uuid.UUID, startDate time.Time, endDate time.Time) ([]model.CapacityPlanEntry, error) {
	op := "List"
	fmt.Println(op)

	user, err := cp.udb.FetchUserByUUID(*cp.db, resourceToGet, orgID)
	if err != nil || user.ID == 0 {
		return nil, err
	}

	planEntries, err := cp.cpdb.List(*cp.db, orgID, user.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return planEntries, nil
}

func (cp CapacityPlan) ListSummary(c echo.Context, orgID int, resourceToGet uuid.UUID, startDate time.Time, endDate time.Time) ([]model.CapacityPlan, error) {
	op := "ListSummary"

	fmt.Println(op)

	returnPlans := make([]model.CapacityPlan, 0)

	user, err := cp.udb.FetchUserByUUID(*cp.db, resourceToGet, orgID)
	if err != nil || user.ID == 0 {
		return nil, err
	}

	for d := startDate; d.After(endDate) == false; d = d.AddDate(0, 0, 1) {
		daySum := 0
		planEntries, err := cp.cpdb.List(*cp.db, orgID, user.ID, d, d)
		if err != nil {
			return nil, err
		}
		for _, planEntry := range planEntries {
			daySum += planEntry.WorkPercent
		}
		//, _ := time.LoadLocation("America/Phoenix")
		tempPlan := model.CapacityPlan{}
		tempPlan.WorkDate = d
		tempPlan.SumWorkPercent = daySum
		tempPlan.ResourceUUID = nil
		returnPlans = append(returnPlans, tempPlan)
	}

	if err != nil {
		return nil, err
	}
	return returnPlans, nil
}
