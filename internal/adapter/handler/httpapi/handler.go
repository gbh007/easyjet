package httpapi

import (
	"context"
	"net/http"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
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

func (h *Handler) NewError(ctx context.Context, err error) *ogenapi.ErrorStatusCode {
	return &ogenapi.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: ogenapi.Error{
			Error: ogenapi.NewOptString(err.Error()),
		},
	}
}
