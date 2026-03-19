package gorm

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/glebarez/sqlite"
	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRepo(logger *slog.Logger, tp, dns string) (Repo, error) {
	var dialector gorm.Dialector

	switch tp {
	case "sqlite":
		dialector = sqlite.Open(dns)
	case "postgres":
		dialector = postgres.Open(dns)
	default:
		return Repo{}, errors.New("unsupported connection type")
	}

	sloggorm.New()

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: sloggorm.New(
			sloggorm.WithHandler(logger.Handler()),
			sloggorm.WithTraceAll(),
			sloggorm.SetLogLevel(sloggorm.DefaultLogType, slog.LevelDebug),
		),
	})
	if err != nil {
		return Repo{}, fmt.Errorf("gorm open: %w", err)
	}

	err = db.AutoMigrate(
		new(entity.Project),
		new(entity.ProjectRun),
		new(entity.ProjectStage),
		new(entity.ProjectRunStage),
		new(entity.ProjectRunGitCommits),
	)
	if err != nil {
		return Repo{}, fmt.Errorf("gorm automigrate: %w", err)
	}

	return Repo{
		db:     db,
		logger: logger,
	}, nil
}
