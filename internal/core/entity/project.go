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
	IsTemplate     bool

	Stages  []ProjectStage
	EnvVars []EnvironmentVariable
}

func (p Project) HasGIT() bool {
	return p.GitURL != "" && p.GitBranch != ""
}

type ProjectStage struct {
	ProjectID uint
	Number    int
	Script    string
}

type ProjectsWithRunInfo struct {
	ID                  uint
	Name                string
	CronEnabled         bool
	IsTemplate          bool
	LastSuccessfulRunAt *time.Time
	LastRun             *ProjectLastRun
}

type ProjectLastRun struct {
	CreatedAt time.Time
	Status    string
	Duration  time.Duration
}
