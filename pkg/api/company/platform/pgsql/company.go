package pgsql

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"net/http"
	"strings"

	"github.com/jpurdie/authapi"
)

// User represents the client for company_user table
type CompanyUser struct{}

// Custom errors
var (
	ErrCompAlreadyExists  = echo.NewHTTPError(http.StatusConflict, "Company name already exists.")
	ErrEmailAlreadyExists = echo.NewHTTPError(http.StatusConflict, "Email already exists.")
)

// Create creates a new user on database
func (c CompanyUser) Create(db *pg.DB, cu authapi.CompanyUser) (authapi.CompanyUser, error) {

	var n string
	tempDb, _ := postgres.Init()
	_, err := tempDb.QueryOne(pg.Scan(&n), "SELECT now() ")
	tempDb.Close()

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database at: " + n)

	var company = new(authapi.Company)
	tempDb, _ = postgres.Init()
	count, err := tempDb.Model(company).Where("lower(name) = ? and deleted_at is null", strings.ToLower(cu.Company.Name)).Count()
	tempDb.Close()
	if err != nil {
		return authapi.CompanyUser{}, err
	}
	if count > 0 {
		return authapi.CompanyUser{}, ErrCompAlreadyExists
	}
	var user = new(authapi.User)

	tempDb, _ = postgres.Init()
	count, err = tempDb.Model(user).Where("lower(email) = ? and deleted_at is null", strings.ToLower(cu.User.Email)).Count()
	tempDb.Close()

	if err != nil {
		return authapi.CompanyUser{}, err
	}
	if count > 0 {
		return authapi.CompanyUser{}, ErrEmailAlreadyExists
	}

	print(db.PoolStats().TotalConns)
	tempDb, _ = postgres.Init()
	tx, err := tempDb.Begin()
	tx.Model(cu.Company).Insert()
	cu.User.CompanyID = cu.Company.ID
	tx.Model(cu.User).Insert()
	cu.UserID = cu.User.ID
	cu.CompanyID = cu.Company.ID
	tx.Model(&cu).Insert()
	trErr := tx.Commit()
	if trErr != nil {
		tx.Rollback()
	}
	tempDb.Close()
	return cu, err
}

//// View returns single user by ID
//func (co Company) View(db orm.DB, id int) (authapi.Company, error) {
//	var company authapi.Company
//	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name"
//	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id"
//	WHERE ("user"."id" = ? and deleted_at is null)`
//	_, err := db.QueryOne(&company, sql, id)
//	return company, err
//}
//
//// Update updates user's contact info
//func (co Company) Update(db orm.DB, company authapi.Company) error {
//	_, err := db.Model(&company).WherePK().UpdateNotZero()
//	return err
//}
//
//// List returns list of all users retrievable for the current user, depending on role
//func (co Company) List(db orm.DB, qp *authapi.ListQuery, p authapi.Pagination) ([]authapi.Company, error) {
//	var companies []authapi.Company
//	q := db.Model(&companies).Relation("Role").Limit(p.Limit).Offset(p.Offset).Where("deleted_at is null").Order("user.id desc")
//	if qp != nil {
//		q.Where(qp.Query, qp.ID)
//	}
//	err := q.Select()
//	return companies, err
//}
//
//// Delete sets deleted_at for a user
//func (co Company) Delete(db orm.DB, company authapi.Company) error {
//	return db.Delete(&company)
//}
