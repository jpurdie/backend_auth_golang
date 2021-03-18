package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
)

type Project struct {
	authapi.Base
	UUID                uuid.UUID             `json:"id,omitempty" db:"uuid"`
	Name                string                `json:"name,omitempty" validate:"required" db:"name"`
	Description         string                `json:"descr,omitempty" validate:"omitempty,max=1000"`
	RGT                 string                `json:"rgt,omitempty"  validate:"required,oneof=R G T" db:"rgt"`
	StatusID            int                   `json:"-" db:"status_id"`
	Status              *ProjectStatus        `json:"status,omitempty" validate:"required" db:"project_status"`
	TypeID              int                   `json:"-"  db:"type_id"`
	Type                *ProjectType          `json:"type,omitempty" validate:"required" db:"project_type"`
	ComplexityID        int                   `json:"-" db:"complexity_id"`
	Complexity          *ProjectComplexity    `json:"complexity,omitempty" validate:"required" db:"project_complexity"`
	SizeID              int                   `json:"-" db:"size_id"`
	Size                *ProjectSize          `json:"size,omitempty" validate:"required" db:"project_size"`
	OrganizationID      int                   `json:"-" db:"organization_id"`
	Organization        *authapi.Organization `json:"-"`
	EstStartDate        *time.Time            `db:"est_start_date" json:"omitempty"`
	EstEndDate          *time.Time            `db:"est_end_date" json:"omitempty"`
	ActStartDate        *time.Time            `db:"act_start_date" json:"omitempty"`
	ActEndDate          *time.Time            `db:"act_end_date" json:"omitempty"`
	StrategicAlignments []*StrategicAlignment `json:"strategicAlignments,omitempty" pg:"many2many:project_alignments"`
	OpenForTimeEntry    *bool                 `db:"open_for_time_entry" json:"openForTimeEntry,omitempty"`
	Compliance          *bool                 `db:"compliance" json:"compliance,omitempty"`
	TimeConstrained     *bool                 `db:"time_constrained" json:"timeConstrained,omitempty"`
	SponsorAreas        []*SponsorArea        `json:"sponsorAreas,omitempty" pg:"many2many:project_sponsor_areas"`
}

type ProjectStatus struct {
	authapi.Base
	UUID  uuid.UUID `json:"id"  db:"uuid"`
	Name  string    `json:"name" db:"name"`
	Order int       `json:"order" db:"order"`
}

type ProjectType struct {
	authapi.Base
	UUID  uuid.UUID `json:"id" db:"uuid"`
	Name  string    `json:"name" db:"name"`
	Order int       `json:"order" db:"order"`
}

type ProjectComplexity struct {
	authapi.Base
	UUID   uuid.UUID `json:"id" db:"uuid"`
	Name   string    `json:"name" db:"name"`
	Weight int       `json:"weight" db:"weight"`
}

type ProjectSize struct {
	authapi.Base
	UUID   uuid.UUID `json:"id"  db:"uuid"`
	Name   string    `json:"name" db:"name"`
	Weight int       `json:"weight" db:"weight"`
}

type ProjectStrategicAlignment struct {
	ProjectID            int                 `db:"project_id""` // pk tag is used to mark field as primary key
	Project              *Project            ``
	StrategicAlignmentID int                 ` db:"strategic_alignment_id"`
	StrategicAlignment   *StrategicAlignment `pg:"rel:has-one"`
}

type ProjectSponsorArea struct {
	ProjectID     int          `db:"project_id""` // pk tag is used to mark field as primary key
	Project       *Project     ``
	SponsorAreaID int          `db:"sponsor_area_id"`
	SponsorArea   *SponsorArea `pg:"rel:has-one"`
}
