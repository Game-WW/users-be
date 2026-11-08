package user_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) FindByIdentity(ctx context.Context, provider user.Provider, subject string) (*user.User, error) {
	sql, args, err := r.builder.
		Select("u.id::text", "u.nickname", "u.created_at").
		From("users u").
		Join("user_identities i ON i.user_id = u.id").
		Where(sq.Eq{"i.provider": provider, "i.subject": subject}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find by identity: %w", err)
	}

	var state user.State
	err = r.db.QueryRow(ctx, sql, args...).
		Scan(&state.ID, &state.Nickname, &state.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user by %s: %w", provider, errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("select user by identity: %w", err)
	}

	return user.Restore(state), nil
}
