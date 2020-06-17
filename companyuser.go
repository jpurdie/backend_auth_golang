package authapi

import "github.com/google/uuid"

// Company represents company model
type CompanyUser struct {
	Base
	UUID      uuid.UUID `json:"uuid"`
	UserID    int       `json:"user"`
	User      *User     `json:"user"`
	CompanyID int
	Company   *Company `json:"company"`
	RoleID    int      `json:"-"`
	Role      *Role    `json:"-"`
}
