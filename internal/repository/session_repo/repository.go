package session_repo

import (
	"github.com/Game-WW/users-be/internal/core/application/commands"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

var (
	_ commands.SessionWriter  = (*Repository)(nil)
	_ commands.SessionRotator = (*Repository)(nil)
	_ commands.SessionRevoker = (*Repository)(nil)
)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
