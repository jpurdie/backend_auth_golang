package authapi

import (
	"github.com/google/uuid"
	"time"
)

type Invitation struct {
	Base
	TokenHash      string        `json:"-" pg:",alias:token, unique:invitation_token"` //represents the hash of the token
	TokenStr       string        `json:"-" pg:"-" sql:"-"`                             // represents the plaintext token string
	ExpiresAt      *time.Time    `json:"expires_at"`                                   //token expiration
	InvitorID      int           `json:"-"`                                            //ID of the person sending the invite
	Invitor        *User         `json:"-"`                                            //person sending the invitation
	OrganizationID int           `json:"-"`
	Organization   *Organization `json:"organization"`
	Email          string        `json:"email"` //email of the user being invited
	Used           bool          `json:"used" pg:"default:FALSE"`
	UUID           uuid.UUID     `json:"-" pg:",unique,type:uuid"`
}
