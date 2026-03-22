package entity

import "time"

type ProjectRun struct {
	ID        uint      `param:"run_id" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Project    Project                `json:"project,omitzero"`
	Stages     []ProjectRunStage      `json:"stages,omitempty,omitzero" validate:"min=1"`
	GitCommits []ProjectRunGitCommits `json:"git_commits,omitempty,omitzero"`

	ProjectID  uint   `json:"project_id"`
	Success    bool   `json:"success"`
	Pending    bool   `json:"pending"`
	Processing bool   `json:"processing"`
	FailLog    string `json:"fail_log"`
}

func (ProjectRun) TableName() string {
	return "runs"
}

type ProjectRunStage struct {
	RunID       uint   `json:"run_id"`
	StageNumber int    `json:"stage_number" validate:"min=1"`
	Success     bool   `json:"success"`
	Log         string `json:"log"`
}

func (ProjectRunStage) TableName() string {
	return "run_stages"
}

type ProjectRunGitCommits struct {
	RunID   uint   `json:"run_id"`
	Number  int    `json:"number"`
	Hash    string `json:"hash"`
	Subject string `json:"subject"`
}

func (ProjectRunGitCommits) TableName() string {
	return "run_git_commits"
}
