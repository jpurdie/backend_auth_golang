package database

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
	"strings"
)

// User represents the client for company_user table

type CompanyStore struct {
	db *pg.DB
}

// NewAdmAccountStore returns an AccountStore.
func NewCompanyStore(db *pg.DB) *CompanyStore {
	return &CompanyStore{
		db: db,
	}
}

//
// Custom errors
var (
	ErrCompAlreadyExists  = errors.New("Company name already exists")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

// Create creates a new user on database
func (s *CompanyStore) Create(cu authapi.CompanyUser) (authapi.CompanyUser, error) {

	var company = new(authapi.Company)

	count, err := s.db.Model(company).Where("lower(name) = ? and deleted_at is null", strings.ToLower(cu.Company.Name)).Count()
	if err != nil {
		return authapi.CompanyUser{}, err
	}
	if count > 0 {
		return authapi.CompanyUser{}, ErrCompAlreadyExists
	}
	var user = new(authapi.User)

	count, err = s.db.Model(user).Where("lower(email) = ? and deleted_at is null", strings.ToLower(cu.User.Email)).Count()

	if err != nil {
		return authapi.CompanyUser{}, err
	}
	if count > 0 {
		return authapi.CompanyUser{}, ErrEmailAlreadyExists
	}

	tx, err := s.db.Begin()
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
	return cu, err
}

//func (env *Env) List() ([]authapi.Company, error) {
//	var companies []authapi.Company
//	tempDb, _ := postgres.Init()
//	defer tempDb.Close()
//
//	_, err := tempDb.QueryOne(pg.Scan(&n), "SELECT now() ")
//	return companies, err
//}

//
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
//
//// Delete sets deleted_at for a user
//func (co Company) Delete(db orm.DB, company authapi.Company) error {
//	return db.Delete(&company)
//}
