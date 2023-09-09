package model

import (
	"time"

	"github.com/emzola/venato/project/pkg/validator"
)

// Project defines the project data.
type Project struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	TargetEndDate time.Time `json:"target_end_date"`
	ActualEndDate time.Time `json:"actual_end_date,omitempty"`
	CreatedOn     time.Time `json:"created_on"`
	CreatedBy     int64     `json:"created_by"`
	ModifiedOn    time.Time `json:"modified_on,omitempty"`
	ModifiedBy    int64     `json:"modified_by,omitempty"`
	Version       int64     `json:"version"`
}

// ValidateProject performs data validation on supplied project data.
func ValidateProject(v *validator.Validator, project *Project) {
	v.Check(project.Name != "", "name", "must be provided")
	v.Check(len(project.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(len(project.Description) <= 1000, "description", "must not be more than 1000 bytes long")
}
