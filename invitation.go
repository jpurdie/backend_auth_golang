package authapi

import "time"

type Invitation struct {
	Base
	Token          string        `json:"-" pg:"unique:invitation_token"` //represents the hash of the plaintext token string
	TokenStr       string        `json:"token" pg:"-" sql:"-"`           // represents the plaintext token string
	ExpiresAt      *time.Time    `json:"expires_at"`                     //token expiration
	InvitorID      int           `json:"-"`                              //ID of the person sending the invite
	Invitor        *User         `json:"-"`                              //person sending the invitation
	OrganizationID int           `json:"_"`
	Organization   *Organization `json:"-"`
	Email          string        `json:"email"` //email of the user being invited
	Used           bool          `json:used`
}
