package metrics

import (
	"context"
	"log/slog"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/gbh007/easyjet/internal/core/port"
)

type Controller struct {
	logger *slog.Logger
	ps     port.PubSub
}

func New(
	logger *slog.Logger,
	ps port.PubSub,
) Controller {
	return Controller{
		logger: logger,
		ps:     ps,
	}
}

func (cnt Controller) Start(ctx context.Context) {
	mChan := cnt.ps.SubscribeEvent("metrics", 100)

	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-mChan:
			if !ok {
				return
			}
			cnt.handle(e)
		}
	}
}

func (cnt Controller) handle(e entity.Event) {
	switch e.Type {
	case entity.EventRunFinished:
		observeRun(e.ProjectID, e.Err == nil, e.Duration)

	case entity.EventProjectCreated,
		entity.EventProjectUpdated,
		entity.EventProjectDeleted,
		entity.EventRunGitFinished,
		entity.EventRunStageFinished,
		entity.EventRunRotateFinished,
		entity.EventRequireAppRestart:
	}
}
