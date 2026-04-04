package mcp

import (
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type ProjectResponse struct {
	ID             uint             `json:"id"`
	Name           string           `json:"name"`
	Dir            string           `json:"dir"`
	GitURL         string           `json:"git_url"`
	GitBranch      string           `json:"git_branch"`
	CronEnabled    bool             `json:"cron_enabled"`
	CronSchedule   string           `json:"cron_schedule"`
	RestartAfter   bool             `json:"restart_after"`
	RetentionCount int              `json:"retention_count"`
	WithRootEnv    bool             `json:"with_root_env"`
	IsTemplate     bool             `json:"is_template"`
	Stages         []StageResponse  `json:"stages"`
	EnvVars        []EnvVarResponse `json:"env_vars"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

func toProjectResponse(p entity.Project) ProjectResponse {
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
		Stages:         toStagesResponse(p.Stages),
		EnvVars:        toEnvVarsResponse(p.EnvVars),
	}
}

type StageResponse struct {
	Number int    `json:"number"`
	Script string `json:"script"`
}

func toStagesResponse(stages []entity.ProjectStage) []StageResponse {
	result := make([]StageResponse, 0, len(stages))
	for _, s := range stages {
		result = append(result, StageResponse{
			Number: s.Number,
			Script: s.Script,
		})
	}
	return result
}

type EnvVarResponse struct {
	ID                 uint      `json:"id"`
	Name               string    `json:"name"`
	Value              string    `json:"value"`
	UsesOtherVariables bool      `json:"uses_other_variables"`
	ProjectID          *uint     `json:"project_id,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func toEnvVarResponse(ev entity.EnvironmentVariable) EnvVarResponse {
	return EnvVarResponse{
		ID:                 ev.ID,
		Name:               ev.Name,
		Value:              ev.Value,
		UsesOtherVariables: ev.UsesOtherVariables,
		ProjectID:          ev.ProjectID,
		CreatedAt:          ev.CreatedAt,
		UpdatedAt:          ev.UpdatedAt,
	}
}

func toEnvVarsResponse(envVars []entity.EnvironmentVariable) []EnvVarResponse {
	result := make([]EnvVarResponse, 0, len(envVars))
	for _, ev := range envVars {
		result = append(result, toEnvVarResponse(ev))
	}
	return result
}

type ProjectRunResponse struct {
	ID         uint                   `json:"id"`
	ProjectID  uint                   `json:"project_id"`
	Project    ProjectResponse        `json:"project"`
	Status     string                 `json:"status"`
	FailLog    string                 `json:"fail_log,omitempty"`
	Duration   string                 `json:"duration"`
	Stages     []RunStageResponse     `json:"stages"`
	GitCommits []RunGitCommitResponse `json:"git_commits"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

func toProjectRunResponse(run entity.ProjectRun) ProjectRunResponse {
	return ProjectRunResponse{
		ID:         run.ID,
		ProjectID:  run.ProjectID,
		Project:    toProjectResponse(run.Project),
		Status:     run.Status,
		FailLog:    run.FailLog,
		Duration:   run.Duration.String(),
		Stages:     toRunStagesResponse(run.Stages),
		GitCommits: toRunGitCommitsResponse(run.GitCommits),
		CreatedAt:  run.CreatedAt,
		UpdatedAt:  run.UpdatedAt,
	}
}

type RunStageResponse struct {
	StageNumber int    `json:"stage_number"`
	Success     bool   `json:"success"`
	Log         string `json:"log,omitempty"`
	Duration    string `json:"duration"`
}

func toRunStagesResponse(stages []entity.ProjectRunStage) []RunStageResponse {
	result := make([]RunStageResponse, 0, len(stages))
	for _, s := range stages {
		result = append(result, RunStageResponse{
			StageNumber: s.StageNumber,
			Success:     s.Success,
			Log:         s.Log,
			Duration:    s.Duration.String(),
		})
	}
	return result
}

type RunGitCommitResponse struct {
	Number  int    `json:"number"`
	Hash    string `json:"hash"`
	Subject string `json:"subject"`
}

func toRunGitCommitsResponse(commits []entity.ProjectRunGitCommits) []RunGitCommitResponse {
	result := make([]RunGitCommitResponse, 0, len(commits))
	for _, c := range commits {
		result = append(result, RunGitCommitResponse{
			Number:  c.Number,
			Hash:    c.Hash,
			Subject: c.Subject,
		})
	}
	return result
}

type RunResponse struct {
	RunID uint `json:"run_id"`
}

type SuccessMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type IDResponse struct {
	ID uint `json:"id"`
}
