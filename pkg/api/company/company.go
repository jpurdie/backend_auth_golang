// Package user contains user application services
package company

import (
	"github.com/labstack/echo"

	"github.com/jpurdie/authapi"
)

// Create creates a new CompanyUser account
func (coU CompanyUser) Create(c echo.Context, req authapi.CompanyUser) (authapi.CompanyUser, error) {
	return coU.cdb.Create(coU.db, req)
}

// List returns list of users
//func (co Company) List(c echo.Context, p authapi.Pagination) ([]authapi.Company, error) {
//	au := co.rbac.User(c)
//	q, err := query.List(au)
//	if err != nil {
//		return nil, err
//	}
//	return co.udb.List(co.db, q, p)
//}

// View returns single user
//func (co Company) View(c echo.Context, id int) (authapi.Company, error) {
//	if err := co.rbac.EnforceUser(c, id); err != nil {
//		return authapi.Company{}, err
//	}
//	return co.udb.View(co.db, id)
//}

// Delete deletes a user
//func (co Company) Delete(c echo.Context, id int) error {
//	user, err := co.udb.View(co.db, id)
//	if err != nil {
//		return err
//	}
//	//if err := co.rbac.IsLowerRole(c, authapi.User.AccessLevel); err != nil {
//	//	return err
//	//}
//	return co.udb.Delete(co.db, user)
//}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName string
	LastName  string
	Mobile    string
	Phone     string
	Address   string
}

// Update updates user's contact information
//func (co Company) Update(c echo.Context, r Update) (authapi.User, error) {
//	if err := co.rbac.EnforceUser(c, r.ID); err != nil {
//		return authapi.User{}, err
//	}
//
//	if err := co.udb.Update(co.db, authapi.User{
//		Base:      authapi.Base{ID: r.ID},
//		FirstName: r.FirstName,
//		LastName:  r.LastName,
//		Mobile:    r.Mobile,
//		Address:   r.Address,
//	}); err != nil {
//		return authapi.User{}, err
//	}
//
//	return co.udb.View(co.db, r.ID)
//}
