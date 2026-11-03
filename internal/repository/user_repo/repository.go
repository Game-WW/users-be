package user_repo

import (
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

var (
	_ commands.UserRepository        = (*Repository)(nil)
	_ queries.UserProfileReader      = (*Repository)(nil)
	_ queries.LinkedIdentitiesReader = (*Repository)(nil)
	_ commands.IdentityLinker        = (*Repository)(nil)
	_ commands.IdentityUnlicker      = (*Repository)(nil)
)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
