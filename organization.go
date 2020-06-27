package authapi

import "github.com/google/uuid"

type Organization struct {
	Base
	Name             string             `json:"name"`
	Active           bool               `json:"active"`
	OrganizationUser []OrganizationUser `json:"-" pg:",many2many:organization_users"`
	UUID             uuid.UUID          `json:"uuid" pg:",unique"`
}
