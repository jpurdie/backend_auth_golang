package model

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
)

type StrategicAlignment struct {
	authapi.Base
	UUID           uuid.UUID             `json:"id" db:"uuid"`
	Name           string                `json:"name" db:"name"`
	DisplayOrder   sql.NullInt32         `json:"-" db:"display_order"`
	Projects       []*Project            `json:"-" pg:"many2many:project_alignments"`
	OrganizationID int                   `json:"-" db:"organization_id"`
	Organization   *authapi.Organization `json:"-"`
}
