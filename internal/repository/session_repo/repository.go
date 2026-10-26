package session_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/pkg/errs"
)

type Repository struct {
	db *pgxpool.Pool
}

var (
	_ commands.SessionWriter  = (*Repository)(nil)
	_ commands.SessionReader  = (*Repository)(nil)
	_ commands.SessionRevoker = (*Repository)(nil)
)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create — id строки ставит БД (DEFAULT gen_random_uuid()), пустые user_agent и ip пишутся как NULL.
func (r *Repository) Create(ctx context.Context, s *session.Session) error {
	st := s.Snapshot()
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip, created_at, expires_at)
		 VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6)`,
		st.UserID, st.TokenHash, st.UserAgent, st.IP, st.CreatedAt, st.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	return nil
}

func (r *Repository) FindByTokenHash(ctx context.Context, hash string) (*session.Session, error) {
	var st session.State
	err := r.db.QueryRow(ctx,
		`SELECT user_id::text, token_hash, COALESCE(user_agent, ''), COALESCE(ip, ''), created_at, expires_at
		 FROM refresh_tokens WHERE token_hash = $1`, hash,
	).Scan(&st.UserID, &st.TokenHash, &st.UserAgent, &st.IP, &st.CreatedAt, &st.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("session: %w", errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("select session: %w", err)
	}

	return session.Restore(st), nil
}

func (r *Repository) DeleteByTokenHash(ctx context.Context, hash string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, hash)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session: %w", errs.ErrNotFound)
	}

	return nil
}
