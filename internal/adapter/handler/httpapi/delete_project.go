package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

func (h *Handler) DeleteProject(ctx context.Context, params ogenapi.DeleteProjectParams) error {
	return h.service.DeleteProject(ctx, params.ProjectID)
}
