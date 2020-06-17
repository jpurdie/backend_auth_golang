package mockdb

import (
	"github.com/go-pg/pg/v9/orm"

	"github.com/jpurdie/authapi"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, authapi.User) (authapi.User, error)
	ViewFn           func(orm.DB, int) (authapi.User, error)
	FindByUsernameFn func(orm.DB, string) (authapi.User, error)
	FindByTokenFn    func(orm.DB, string) (authapi.User, error)
	ListFn           func(orm.DB, *authapi.ListQuery, authapi.Pagination) ([]authapi.User, error)
	DeleteFn         func(orm.DB, authapi.User) error
	UpdateFn         func(orm.DB, authapi.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr authapi.User) (authapi.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (authapi.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (authapi.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (authapi.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *authapi.ListQuery, p authapi.Pagination) ([]authapi.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr authapi.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr authapi.User) error {
	return u.UpdateFn(db, usr)
}
