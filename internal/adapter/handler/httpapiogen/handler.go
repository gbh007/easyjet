package httpapiogen

import (
	"context"
	"net/http"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapiogen/ogenapi"
	"github.com/gbh007/easyjet/internal/core/port"
)

type Handler struct {
	service port.Service
}

func NewHandler(service port.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateProject(ctx context.Context, req *ogenapi.ProjectCreate) (*ogenapi.CreateProjectCreated, error) {
	project := convertProjectCreate(req)

	id, err := h.service.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return &ogenapi.CreateProjectCreated{
		ID: ogenapi.NewOptUint(id),
	}, nil
}

func (h *Handler) CreateProjectRun(ctx context.Context, params ogenapi.CreateProjectRunParams) (*ogenapi.CreateProjectRunCreated, error) {
	runID, err := h.service.RunProject(ctx, params.ProjectID)
	if err != nil {
		return nil, err
	}

	return &ogenapi.CreateProjectRunCreated{
		ID: ogenapi.NewOptUint(runID),
	}, nil
}

func (h *Handler) GetProject(ctx context.Context, params ogenapi.GetProjectParams) (*ogenapi.Project, error) {
	project, err := h.service.Project(ctx, params.ProjectID)
	if err != nil {
		return nil, err
	}

	res := convertProjectToOgen(project)

	return &res, nil
}

func (h *Handler) GetProjectRun(ctx context.Context, params ogenapi.GetProjectRunParams) (*ogenapi.ProjectRun, error) {
	run, err := h.service.ProjectRun(ctx, params.RunID)
	if err != nil {
		return nil, err
	}

	res := convertProjectRunToOgen(run)

	return &res, nil
}

func (h *Handler) GetProjectRuns(ctx context.Context, params ogenapi.GetProjectRunsParams) (*ogenapi.GetProjectRunsOK, error) {
	runs, err := h.service.ProjectRuns(ctx, params.ProjectID)
	if err != nil {
		return nil, err
	}

	return &ogenapi.GetProjectRunsOK{
		Runs: convertProjectRunsToOgen(runs),
	}, nil
}

func (h *Handler) GetProjects(ctx context.Context) (*ogenapi.GetProjectsOK, error) {
	projects, err := h.service.Projects(ctx)
	if err != nil {
		return nil, err
	}

	return &ogenapi.GetProjectsOK{
		Projects: convertProjectsToOgen(projects),
	}, nil
}

func (h *Handler) UpdateProject(ctx context.Context, req *ogenapi.ProjectUpdate, params ogenapi.UpdateProjectParams) error {
	project := convertProjectUpdate(req, params.ProjectID)

	err := h.service.UpdateProject(ctx, project)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) NewError(ctx context.Context, err error) *ogenapi.ErrorStatusCode {
	return &ogenapi.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: ogenapi.Error{
			Error: ogenapi.NewOptString(err.Error()),
		},
	}
}
