package user_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) UpdateProfile(ctx context.Context, cmd commands.UpdateProfileCommand) error {
	update := r.builder.
		Update("users").
		Where(sq.Eq{"id": cmd.UserID})

	if cmd.Nickname != nil {
		update = update.Set("nickname", *cmd.Nickname)
	}
	if cmd.AvatarURL != nil {
		if *cmd.AvatarURL == "" {
			update = update.Set("avatar_url", nil)
		} else {
			update = update.Set("avatar_url", *cmd.AvatarURL)
		}
	}
	if cmd.Location != nil {
		update = update.Set("location", *cmd.Location)
	}
	if cmd.ShowOnlineStatus != nil {
		update = update.Set("show_online_status", *cmd.ShowOnlineStatus)
	}

	sql, args, err := update.ToSql()
	if err != nil {
		return fmt.Errorf("build update profile: %w", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		isNicknameUniqueViolation := errors.As(err, &pgErr) &&
			pgErr.Code == uniqueViolation &&
			pgErr.ConstraintName == "users_nickname_lower_unique"

		if isNicknameUniqueViolation {
			return errs.ErrNicknameTaken
		}
		return fmt.Errorf("execute update profile: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}
