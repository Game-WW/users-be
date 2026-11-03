package session_repo

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteByTokenHash(ctx context.Context, hash string) error {
	query, args, err := r.builder.
		Delete("refresh_tokens").
		Where(sq.Eq{"token_hash": hash}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete session query: %w", err)
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session: %w", errs.ErrNotFound)
	}

	return nil
}
