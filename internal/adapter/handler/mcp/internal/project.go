package internal

import (
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type ProjectResponse struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Dir            string           `json:"dir,omitempty"`
	GitURL         string           `json:"git_url,omitempty"`
	GitBranch      string           `json:"git_branch,omitempty"`
	CronEnabled    bool             `json:"cron_enabled,omitempty"`
	CronSchedule   string           `json:"cron_schedule,omitempty"`
	RestartAfter   bool             `json:"restart_after,omitempty"`
	RetentionCount int              `json:"retention_count,omitempty"`
	WithRootEnv    bool             `json:"with_root_env,omitempty"`
	IsTemplate     bool             `json:"is_template,omitempty"`
	Stages         []StageResponse  `json:"stages,omitempty"`
	EnvVars        []EnvVarResponse `json:"env_vars,omitempty"`
	CreatedAt      time.Time        `json:"created_at,omitzero"`
	UpdatedAt      time.Time        `json:"updated_at,omitzero"`
}

func ToProjectResponse(p entity.Project) ProjectResponse {
	return ProjectResponse{
		ID:             p.ID,
		Name:           p.Name,
		Dir:            p.Dir,
		GitURL:         p.GitURL,
		GitBranch:      p.GitBranch,
		CronEnabled:    p.CronEnabled,
		CronSchedule:   p.CronSchedule,
		RestartAfter:   p.RestartAfter,
		RetentionCount: p.RetentionCount,
		WithRootEnv:    p.WithRootEnv,
		IsTemplate:     p.IsTemplate,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		Stages:         ToStagesResponse(p.Stages),
		EnvVars:        ToEnvVarsResponse(p.EnvVars),
	}
}

type StageResponse struct {
	Number int    `json:"number"`
	Script string `json:"script"`
}

func ToStagesResponse(stages []entity.ProjectStage) []StageResponse {
	result := make([]StageResponse, 0, len(stages))
	for _, s := range stages {
		result = append(result, StageResponse{
			Number: s.Number,
			Script: s.Script,
		})
	}
	return result
}
