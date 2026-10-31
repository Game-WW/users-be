package user_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolation = "23505"

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

var (
	_ commands.UserRepository        = (*Repository)(nil)
	_ queries.UserProfileReader      = (*Repository)(nil)
	_ queries.LinkedIdentitiesReader = (*Repository)(nil)
)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) GetProfile(ctx context.Context, id string) (queries.GetMeResult, error) {
	sql, args, err := r.builder.
		Select("id::text", "username", "email", "phone_number", "created_at").
		From("users").Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return queries.GetMeResult{}, fmt.Errorf("build get profile: %w", err)
	}

	var res queries.GetMeResult
	err = r.db.QueryRow(ctx, sql, args...).
		Scan(&res.ID, &res.UserName, &res.Email, &res.PhoneNumber, &res.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return queries.GetMeResult{}, fmt.Errorf("user %s: %w", id, errs.ErrNotFound)
	}
	if err != nil {
		return queries.GetMeResult{}, fmt.Errorf("select profile: %w", err)
	}

	return res, nil
}

func (r *Repository) FindByIdentity(ctx context.Context, provider user.Provider, subject string) (*user.User, error) {
	sql, args, err := r.builder.
		Select("u.id::text", "u.username", "u.email", "u.phone_number", "u.created_at").
		From("users u").
		Join("user_identities i ON i.user_id = u.id").
		Where(sq.Eq{"i.provider": provider, "i.subject": subject}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build find by identity: %w", err)
	}

	var state user.State
	err = r.db.QueryRow(ctx, sql, args...).
		Scan(&state.ID, &state.UserName, &state.Email, &state.PhoneNumber, &state.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user by %s: %w", provider, errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("select user by identity: %w", err)
	}

	return user.Restore(state), nil
}

func (r *Repository) Create(ctx context.Context, u *user.User, identity user.ExternalIdentity) error {
	state := u.Snapshot()

	err := pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		sql, args, err := r.builder.
			Insert("users").
			Columns("id", "username", "email", "phone_number", "created_at").
			Values(state.ID, state.UserName, state.Email, state.PhoneNumber, state.CreatedAt).
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
		return fmt.Errorf("create user: %w", commands.ErrIdentityTaken)
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *Repository) GetLinkedIdentities(ctx context.Context, userID string) ([]queries.LinkedIdentity, error) {
	sql, args, err := r.builder.
		Select("provider", "email", "created_at").
		From("user_identities").
		Where(sq.Eq{"user_id": userID}).
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
