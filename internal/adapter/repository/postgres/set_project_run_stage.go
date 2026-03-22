package postgres

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) SetProjectRunStage(ctx context.Context, rs entity.ProjectRunStage) error {
	query, args, err := repo.psql.
		Insert("run_stages").
		SetMap(map[string]any{
			"run_id":    rs.RunID,
			"stage_num": rs.StageNumber,
			"success":   rs.Success,
			"log":       rs.Log,
		}).
		Suffix("ON CONFLICT (run_id, stage_num) DO UPDATE SET success = EXCLUDED.success, log = EXCLUDED.log").
		ToSql()
	if err != nil {
		return fmt.Errorf("build upsert query: %w", err)
	}

	_, err = repo.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("upsert run stage: %w", err)
	}

	return nil
}
