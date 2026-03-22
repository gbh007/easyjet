package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// UpdateProject handles PUT /project/{projectID} endpoint.
func (h *Handler) UpdateProject(ctx context.Context, req *ogenapi.ProjectUpdate, params ogenapi.UpdateProjectParams) error {
	project := convertProjectUpdate(req, params.ProjectID)

	err := h.service.UpdateProject(ctx, project)
	if err != nil {
		return err
	}

	return nil
}
