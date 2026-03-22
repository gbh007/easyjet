package entity

import (
	"time"
)

type Project struct {
	ID        uint      `param:"project_id" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CronEnabled  bool   `json:"cron_enabled"`
	CronSchedule string `json:"cron_schedule" validate:"omitzero,cron"`

	Dir            string `json:"dir"`
	GitURL         string `json:"git_url"`
	GitBranch      string `json:"git_branch"`
	Name           string `json:"name"`
	RestartAfter   bool   `json:"restart_after"`
	RetentionCount int    `json:"retention_count"`
	WithRootEnv    bool   `json:"with_root_env"`

	Stages []ProjectStage `json:"stages" validate:"min=1,dive"`
}

func (Project) TableName() string {
	return "projects"
}

func (p Project) HasGIT() bool {
	return p.GitURL != "" && p.GitBranch != ""
}

type ProjectStage struct {
	ProjectID uint   `json:"project_id"`
	Number    int    `json:"number" validate:"min=1"`
	Script    string `json:"script" validate:"required,gt=1"`
}

func (ProjectStage) TableName() string {
	return "stages"
}
