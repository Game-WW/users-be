package user_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) LinkIdentity(ctx context.Context, userID string, identity user.ExternalIdentity, now time.Time) error {
	var email any
	if identity.Email != "" {
		email = identity.Email
	}

	sql, args, err := r.builder.
		Insert("user_identities").
		Columns("user_id", "provider", "subject", "email", "created_at").
		Values(userID, identity.Provider, identity.Subject, email, now).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert identity: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return fmt.Errorf("link identity: %w", user.ErrIdentityTaken)
		}
		return fmt.Errorf("insert identity: %w", err)
	}

	return nil
}
