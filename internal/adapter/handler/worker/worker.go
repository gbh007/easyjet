package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/gbh007/easyjet/internal/core/port"
)

type Controller struct {
	logger  *slog.Logger
	service port.Service
}

func New(
	logger *slog.Logger,
	service port.Service,
) Controller {
	return Controller{
		logger:  logger,
		service: service,
	}
}

func (cnt Controller) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cnt.run(ctx)
		}
	}
}

func (cnt Controller) run(ctx context.Context) {
	// TODO: вынести настройку таймаута в конфиг
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Minute)
	defer cancel()

	ids, err := cnt.service.PendingProjectRuns(ctx)
	if err != nil {
		cnt.logger.Error("get pending runs", "error", err.Error())
		return
	}

	for _, id := range ids {
		err = cnt.service.HandleRun(ctx, id)
		if err != nil {
			cnt.logger.Error("handle project run", "error", err.Error(), "run_id", id)
		}
	}
}
