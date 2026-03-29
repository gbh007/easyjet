package entity

import (
	"time"
)

type EnvironmentVariable struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID          *uint
	Name               string
	Value              string
	UsesOtherVariables bool
}

func (e EnvironmentVariable) IsGlobal() bool {
	return e.ProjectID == nil
}
