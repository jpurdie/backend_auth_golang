package model

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"time"
)

type Project struct {
	authapi.Base
	UUID           uuid.UUID             `json:"id" pg:",unique,type:uuid,notnull"`
	Name           string                `json:"name" validate:"required"`
	Description    string                `json:"descr" validate:"omitempty,max=1000"`
	RGT            string                `json:"rgt"  validate:"required,oneof=R G T"`
	StatusID       int                   `json:"-" pg:",notnull"`
	Status         *ProjectStatus        `json:"status" validate:"required"`
	TypeID         int                   `json:"-" pg:",notnull"`
	Type           *ProjectType          `json:"type" validate:"required"`
	ComplexityID   int                   `json:"-" pg:",notnull"`
	Complexity     *ProjectComplexity    `json:"complexity" validate:"required"`
	SizeID         int                   `json:"-" pg:",notnull"`
	Size           *ProjectSize          `json:"size" validate:"required"`
	OrganizationID int                   `json:"-"`
	Organization   *authapi.Organization `json:"-"`
	EstStartDate   *time.Time            `pg:"type:date"`
	EstEndDate     *time.Time            `pg:"type:date"`
	ActStartDate   *time.Time            `pg:"type:date"`
	ActEndDate     *time.Time            `pg:"type:date"`
}

type ProjectStatus struct {
	authapi.Base
	UUID  uuid.UUID `json:"id" pg:",unique,type:uuid,notnull"`
	Name  string    `json:"name"`
	Order int       `json:"order"`
}

type ProjectType struct {
	authapi.Base
	UUID  uuid.UUID `json:"id" pg:",unique,type:uuid,notnull"`
	Name  string    `json:"name"`
	Order int       `json:"order"`
}

type ProjectComplexity struct {
	authapi.Base
	UUID   uuid.UUID `json:"id" pg:",unique,type:uuid,notnull"`
	Name   string    `json:"name"`
	Weight int       `json:"weight"`
}

type ProjectSize struct {
	authapi.Base
	UUID   uuid.UUID `json:"id" pg:",unique,type:uuid,notnull"`
	Name   string    `json:"name"`
	Weight int       `json:"weight"`
}
