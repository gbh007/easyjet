package internal

import (
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type ProjectRunResponse struct {
	ID         int                    `json:"id"`
	ProjectID  int                    `json:"project_id"`
	Status     string                 `json:"status"`
	FailLog    string                 `json:"fail_log,omitempty"`
	Duration   string                 `json:"duration,omitempty"`
	Stages     []RunStageResponse     `json:"stages,omitempty"`
	GitCommits []RunGitCommitResponse `json:"git_commits,omitempty"`
	CreatedAt  time.Time              `json:"created_at,omitzero"`
	UpdatedAt  time.Time              `json:"updated_at,omitzero"`
}

func ToProjectRunResponse(run entity.ProjectRun) ProjectRunResponse {
	return ProjectRunResponse{
		ID:         run.ID,
		ProjectID:  run.ProjectID,
		Status:     run.Status,
		FailLog:    run.FailLog,
		Duration:   run.Duration.String(),
		Stages:     ToRunStagesResponse(run.Stages),
		GitCommits: ToRunGitCommitsResponse(run.GitCommits),
		CreatedAt:  run.CreatedAt,
		UpdatedAt:  run.UpdatedAt,
	}
}

type RunStageResponse struct {
	StageNumber int    `json:"stage_number"`
	Success     bool   `json:"success"`
	Log         string `json:"log,omitempty"`
	Duration    string `json:"duration,omitempty"`
}

func ToRunStagesResponse(stages []entity.ProjectRunStage) []RunStageResponse {
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

func ToRunGitCommitsResponse(commits []entity.ProjectRunGitCommits) []RunGitCommitResponse {
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
