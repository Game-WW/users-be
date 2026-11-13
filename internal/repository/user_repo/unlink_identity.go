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

func (r *Repository) UnlinkIdentity(ctx context.Context, userID string, provider user.Provider) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin unlink identity transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	lockQuery, lockArgs, err := r.builder.
		Select("id::text").
		From("users").
		Where(sq.Eq{"id": userID}).
		Suffix("FOR UPDATE").
		ToSql()
	if err != nil {
		return fmt.Errorf("build lock user query: %w", err)
	}

	var lockedUserID string
	err = tx.QueryRow(ctx, lockQuery, lockArgs...).Scan(&lockedUserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("user: %w", errs.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("lock user: %w", err)
	}

	deleteQuery, deleteArgs, err := r.builder.
		Delete("user_identities").
		Where(sq.Eq{"user_id": userID, "provider": provider}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete identity query: %w", err)
	}

	tag, err := tx.Exec(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return fmt.Errorf("delete identity: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil
	}

	countQuery, countArgs, err := r.builder.
		Select("COUNT(*)").
		From("user_identities").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build count identities query: %w", err)
	}

	var remaining int
	if err := tx.QueryRow(ctx, countQuery, countArgs...).Scan(&remaining); err != nil {
		return fmt.Errorf("count remaining identities: %w", err)
	}
	if remaining == 0 {
		return fmt.Errorf("unlink identity: %w", user.ErrLastIdentity)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit unlink identity transaction: %w", err)
	}

	return nil
}
