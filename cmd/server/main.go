package main

import (
	"context"
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
	"github.com/gbh007/easyjet/internal/adapter/handler/worker"
	"github.com/gbh007/easyjet/internal/adapter/repository/gorm"
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

	cfg, err := config.Read("config.toml")
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

	db, err := gorm.NewRepo(logger, cfg.Database.Type, cfg.Database.DNS)
	if err != nil {
		logger.Error("create database adapter", "error", err.Error())
		os.Exit(1)
	}

	ex := shellexec.New(logger)
	fs := filesystem.New(logger, cfg.App.ProjectDir)
	git := shellgit.New(logger)

	srv := service.New(logger, ex, fs, git, db)

	apiCnt := httpapi.New(
		logger,
		httpapi.Config{
			Addr: cfg.Server.Addr,
			User: cfg.Server.User,
			Pass: cfg.Server.Pass,
		},
		srv,
	)
	workerCnt := worker.New(logger, srv)

	logger.Info("EasyJet starting")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		workerCnt.Start(ctx)
		return nil
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
