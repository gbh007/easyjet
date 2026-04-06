package service

import (
	"context"
	"fmt"
)

func (srv Service) DeleteProject(ctx context.Context, id int) error {
	err := srv.fs.RemoveProjectDir(ctx, id)
	if err != nil {
		return fmt.Errorf("remove project dir: %w", err)
	}

	err = srv.db.DeleteProject(ctx, id)
	if err != nil {
		return fmt.Errorf("remove project: %w", err)
	}

	return nil
}
