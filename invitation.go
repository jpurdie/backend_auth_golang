package authapi

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	Base
	TokenHash      string        `json:"-" db:"token_hash"'`         //represents the hash of the token
	TokenStr       string        `json:"-" pg:"-" sql:"-"`           // represents the plaintext token string
	ExpiresAt      *time.Time    `json:"expires_at" db:"expires_at"` //token expiration
	InvitorID      int           `json:"-" db:"invitor_id"`          //ID of the person sending the invite
	Invitor        *User         `json:"-"`                          //person sending the invitation
	OrganizationID int           `json:"-" db:"organization_id"`
	Organization   *Organization `json:"organization"`
	Email          string        `json:"email" db:"email"` //email of the user being invited
	Used           bool          `json:"used" db:"used"`
	UUID           uuid.UUID     `json:"-" db:"uuid"`
}
