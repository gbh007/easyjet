package entity

import (
	"time"
)

type Project struct {
	ID        uint      `param:"project_id" json:"id" gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`

	CronEnabled  bool   `json:"cron_enabled" gorm:"column:cron_enabled;not null;default:false"`
	CronSchedule string `json:"cron_schedule" gorm:"column:cron_schedule;type:text;not null;default:''" validate:"omitzero,cron"`

	Dir            string `json:"dir" gorm:"column:dir;not null"`
	GitURL         string `json:"git_url" gorm:"column:git_url;not null"`
	GitBranch      string `json:"git_branch" gorm:"column:git_branch;not null"`
	Name           string `json:"name" gorm:"column:name;not null"`
	RetentionCount int    `json:"retention_count" gorm:"column:retention_count;not null;default:0"`

	Stages []ProjectStage `json:"stages" gorm:"foreignKey:ProjectID" validate:"min=1"`
}

func (Project) TableName() string {
	return "projects"
}

func (p Project) HasGIT() bool {
	return p.GitURL != "" && p.GitBranch != ""
}

type ProjectStage struct {
	ProjectID uint   `json:"project_id" gorm:"column:project_id;not null;index:idx_stages_project_id"`
	Number    int    `json:"number" gorm:"column:num;not null" validate:"min=1"`
	Script    string `json:"script" gorm:"column:script;not null" validate:"required"`
}

func (ProjectStage) TableName() string {
	return "stages"
}
