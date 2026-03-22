package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Masterminds/squirrel"
	_ "github.com/glebarez/go-sqlite" // SQLite driver registration
)

type Repo struct {
	db     *sql.DB
	logger *slog.Logger
	psql   squirrel.StatementBuilderType
}

func NewRepo(ctx context.Context, logger *slog.Logger, dns string) (Repo, error) {
	if !strings.Contains(dns, "?") {
		dns += "?_pragma=foreign_keys(1)"
	}

	db, err := sql.Open("sqlite", dns)
	if err != nil {
		return Repo{}, fmt.Errorf("open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return Repo{}, fmt.Errorf("ping database: %w", err)
	}

	err = migrate(ctx, logger, db)
	if err != nil {
		return Repo{}, fmt.Errorf("migrate repo: %w", err)
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	return Repo{
		db:     db,
		logger: logger,
		psql:   psql,
	}, nil
}

func (repo Repo) Close() {
	_ = repo.db.Close()
}
