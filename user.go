package authapi

import "github.com/google/uuid"

type User struct {
	Base
	FirstName        string             `json:"first_name"`
	LastName         string             `json:"last_name"`
	Password         string             `json:"-" pg:"-" sql:"-"`
	Email            string             `json:"email"`
	Username         string             `json:"username"`
	Mobile           string             `json:"mobile,omitempty"`
	Phone            string             `json:"phone,omitempty"`
	Address          string             `json:"address,omitempty"`
	Active           bool               `json:"active"`
	ExternalID       string             `json:"auth0id" pg:",unique"`
	UUID             uuid.UUID          `json:"uuid" pg:",unique,type:uuid"`
	OrganizationID   int                `json:"organization_id" pg:"-" sql:"-"`
	OrganizationUser []OrganizationUser `json:"-"  pg:",many2many:organization_users"`
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
