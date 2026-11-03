package session_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindByTokenHash(ctx context.Context, hash string) (*session.Session, error) {
	query, args, err := r.builder.
		Select(
			"user_id::text",
			"token_hash",
			"COALESCE(user_agent, '')",
			"COALESCE(ip, '')",
			"created_at",
			"expires_at",
		).
		From("refresh_tokens").
		Where(sq.Eq{"token_hash": hash}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select session query: %w", err)
	}

	var st session.State
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&st.UserID, &st.TokenHash, &st.UserAgent, &st.IP, &st.CreatedAt, &st.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("session: %w", errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("select session: %w", err)
	}

	return session.Restore(st), nil
}
