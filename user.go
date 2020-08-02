package authapi

import "github.com/google/uuid"

type User struct {
	Base
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"-" pg:"-" sql:"-"`
	Email     string `json:"email"`
	Username  string `json:"-"`
	Mobile    string `json:"mobile"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	//	Active           bool               `json:"active"`
	ExternalID     string    `json:"-" pg:",unique"`
	UUID           uuid.UUID `json:"userID" pg:",unique,type:uuid,notnull"`
	OrganizationID int       `json:"-" pg:"-" sql:"-"`
	Profile        []Profile `json:"profiles"  pg:",many2many:profiles,joinFK:id""`
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
