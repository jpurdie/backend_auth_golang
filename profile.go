package authapi

import "github.com/google/uuid"

// Company represents Organization model
type Profile struct {
	Base
	UUID           uuid.UUID     `json:"profileID" pg:",unique,type:uuid,notnull"`
	UserID         int           `json:"-" pg:",notnull"`
	User           *User         `json:"-"`
	OrganizationID int           `json:"-" pg:",notnull"`
	Organization   *Organization `json:"organization"`
	RoleID         int           `json:"-" pg:",notnull"`
	Role           *Role         `json:"role"`
	Active         bool          `json:"-" pg:",notnull"`
}
