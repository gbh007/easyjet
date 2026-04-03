package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// GetProjectRuns handles GET /project/{projectID}/runs endpoint.
func (h *Handler) GetProjectRuns(ctx context.Context, params ogenapi.GetProjectRunsParams) (*ogenapi.GetProjectRunsOK, error) {
	runs, err := h.service.ProjectRuns(ctx, params.ProjectID)
	if err != nil {
		return nil, err
	}

	return &ogenapi.GetProjectRunsOK{
		Runs: ogenapi.NewOptNilProjectRunSummaryArray(convertProjectRunsSummaryToOgen(runs)),
	}, nil
}
