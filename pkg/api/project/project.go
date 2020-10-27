package project

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/model"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var (
	ErrNotFound = authapi.Error{CodeInt: http.StatusNotFound, Message: "Not found"}
)

func (p Project) ListStatuses(c echo.Context) ([]model.ProjectStatus, error) {
	return p.pdb.ListStatuses(p.db)
}

func (p Project) ListTypes(c echo.Context) ([]model.ProjectType, error) {
	return p.pdb.ListTypes(p.db)
}

func (p Project) ListComplexities(c echo.Context) ([]model.ProjectComplexity, error) {
	return p.pdb.ListComplexities(p.db)
}

func (p Project) ListSizes(c echo.Context) ([]model.ProjectSize, error) {
	return p.pdb.ListSizes(p.db)
}

func (p Project) List(c echo.Context, oID string) ([]model.Project, error) {
	return p.pdb.List(p.db, oID)
}

func (p Project) View(c echo.Context, pUUID string) (model.Project, error) {
	return p.pdb.View(p.db, pUUID)
}

func (p Project) UpdateName(c echo.Context, proj model.Project) error {
	op := "UpdateName"
	tempProj, err := p.pdb.View(p.db, proj.UUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}
	return p.pdb.UpdateName(p.db, proj)
}
func (p Project) UpdateDescr(c echo.Context, proj model.Project) error {
	op := "UpdateDescr"
	tempProj, err := p.pdb.View(p.db, proj.UUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}
	return p.pdb.UpdateDescr(p.db, proj)
}

func (p Project) UpdateSize(c echo.Context, pUUID uuid.UUID, ps model.ProjectSize) error {
	op := "UpdateSize"

	tempProj, err := p.pdb.View(p.db, pUUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}

	//get status ID
	sizes, _ := p.pdb.ListSizes(p.db)
	for _, tempS := range sizes {
		if tempS.UUID == ps.UUID {
			ps.ID = tempS.ID
			break
		}
	}
	if ps.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINVALID,
			Err:  err,
		}
	}

	return p.pdb.UpdateSize(p.db, pUUID, ps)

}
func (p Project) UpdateStatus(c echo.Context, pUUID uuid.UUID, ps model.ProjectStatus) error {
	op := "UpdateStatus"

	tempProj, err := p.pdb.View(p.db, pUUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}

	//get status ID
	statuses, _ := p.pdb.ListStatuses(p.db)
	for _, tempS := range statuses {
		if tempS.UUID == ps.UUID {
			ps.ID = tempS.ID
			break
		}
	}
	if ps.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINVALID,
			Err:  err,
		}
	}

	return p.pdb.UpdateStatus(p.db, pUUID, ps)

}

func (p Project) UpdateRGT(c echo.Context, pUUID uuid.UUID, rgt string) error {
	op := "UpdateRGT"

	tempProj, err := p.pdb.View(p.db, pUUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}
	return p.pdb.UpdateRGT(p.db, pUUID, rgt)

}


func (p Project) UpdateType(c echo.Context, pUUID uuid.UUID, pt model.ProjectType) error {
	op := "UpdateComplexity"

	tempProj, err := p.pdb.View(p.db, pUUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}
	//get type ID
	types, _ := p.pdb.ListTypes(p.db)
	for _, tempType := range types {
		if tempType.UUID == pt.UUID {
			pt.ID = tempType.ID
			break
		}
	}
	return p.pdb.UpdateType(p.db, pUUID, pt)
}

func (p Project) UpdateComplexity(c echo.Context, pUUID uuid.UUID, pc model.ProjectComplexity) error {
	op := "UpdateComplexity"

	tempProj, err := p.pdb.View(p.db, pUUID.String())
	if err != nil {
		return &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  err,
		}
	}
	if tempProj.ID == 0 {
		return &authapi.Error{
			Op:   op,
			Code: authapi.ENOTFOUND,
			Err:  err,
		}
	}
	//get complexity ID
	complexities, _ := p.pdb.ListComplexities(p.db)
	for _, tempComplexity := range complexities {
		if tempComplexity.UUID == pc.UUID {
			pc.ID = tempComplexity.ID
			break
		}
	}
	return p.pdb.UpdateComplexity(p.db, pUUID, pc)
}

func (p Project) Create(c echo.Context, myProj model.Project) (model.Project, error) {

	//checking if status exists
	if myProj.StatusID == 0 {
		statuses, _ := p.pdb.ListStatuses(p.db)
		for _, tempStatus := range statuses {
			if tempStatus.UUID == myProj.Status.UUID {
				myProj.StatusID = tempStatus.ID
				break
			}
		}
	}

	//checking if type exists
	if myProj.TypeID == 0 {
		types, _ := p.pdb.ListTypes(p.db)
		for _, tempType := range types {
			if tempType.UUID == myProj.Type.UUID {
				myProj.TypeID = tempType.ID
				break
			}
		}
	}

	//checking if complexity exists
	if myProj.ComplexityID == 0 {
		complexities, _ := p.pdb.ListComplexities(p.db)
		for _, tempComplexity := range complexities {
			if tempComplexity.UUID == myProj.Complexity.UUID {
				myProj.ComplexityID = tempComplexity.ID
				break
			}
		}
	}

	//checking if size exists
	if myProj.SizeID == 0 {
		sizes, _ := p.pdb.ListSizes(p.db)
		for _, tempSize := range sizes {
			if tempSize.UUID == myProj.Size.UUID {
				myProj.SizeID = tempSize.ID
				break
			}
		}
	}

	log.Println(myProj)

	myProj.UUID = uuid.New()
	return p.pdb.Create(p.db, myProj)
}
