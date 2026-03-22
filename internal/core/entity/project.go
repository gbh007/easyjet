package entity

import (
	"time"
)

type Project struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	CronEnabled  bool
	CronSchedule string

	Dir            string
	GitURL         string
	GitBranch      string
	Name           string
	RestartAfter   bool
	RetentionCount int
	WithRootEnv    bool

	Stages []ProjectStage
}

func (p Project) HasGIT() bool {
	return p.GitURL != "" && p.GitBranch != ""
}

type ProjectStage struct {
	ProjectID uint
	Number    int
	Script    string
}
