package entities

type Project struct {
	ID     uint
	Dir    string
	GitURL string
	Name   string
	Stages []string
}

type ProjectRun struct {
	ID        uint
	ProjectID uint
	Success   bool
	Stages    []ProjectRunStage
}

type ProjectRunStage struct {
	StageNumber int
	Success     bool
	Log         string
}
