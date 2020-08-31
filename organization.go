package authapi

import (
	"github.com/google/uuid"
)

type Organization struct {
	Base
	Name    string    `json:"name"`
	Active  bool      `json:"active"`
	Profile []Profile `json:"-" pg:",many2many:profiles"`
	UUID    uuid.UUID `json:"organizationID" pg:",unique,type:uuid,notnull"`
}
