package authapi

import (
	"database/sql"
	"github.com/google/uuid"
)

type User struct {
	Base
	FirstName      string         `json:"firstName,omitempty" db:"first_name"`
	LastName       string         `json:"lastName,omitempty" db:"last_name"`
	Password  string `json:"-" pg:"-" sql:"-"`
	Email          string         `json:"email,omitempty" db:"email"`
	Username  sql.NullString `json:"-"`
	Mobile    sql.NullString  `json:"-" db:"mobile"`
	Phone     sql.NullString  `json:"-" db:"phone"`
	Address   sql.NullString  `json:"-" db:"address"`
	ExternalID     string    `json:"-" db:"external_id"`
	UUID           uuid.UUID      `json:"id,omitempty" db:"uuid"`
	OrganizationID int       `json:"-" db:"organization_id"`
	Profile        []Profile      `json:"profiles,omitempty"`
	TimeZone       *string        `json:"timeZone" db:"timezone"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID        int
	CompanyID int
	Username  string
	Email     string
	Role      AccessRole
}

//func (u *User) FirstName(firstName string) {
//	u.FirstName = firstName
//}
//
//func (u *User) SetLastName(lastName string) {
//	u.LastName = lastName
//}
//
//func (u *User) SetTimeZone(tz string) {
//	u.TimeZone = &tz
//}
