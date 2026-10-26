package user_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

const uniqueViolation = "23505"

type Repository struct {
	db *pgxpool.Pool
}

var (
	_ commands.UserRepository   = (*Repository)(nil)
	_ queries.UserProfileReader = (*Repository)(nil)
)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func identityColumn(p user.Provider) (string, error) {
	switch p {
	case user.ProviderGoogle:
		return "google_id", nil
	default:
		return "", fmt.Errorf("no identity column for provider %q", p)
	}
}

func (r *Repository) GetProfile(ctx context.Context, id string) (queries.GetMeResult, error) {
	var res queries.GetMeResult
	err := r.db.QueryRow(ctx,
		`SELECT id::text, username, email, phone_number, created_at FROM users WHERE id = $1`, id,
	).Scan(&res.ID, &res.UserName, &res.Email, &res.PhoneNumber, &res.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return queries.GetMeResult{}, fmt.Errorf("user %s: %w", id, errs.ErrNotFound)
	}
	if err != nil {
		return queries.GetMeResult{}, fmt.Errorf("select profile: %w", err)
	}

	return res, nil
}

func (r *Repository) FindByIdentity(ctx context.Context, provider user.Provider, subject string) (*user.User, error) {
	column, err := identityColumn(provider)
	if err != nil {
		return nil, err
	}

	var s user.State
	err = r.db.QueryRow(ctx,
		`SELECT id::text, username, email, phone_number, created_at FROM users WHERE `+column+` = $1`, subject,
	).Scan(&s.ID, &s.UserName, &s.Email, &s.PhoneNumber, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user by %s: %w", provider, errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("select user by identity: %w", err)
	}

	return user.Restore(s), nil
}

func (r *Repository) Create(ctx context.Context, u *user.User, identity user.ExternalIdentity) error {
	column, err := identityColumn(identity.Provider)
	if err != nil {
		return err
	}

	s := u.Snapshot()
	_, err = r.db.Exec(ctx,
		`INSERT INTO users (id, username, email, phone_number, created_at, `+column+`)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		s.ID, s.UserName, s.Email, s.PhoneNumber, s.CreatedAt, identity.Subject,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return fmt.Errorf("insert user: %w", commands.ErrIdentityTaken)
	}
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}
