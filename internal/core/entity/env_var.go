package entity

import (
	"time"
)

type EnvironmentVariable struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID          *int
	Name               string
	Value              string
	UsesOtherVariables bool
}

func (e EnvironmentVariable) IsGlobal() bool {
	return e.ProjectID == nil
}
