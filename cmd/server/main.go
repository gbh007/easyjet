package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/gbh007/easyjet/config"
	"github.com/gbh007/easyjet/internal/adapter/exec/shellexec"
	"github.com/gbh007/easyjet/internal/adapter/filesystem/filesystem"
	"github.com/gbh007/easyjet/internal/adapter/git/shellgit"
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi"
	"github.com/gbh007/easyjet/internal/adapter/handler/mcp"
	"github.com/gbh007/easyjet/internal/adapter/handler/metrics"
	schedulerhandler "github.com/gbh007/easyjet/internal/adapter/handler/scheduler"
	"github.com/gbh007/easyjet/internal/adapter/handler/worker"
	"github.com/gbh007/easyjet/internal/adapter/pubsub/eventbus"
	"github.com/gbh007/easyjet/internal/adapter/repository/postgres"
	"github.com/gbh007/easyjet/internal/adapter/repository/sqlite"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/gbh007/easyjet/internal/core/port"
	"github.com/gbh007/easyjet/internal/core/service"
	"github.com/golang-cz/devslog"
	"github.com/lmittmann/tint"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	cfg, err := config.Read[config.Config]("config.toml")
	if err != nil {
		panic(err)
	}

	logger := slog.Default()
	llv := cfg.Log.SlogLevel()

	switch cfg.Log.Format {
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: llv,
		}))
	case "dev":
		logger = slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
			HandlerOptions: &slog.HandlerOptions{
				Level: llv,
			},
		}))
	case "tint":
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level: llv,
		}))
	}

	var db port.Database

	switch cfg.Database.Type {
	case "postgres":
		db, err = postgres.NewRepo(ctx, logger, cfg.Database.DNS)
	case "sqlite":
		db, err = sqlite.NewRepo(ctx, logger, cfg.Database.DNS)
	default:
		err = fmt.Errorf("unsupported db type: %s", cfg.Database.Type)
	}

	if err != nil {
		logger.Error("create database adapter", "error", err.Error())
		os.Exit(1) //nolint:gocritic // будет исправлено позднее
	}

	ex := shellexec.New(logger)
	fs := filesystem.New(logger, cfg.App.ProjectDir)
	git := shellgit.New(logger)
	ps := eventbus.New(logger)
	srv := service.New(logger, ex, fs, git, db, ps, cfg.Server.External.Web)

	apiCnt := httpapi.New(
		logger,
		httpapi.Config{
			Addr:            cfg.Server.Addr,
			User:            cfg.Server.User,
			Pass:            cfg.Server.Pass,
			StaticFilesPath: cfg.Server.StaticFilesPath,
			MCP: mcp.Config{
				Enabled:        cfg.MCP.Enabled,
				AllowRuns:      cfg.MCP.AllowRuns,
				AllowMutations: cfg.MCP.AllowMutations,
			},
		},
		srv,
	)
	workerCnt := worker.New(logger, srv)
	metricCnt := metrics.New(logger, ps)

	schedulerCnt := schedulerhandler.NewScheduler(logger, ps, srv)

	logger.Info("EasyJet starting")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		ps.Close()

		return nil
	})

	g.Go(func() error {
		for event := range ps.SubscribeEvent("main", 10) {
			if event.Type == entity.EventRequireAppRestart {
				cancel()
				return nil
			}
		}

		return nil
	})

	g.Go(func() error {
		workerCnt.Start(ctx)
		return nil
	})

	g.Go(func() error {
		metricCnt.Start(ctx)
		return nil
	})

	g.Go(func() error {
		return schedulerCnt.Serve(ctx)
	})

	g.Go(func() error {
		return apiCnt.Serve(ctx)
	})

	err = g.Wait()
	if err != nil {
		logger.Error("server error", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("EasyJet stopped")
}
