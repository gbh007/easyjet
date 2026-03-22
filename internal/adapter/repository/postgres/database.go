package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Repo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	psql   squirrel.StatementBuilderType
}

func NewRepo(ctx context.Context, logger *slog.Logger, dns string) (Repo, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return Repo{}, fmt.Errorf("parse dns: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return Repo{}, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return Repo{}, fmt.Errorf("ping database: %w", err)
	}

	err = migrate(ctx, logger, stdlib.OpenDBFromPool(pool))
	if err != nil {
		return Repo{}, fmt.Errorf("migrate repo: %w", err)
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return Repo{
		pool:   pool,
		logger: logger,
		psql:   psql,
	}, nil
}

func (repo Repo) Close() {
	repo.pool.Close()
}
