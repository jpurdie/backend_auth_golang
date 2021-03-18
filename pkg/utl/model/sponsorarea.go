package model

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
)

type SponsorArea struct {
	authapi.Base
	UUID           uuid.UUID             `json:"id" db:"uuid"`
	Name           string                `json:"name" db:"name"`
	Projects       []*Project            `json:"-" pg:"many2many:project_sponsor_areas"`
	OrganizationID int                   `json:"-" db:"organization_id"`
	Organization   *authapi.Organization `json:"-"`
}
