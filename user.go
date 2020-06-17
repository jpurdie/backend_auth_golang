package authapi

import (
	"time"
)

// User represents user domain model
type User struct {
	Base
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Password   string `json:"-"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Mobile     string `json:"mobile,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Address    string `json:"address,omitempty"`
	Active     bool   `json:"active"`
	ExternalID string `json:"auth0id"`

	LastLogin          time.Time `json:"last_login,omitempty"`
	LastPasswordChange time.Time `json:"last_password_change,omitempty"`

	Token  string     `json:"-"`
	RoleID AccessRole `json:"-"`

	Role        *Role         `json:"role,omitempty",pg:"-"`
	CompanyID   int           `json:"company_id",pg:"-"`
	LocationID  int           `json:"location_id, pg:"-""`
	CompanyUser []CompanyUser `json:"-", pg:",many2many:company_users"`
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
func (u *User) ChangePassword(hash string) {
	u.Password = hash
	u.LastPasswordChange = time.Now()
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin(token string) {
	u.Token = token
	u.LastLogin = time.Now()
}
