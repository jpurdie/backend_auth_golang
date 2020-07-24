package authapi

import "github.com/google/uuid"

type User struct {
	Base
	FirstName        string             `json:"first_name"`
	LastName         string             `json:"last_name"`
	Password         string             `json:"-" pg:"-" sql:"-"`
	Email            string             `json:"email"`
	Username         string             `json:"-"`
	Mobile           string             `json:"mobile"`
	Phone            string             `json:"phone"`
	Address          string             `json:"address"`
	Active           bool               `json:"active"`
	ExternalID       string             `json:"-" pg:",unique"`
	UUID             uuid.UUID          `json:"id" pg:",unique,type:uuid,notnull"`
	OrganizationID   int                `json:"-" pg:"-" sql:"-"`
	OrganizationUser []OrganizationUser `json:"-"  pg:",many2many:organization_users,joinFK:id""`
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
