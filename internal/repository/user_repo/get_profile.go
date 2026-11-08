package user_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Game-WW/users-be/pkg/errs"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetProfile(ctx context.Context, id string) (queries.GetMeResult, error) {
	sql, args, err := r.builder.
		Select("id::text", "nickname", "created_at", "avatar_url", "location", "show_online_status").
		From("users").Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return queries.GetMeResult{}, fmt.Errorf("build get profile: %w", err)
	}

	var res queries.GetMeResult
	err = r.db.QueryRow(ctx, sql, args...).
		Scan(&res.ID, &res.Nickname, &res.CreatedAt, &res.AvatarURL, &res.Location, &res.ShowOnlineStatus)

	if errors.Is(err, pgx.ErrNoRows) {
		return queries.GetMeResult{}, fmt.Errorf("user %s: %w", id, errs.ErrNotFound)
	}
	if err != nil {
		return queries.GetMeResult{}, fmt.Errorf("select profile: %w", err)
	}

	return res, nil
}
