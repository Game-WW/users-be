package user_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Create(
	ctx context.Context,
	u *user.User,
	identity user.ExternalIdentity,
	location string,
) error {
	state := u.Snapshot()

	err := pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		sql, args, err := r.builder.
			Insert("users").
			Columns("id", "created_at", "location").
			Values(state.ID, state.CreatedAt, location).
			ToSql()
		if err != nil {
			return fmt.Errorf("build insert user: %w", err)
		}

		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("insert user: %w", err)
		}

		var emailArg any
		if identity.Email != "" {
			emailArg = identity.Email
		}

		sql, args, err = r.builder.
			Insert("user_identities").
			Columns("user_id", "provider", "subject", "email", "created_at").
			Values(state.ID, identity.Provider, identity.Subject, emailArg, state.CreatedAt).
			ToSql()
		if err != nil {
			return fmt.Errorf("build insert identity: %w", err)
		}

		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("insert identity: %w", err)
		}

		return nil
	})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return fmt.Errorf("create user: %w", user.ErrIdentityTaken)
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
