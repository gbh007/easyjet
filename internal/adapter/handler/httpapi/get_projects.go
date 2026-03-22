package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// GetProjects handles GET /projects endpoint.
func (h *Handler) GetProjects(ctx context.Context) (*ogenapi.GetProjectsOK, error) {
	projects, err := h.service.Projects(ctx)
	if err != nil {
		return nil, err
	}

	return &ogenapi.GetProjectsOK{
		Projects: ogenapi.NewOptNilProjectArray(convertProjectsToOgen(projects)),
	}, nil
}
