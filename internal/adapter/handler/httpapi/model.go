package httpapi

type projectIDRequest struct {
	ProjectID uint `param:"project_id" validate:"min=1"`
}

type projectIDAndRunIDRequest struct {
	ProjectID uint `param:"project_id" validate:"min=1"`
	RunID     uint `param:"run_id" validate:"min=1"`
}
