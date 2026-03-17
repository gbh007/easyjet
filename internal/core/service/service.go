package service

import (
	"log/slog"

	"github.com/gbh007/easyjet/internal/adapter/scheduler"
	"github.com/gbh007/easyjet/internal/core/port"
)

type Service struct {
	ex        port.Exec
	fs        port.FileSystem
	git       port.Git
	db        port.Database
	publisher scheduler.EventPublisher
	logger    *slog.Logger
}

func New(
	logger *slog.Logger,
	ex port.Exec,
	fs port.FileSystem,
	git port.Git,
	db port.Database,
	publisher scheduler.EventPublisher,
) Service {
	return Service{
		ex:        ex,
		fs:        fs,
		git:       git,
		db:        db,
		publisher: publisher,
		logger:    logger,
	}
}
