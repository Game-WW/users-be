package session_repo

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) Rotate(ctx context.Context, oldHash string, sess *session.Session) error {
	state := sess.Snapshot()
	query, args, err := r.builder.
		Update("refresh_tokens").
		Set("token_hash", state.TokenHash).
		Set("expires_at", state.ExpiresAt).
		Where(sq.Eq{"token_hash": oldHash}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build rotate session query: %w", err)
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("rotate session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session: %w", errs.ErrNotFound)
	}

	return nil
}
