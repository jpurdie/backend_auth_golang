package authapi

import (
	"database/sql"
	"github.com/google/uuid"
)

type User struct {
	Base
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Password  string `json:"-" pg:"-" sql:"-"`
	Email     string `json:"email" db:"email"`
	Username  sql.NullString `json:"-"`
	Mobile    sql.NullString  `json:"-" db:"mobile"`
	Phone     sql.NullString  `json:"-" db:"phone"`
	Address   sql.NullString  `json:"-" db:"address"`
	//	Active           bool               `json:"active"`
	ExternalID     string    `json:"-" db:"external_id"`
	UUID           uuid.UUID `json:"userID" db:"uuid"`
	OrganizationID int       `json:"-" db:"organization_id"`
	Profile        []Profile `json:"profiles"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID        int
	CompanyID int
	Username  string
	Email     string
	Role      AccessRole
}

// ChangePassword updates user's password related fields
//func (u *User) ChangePassword(hash string) {
//	u.Password = hash
//	u.LastPasswordChange = time.Now()
//}
//
//// UpdateLastLogin updates last login field
//func (u *User) UpdateLastLogin(token string) {
//	u.Token = token
//	u.LastLogin = time.Now()
//}
