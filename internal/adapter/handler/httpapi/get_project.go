package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// GetProject handles GET /project/{projectID} endpoint.
func (h *Handler) GetProject(ctx context.Context, params ogenapi.GetProjectParams) (*ogenapi.Project, error) {
	project, err := h.service.Project(ctx, params.ProjectID)
	if err != nil {
		return nil, err
	}

	res := convertProjectToOgen(project)

	return &res, nil
}
