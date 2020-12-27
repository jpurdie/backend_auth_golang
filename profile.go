package authapi

import (
	"github.com/google/uuid"
)

// Company represents Profile model
type Profile struct {
	Base
	UUID           uuid.UUID     `json:"profileID" db:"uuid"`
	UserID         int           `json:"-" db:"user_id"`
	User           *User         `json:"-"`
	OrganizationID int           `json:"-" db:"organization_id"`
	Organization   *Organization `json:"organization"`
	RoleID         int           `json:"-" db:"role_id"`
	Role           *Role         `json:"role"`
	Active         bool          `json:"-" db:"active"`
}
