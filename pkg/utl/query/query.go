package query

import (
	"github.com/labstack/echo"

	"github.com/jpurdie/authapi"
)

// List prepares data for list queries
func List(u authapi.AuthUser) (*authapi.ListQuery, error) {
	switch true {
	case u.Role <= authapi.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == authapi.CompanyAdminRole:
		return &authapi.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	//case u.Role == authapi.LocationAdminRole:
	//	return &authapi.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
