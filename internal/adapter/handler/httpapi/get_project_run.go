package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// GetProjectRun handles GET /run/{runID} endpoint.
func (h *Handler) GetProjectRun(ctx context.Context, params ogenapi.GetProjectRunParams) (*ogenapi.ProjectRun, error) {
	run, err := h.service.ProjectRun(ctx, params.RunID)
	if err != nil {
		return nil, err
	}

	res := convertProjectRunToOgen(run)

	return &res, nil
}
