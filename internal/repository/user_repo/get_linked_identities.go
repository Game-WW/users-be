package user_repo

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Masterminds/squirrel"
)

func (r *Repository) GetLinkedIdentities(ctx context.Context, userID string) ([]queries.LinkedIdentity, error) {
	sql, args, err := r.builder.
		Select("provider", "email", "created_at").
		From("user_identities").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get linked identities: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("select linked identities: %w", err)
	}
	defer rows.Close()

	var res []queries.LinkedIdentity
	for rows.Next() {
		var li queries.LinkedIdentity
		if err := rows.Scan(&li.Provider, &li.Email, &li.LinkedAt); err != nil {
			return nil, fmt.Errorf("scan linked identity: %w", err)
		}

		res = append(res, li)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate linked identities: %w", err)
	}

	return res, nil
}
