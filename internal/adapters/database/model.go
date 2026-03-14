package database

import "gorm.io/gorm"

type modelProject struct {
	gorm.Model

	Dir    string `gorm:"column:dir"`
	GitURL string `gorm:"column:git_url"`
	Name   string `gorm:"column:name"`

	Stages []modelProjectStage `gorm:"foreignKey:ProjectID"`
}

func (modelProject) TableName() string {
	return "projects"
}

type modelProjectRun struct {
	gorm.Model

	Project modelProject `gorm:"foreignKey:ProjectID"`

	ProjectID   uint   `gorm:"column:project_id"`
	StageNumber int    `gorm:"column:stage_num"`
	Success     bool   `gorm:"column:success"`
	Log         string `gorm:"column:log"`
}

func (modelProjectRun) TableName() string {
	return "runs"
}

type modelProjectStage struct {
	ProjectID uint   `gorm:"column:project_id"`
	Number    int    `gorm:"column:num"`
	Script    string `gorm:"column:script"`
}

func (modelProjectStage) TableName() string {
	return "stages"
}
