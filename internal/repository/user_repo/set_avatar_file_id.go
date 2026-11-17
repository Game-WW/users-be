package user_repo

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
)

var _ commands.AvatarFileUpdater = (*Repository)(nil)

func (r *Repository) SetAvatarFileID(
	ctx context.Context,
	userID string,
	fileID string,
) error {
	sql, args, err := r.builder.
		Update("users").
		Set("avatar_file_id", sq.Expr("NULLIF(?, '')", fileID)).
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build set avatar file id: %w", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute set avatar file id: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}
