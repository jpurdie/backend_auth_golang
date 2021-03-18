package project

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound = authapi.Error{CodeInt: http.StatusNotFound, Message: "Not found"}
)

func (p Project) ListStatuses(c echo.Context) ([]model.ProjectStatus, error) {
	return p.pdb.ListStatuses(*p.db)
}

func (p Project) ListTypes(c echo.Context) ([]model.ProjectType, error) {
	return p.pdb.ListTypes(*p.db)
}

func (p Project) ListComplexities(c echo.Context) ([]model.ProjectComplexity, error) {
	return p.pdb.ListComplexities(*p.db)
}

func (p Project) ListSizes(c echo.Context) ([]model.ProjectSize, error) {
	return p.pdb.ListSizes(*p.db)
}

func (p Project) List(c echo.Context, oID int, filters map[string]string) ([]model.Project, error) {
	return p.pdb.List(*p.db, oID, filters)
}

func (p Project) View(c echo.Context, pUUID uuid.UUID) (model.Project, error) {
	return p.pdb.View(*p.db, pUUID)
}

func (p Project) Update(c echo.Context, orgID int, pUUID uuid.UUID, key string, val interface{}) error {
	op := "Update"

	proj, err := p.pdb.View(*p.db, pUUID)

	if proj.ID == 0 || err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err}
	}

	switch key {
	case "name":
		str := fmt.Sprintf("%v", val)
		if len(str) < 2 {
			return &authapi.Error{
				Op:   op,
				Code: authapi.EINVALID,
			}
		}
		return p.pdb.Update(*p.db, orgID, pUUID, "name", str)
	case "description":
		str := strings.ToUpper(fmt.Sprintf("%v", val))
		return p.pdb.Update(*p.db, orgID, pUUID, "description", str)
	case "rgt":
		str := strings.ToUpper(fmt.Sprintf("%v", val))
		if str != "R" && str != "G" && str != "T" {
			return &authapi.Error{
				Op:   op,
				Code: authapi.ENOTFOUND,
			}
		}
		return p.pdb.Update(*p.db, orgID, pUUID, "rgt", str)
	case "status":
		str := fmt.Sprintf("%v", val)
		pStatUUID, err := uuid.Parse(str)
		if err != nil {
			return &authapi.Error{
				Op:   op,
				Code: authapi.ENOTFOUND,
				Err:  err}
		}
		statuses, _ := p.pdb.ListStatuses(*p.db)
		for _, tempStat := range statuses {
			if tempStat.UUID == pStatUUID {
				return p.pdb.Update(*p.db, orgID, pUUID, "status", tempStat.ID)
			}
		}

	case "type":
		str := fmt.Sprintf("%v", val)
		pTypeUUID, err := uuid.Parse(str)
		if err != nil {
			return &authapi.Error{
				Op:   op,
				Code: authapi.ENOTFOUND,
				Err:  err}
		}
		types, _ := p.pdb.ListTypes(*p.db)
		for _, tempType := range types {
			if tempType.UUID == pTypeUUID {
				return p.pdb.Update(*p.db, orgID, pUUID, "type", tempType.ID)
			}
		}
	case "complexity":
		str := fmt.Sprintf("%v", val)
		pCompUUID, err := uuid.Parse(str)
		if err != nil {
			return &authapi.Error{
				Op:   op,
				Code: authapi.ENOTFOUND,
				Err:  err}
		}
		complexities, _ := p.pdb.ListComplexities(*p.db)
		for _, tempComplexity := range complexities {
			if tempComplexity.UUID == pCompUUID {
				return p.pdb.Update(*p.db, orgID, pUUID, "complexity", tempComplexity.ID)
			}
		}
	case "size":
		str := fmt.Sprintf("%v", val)
		pStatUUID, err := uuid.Parse(str)
		if err != nil {
			return &authapi.Error{
				Op:   op,
				Code: authapi.ENOTFOUND,
				Err:  err}
		}
		sizes, _ := p.pdb.ListSizes(*p.db)
		for _, tempSize := range sizes {
			if tempSize.UUID == pStatUUID {
				return p.pdb.Update(*p.db, orgID, pUUID, "size", tempSize.ID)
			}
		}
	case "openForTimeEntry":
		if val == true || val == false {
			return p.pdb.Update(*p.db, orgID, pUUID, "openForTimeEntry", val)
		}
	case "timeConstrained":
		if val == true || val == false {
			return p.pdb.Update(*p.db, orgID, pUUID, "timeConstrained", val)
		}
	case "compliance":
		if val == true || val == false {
			return p.pdb.Update(*p.db, orgID, pUUID, "compliance", val)
		}
	case "strategicAlignments":

		strategicAlignmentsSent := val.([]interface{})
		var strategicAlignmentsIDsToSave []int
		strategicAlignments, _ := p.saDB.List(*p.db, orgID)
		for _, saFromReq := range strategicAlignmentsSent {
			for _, saFromDB := range strategicAlignments {
				tempSAToSaveStr := fmt.Sprintf("%v", saFromReq)
				if saFromDB.UUID.String() == tempSAToSaveStr {
					strategicAlignmentsIDsToSave = append(strategicAlignmentsIDsToSave, int(saFromDB.ID))
					break
				}
			}
		}
		return p.pdb.UpdateStrategicAlignments(*p.db, orgID, int(proj.ID), strategicAlignmentsIDsToSave)

	case "sponsorAreas":
		sponsorAreasSent := val.([]interface{})
		var sponsorAreaIDsToSave []int
		sponsorAreasFromDB, _ := p.sponsAreaDB.List(*p.db, orgID)
		for _, saFromReq := range sponsorAreasSent {
			for _, saFromDB := range sponsorAreasFromDB {
				tempSAToSaveStr := fmt.Sprintf("%v", saFromReq)
				if saFromDB.UUID.String() == tempSAToSaveStr {
					sponsorAreaIDsToSave = append(sponsorAreaIDsToSave, int(saFromDB.ID))
					break
				}
			}
		}
		return p.pdb.UpdateSponsorAreas(*p.db, orgID, int(proj.ID), sponsorAreaIDsToSave)
	default:

	}
	return &authapi.Error{
		Op:   op,
		Code: authapi.ENOTFOUND,
	}
}

func (p Project) Create(c echo.Context, myProj model.Project) (uuid.UUID, error) {
	op := "create"
	//checking if status exists
	if myProj.StatusID == 0 {
		statuses, _ := p.pdb.ListStatuses(*p.db)
		for i, _ := range statuses {
			if statuses[i].UUID == myProj.Status.UUID {
				myProj.StatusID = statuses[i].ID
				break
			}
		}
	}

	//checking if type exists
	if myProj.TypeID == 0 {
		types, _ := p.pdb.ListTypes(*p.db)
		for _, tempType := range types {
			if tempType.UUID == myProj.Type.UUID {
				myProj.TypeID = tempType.ID
				break
			}
		}
	}

	//checking if complexity exists
	if myProj.ComplexityID == 0 {
		complexities, _ := p.pdb.ListComplexities(*p.db)
		for _, tempComplexity := range complexities {
			if tempComplexity.UUID == myProj.Complexity.UUID {
				myProj.ComplexityID = tempComplexity.ID
				break
			}
		}
	}

	//checking if size exists
	if myProj.SizeID == 0 {
		sizes, _ := p.pdb.ListSizes(*p.db)
		for _, tempSize := range sizes {
			if tempSize.UUID == myProj.Size.UUID {
				myProj.SizeID = tempSize.ID
				break
			}
		}
	}

	savedSAs, err := p.saDB.List(*p.db, int(myProj.OrganizationID))
	if err != nil {
		return uuid.Nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	for _, tempSA := range savedSAs {
		for _, tempProjSA := range myProj.StrategicAlignments {
			if tempSA.UUID == tempProjSA.UUID {
				tempProjSA.ID = tempSA.ID
				continue
			}
		}
	}

	savedSponsorAreas, err := p.sponsAreaDB.List(*p.db, int(myProj.OrganizationID))
	if err != nil {
		return uuid.Nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	for _, tempSA := range savedSponsorAreas {
		for _, tempProjSA := range myProj.SponsorAreas {
			if tempSA.UUID == tempProjSA.UUID {
				tempProjSA.ID = tempSA.ID
				continue
			}
		}
	}

	myProj.UUID = uuid.New()

	err = p.pdb.Create(*p.db, myProj)
	if err != nil {
		return uuid.Nil, &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	return myProj.UUID, nil

}
