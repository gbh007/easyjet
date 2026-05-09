package hook

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"text/template"

	"github.com/gbh007/easyjet/config"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/gbh007/easyjet/internal/core/port"
)

type Controller struct {
	logger *slog.Logger
	ps     port.PubSub
	hooks  []config.Hook
	db     port.Database
}

func New(
	logger *slog.Logger,
	ps port.PubSub,
	hooks []config.Hook,
	db port.Database,
) Controller {
	return Controller{
		logger: logger,
		ps:     ps,
		hooks:  hooks,
		db:     db,
	}
}

func (cnt Controller) Start(ctx context.Context) {
	if len(cnt.hooks) == 0 {
		return
	}

	mChan := cnt.ps.SubscribeEvent("hooks", 100)

	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-mChan:
			if !ok {
				return
			}
			cnt.handle(ctx, e)
		}
	}
}

func (cnt Controller) handle(ctx context.Context, e entity.Event) {
	switch e.Type {
	case entity.EventRunFinished, entity.EventRunStageFinished:
		cnt.sendHook(ctx, e)

	case entity.EventProjectCreated,
		entity.EventProjectUpdated,
		entity.EventProjectDeleted,
		entity.EventRunGitFinished,
		entity.EventRunRotateFinished,
		entity.EventRequireAppRestart:
	}
}

func (cnt Controller) sendHook(ctx context.Context, e entity.Event) {
	project, err := cnt.db.Project(ctx, e.ProjectID)
	if err != nil {
		cnt.logger.Error("send hook: get project", "error", err)
		return
	}

	for _, hook := range cnt.hooks {
		if !filter(hook, e) {
			continue
		}

		u, err := url.Parse(hook.URL)
		if err != nil {
			cnt.logger.Error("send hook: parse url", "error", err)
			return
		}

		tmlp, err := template.New("").Parse(hook.Body)
		if err != nil {
			cnt.logger.Error("send hook: parse body template", "error", err)
			return
		}

		buff := &bytes.Buffer{}

		err = tmlp.Execute(buff, map[string]any{
			"err":      e.Err,
			"duration": e.Duration,
			"project": map[string]any{
				"name": project.Name,
			},
		})
		if err != nil {
			cnt.logger.Error("send hook: render body", "error", err)
			return
		}

		resp, err := http.DefaultClient.Do(&http.Request{
			Method: hook.Method,
			URL:    u,
			Header: hook.Headers,
			Body:   io.NopCloser(buff),
		})
		if err != nil {
			cnt.logger.Error("send hook", "error", err)
			return
		}

		if resp.Body != nil {
			_ = resp.Body.Close()
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			cnt.logger.Error("send hook: unsuccess", "code", resp.StatusCode)
		}
	}
}

func filter(hook config.Hook, e entity.Event) bool {
	if len(hook.Filter.ProjectIDs) > 0 {
		if !slices.Contains(hook.Filter.ProjectIDs, e.ProjectID) {
			return false
		}
	}

	if len(hook.Filter.Status) > 0 {
		allow := false

	loop:
		for _, status := range hook.Filter.Status {
			switch status {
			case "all":
				allow = true
				break loop
			case "job_all":
				if e.Type == entity.EventRunFinished {
					allow = true
					break loop
				}
			case "job_success":
				if e.Type == entity.EventRunFinished &&
					e.Err == nil {
					allow = true
					break loop
				}
			case "job_fail":
				if e.Type == entity.EventRunFinished &&
					e.Err != nil {
					allow = true
					break loop
				}
			case "stage_all":
				if e.Type == entity.EventRunStageFinished {
					allow = true
					break loop
				}
			case "stage_success":
				if e.Type == entity.EventRunStageFinished &&
					e.Err == nil {
					allow = true
					break loop
				}
			case "stage_fail":
				if e.Type == entity.EventRunStageFinished &&
					e.Err != nil {
					allow = true
					break loop
				}
			}
		}

		if !allow {
			return false
		}
	}

	return true
}
