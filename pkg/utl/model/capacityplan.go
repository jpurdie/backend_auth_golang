package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"gopkg.in/guregu/null.v4"
)

type CapacityPlan struct {
	CapacityPlans  []CapacityPlanEntry `json:"capacityPlanEntries,array,omitempty" validate:"required"`
	ResourceID     int                 `json:"-" validate:"required"`
	ResourceUUID   *uuid.UUID          `json:"resourceID,omitempty"`
	WorkDate       time.Time           `json:"workDate" validate:"required"`
	SumWorkPercent int                 `json:"sumWorkPercent"`
}

type CapacityPlanEntry struct {
	authapi.Base
	UUID         uuid.UUID     `json:"id" db:"uuid"`
	ProjectID    int           `json:"-"  db:"project_id"` // pk tag is used to mark field as primary key
	ProjectUUID  uuid.UUID     `db:"projectUUID" json:"projectID,omitempty"`
	Project      *Project      `db:"project" json:"project,omitempty"`
	ResourceID   int           `json:"-" db:"resource_id"`
	ResourceUUID uuid.UUID     `db:"resourceUUID" json:"resourceID,omitempty"`
	Resource     *authapi.User `json:"resource,omitempty" db:"resource"`
	WorkDate     null.Time     `json:"workDate" db:"work_date"`
	WorkPercent  int           `json:"workPercent" db:"work_percent"`
	CreatedAt    null.Time     `json:"-" db:"created_at"`
	UpdatedAt    null.Time     `json:"-" db:"updated_at"`
	DeletedAt    null.Time     `json:"-" db:"deleted_at"`
}
