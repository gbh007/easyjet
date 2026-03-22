package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gbh007/easyjet/migrations"
	"github.com/pressly/goose/v3"
)

type slogGooseAdapter struct {
	logger *slog.Logger
}

func (a slogGooseAdapter) Fatalf(format string, v ...any) {
	a.logger.Error(fmt.Sprintf(format, v...))
}

func (a slogGooseAdapter) Printf(format string, v ...any) {
	a.logger.Info(fmt.Sprintf(format, v...))
}

func migrate(ctx context.Context, logger *slog.Logger, db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)

	err := goose.SetDialect("sqlite3")
	if err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	goose.SetLogger(slogGooseAdapter{
		logger: logger,
	})

	err = goose.UpContext(
		ctx, db, "sqlite",
		goose.WithNoColor(true),
	)
	if err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}

	return nil
}
