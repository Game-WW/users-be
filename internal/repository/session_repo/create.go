package session_repo

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) Create(ctx context.Context, s *session.Session) error {
	st := s.Snapshot()
	query, args, err := r.builder.
		Insert("refresh_tokens").
		Columns("user_id", "token_hash", "user_agent", "ip", "created_at", "expires_at").
		Values(
			st.UserID,
			st.TokenHash,
			sq.Expr("NULLIF(?, '')", st.UserAgent),
			sq.Expr("NULLIF(?, '')", st.IP),
			st.CreatedAt,
			st.ExpiresAt,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert session query: %w", err)
	}
	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	return nil
}
