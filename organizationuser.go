package authapi

import "github.com/google/uuid"

// Company represents Organization model
type OrganizationUser struct {
	Base
	UUID           uuid.UUID `json:"uuid" pg:",unique,type:uuid,notnull"`
	UserID         int       `json:"user"`
	User           *User     `json:"user"`
	OrganizationID int
	Organization   *Organization `json:"organization"`
	RoleID         int           `json:"-"`
	Role           *Role         `json:"-"`
	Active         bool          `json:"active"`
}
