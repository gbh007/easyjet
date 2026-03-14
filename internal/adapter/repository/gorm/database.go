package gorm

import (
	"errors"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(tp, dns string) (Repo, error) {
	var dialector gorm.Dialector

	switch tp {
	case "sqlite":
		dialector = sqlite.Open(dns)
	case "postgres":
		dialector = postgres.Open(dns)
	default:
		return Repo{}, errors.New("unsupported connection type")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return Repo{}, fmt.Errorf("gorm open: %w", err)
	}

	err = db.AutoMigrate(
		new(modelProject),
		new(modelProjectRun),
		new(modelProjectStage),
		new(modelProjectStageRun),
	)
	if err != nil {
		return Repo{}, fmt.Errorf("gorm automigrate: %w", err)
	}

	return Repo{
		db: db,
	}, nil
}
