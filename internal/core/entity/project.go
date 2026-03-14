package entity

import (
	"time"
)

type Project struct {
	ID        uint      `gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`

	Dir       string `gorm:"column:dir;not null"`
	GitURL    string `gorm:"column:git_url;not null"`
	GitBranch string `gorm:"column:git_branch;not null"`
	Name      string `gorm:"column:name;not null"`

	Stages []ProjectStage `gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string {
	return "projects"
}

type ProjectStage struct {
	ProjectID uint   `gorm:"column:project_id;not null"`
	Number    int    `gorm:"column:num;not null"`
	Script    string `gorm:"column:script;not null"`
}

func (ProjectStage) TableName() string {
	return "stages"
}

type ProjectRun struct {
	ID        uint      `gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`

	Project Project           `gorm:"foreignKey:ProjectID"`
	Stages  []ProjectStageRun `gorm:"foreignKey:RunID"`

	ProjectID uint `gorm:"column:project_id;not null"`
	Success   bool `gorm:"column:success;not null"`
}

func (ProjectRun) TableName() string {
	return "runs"
}

type ProjectStageRun struct {
	RunID       uint   `gorm:"column:run_id;not null"`
	StageNumber int    `gorm:"column:stage_num;not null"`
	Success     bool   `gorm:"column:success;not null"`
	Log         string `gorm:"column:log;not null"`
}

func (ProjectStageRun) TableName() string {
	return "stage_runs"
}
