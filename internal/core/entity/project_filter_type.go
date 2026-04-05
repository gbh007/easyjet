package entity

type ProjectFilterType string

const (
	ProjectFilterTypeAll      ProjectFilterType = "all"
	ProjectFilterTypeProject  ProjectFilterType = "project"
	ProjectFilterTypeTemplate ProjectFilterType = "template"
)
