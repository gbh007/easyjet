package service

import (
	"log/slog"

	"github.com/gbh007/easyjet/internal/core/port"
)

type Service struct {
	ex              port.Exec
	fs              port.FileSystem
	git             port.Git
	db              port.Database
	pubsub          port.PubSub
	logger          *slog.Logger
	externalWebAddr string
}

func New(
	logger *slog.Logger,
	ex port.Exec,
	fs port.FileSystem,
	git port.Git,
	db port.Database,
	publisher port.PubSub,
	externalWebAddr string,
) Service {
	return Service{
		ex:              ex,
		fs:              fs,
		git:             git,
		db:              db,
		pubsub:          publisher,
		logger:          logger,
		externalWebAddr: externalWebAddr,
	}
}
