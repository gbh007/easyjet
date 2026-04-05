package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

// CreateProject handles POST /project endpoint.
func (h *Handler) CreateProject(ctx context.Context, req *ogenapi.ProjectCreate) (*ogenapi.CreateProjectCreated, error) {
	project := convertProjectCreate(req)

	id, err := h.service.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return &ogenapi.CreateProjectCreated{
		ID: ogenapi.NewOptInt(id),
	}, nil
}
