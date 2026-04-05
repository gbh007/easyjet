package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// CreateProjectRun handles POST /project/{projectID}/run endpoint.
func (h *Handler) CreateProjectRun(ctx context.Context, params ogenapi.CreateProjectRunParams) (*ogenapi.CreateProjectRunCreated, error) {
	runID, err := h.service.RunProject(ctx, params.ProjectID)
	if err != nil {
		return nil, err
	}

	return &ogenapi.CreateProjectRunCreated{
		ID: ogenapi.NewOptInt(runID),
	}, nil
}
