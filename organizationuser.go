package authapi

import "github.com/google/uuid"

// Company represents Organization model
type OrganizationUser struct {
	Base
	UUID           uuid.UUID `json:"uuid"`
	UserID         int       `json:"user"`
	User           *User     `json:"user"`
	OrganizationID int
	Organization   *Organization `json:"organization"`
	RoleID         int           `json:"-"`
	Role           *Role         `json:"-"`
}
