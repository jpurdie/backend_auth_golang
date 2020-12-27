package authapi

import (
	"github.com/google/uuid"
)

type Organization struct {
	Base
	Name    string    `json:"name"  db:"name"`
	Active  bool      `json:"active"  db:"active"`
	Profile []Profile `json:"-" pg:",many2many:profiles"`
	UUID    uuid.UUID `json:"organizationID"  db:"uuid"`
}
