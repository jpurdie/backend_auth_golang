package authapi

import "github.com/google/uuid"

// Company represents Organization model
type OrganizationUser struct {
	Base
	UUID           uuid.UUID     `json:"uuid" pg:",unique,type:uuid,notnull"`
	UserID         int           `json:"user" pg:"notnull"`
	User           *User         `json:"user"`
	OrganizationID int           `pg:"notnull"`
	Organization   *Organization `json:"organization"`
	RoleID         int           `json:"-" pg:"notnull"`
	Role           *Role         `json:"-"`
	Active         bool          `json:"active" pg:"notnull"`
}
