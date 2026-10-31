package main

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/internal/infrastructure/clock"
	"github.com/Game-WW/users-be/internal/infrastructure/github"
	"github.com/Game-WW/users-be/internal/infrastructure/google"
	"github.com/Game-WW/users-be/internal/infrastructure/jwt"
	"github.com/Game-WW/users-be/internal/infrastructure/yandex"
	"github.com/Game-WW/users-be/internal/repository/session_repo"
	"github.com/Game-WW/users-be/internal/repository/user_repo"
	"github.com/Game-WW/users-be/internal/rpc"
)

func buildServer(cfg config, pool *pgxpool.Pool) (*rpc.Server, error) {
	var (
		users    = user_repo.New(pool)
		sessions = session_repo.New(pool)
		access   = jwt.NewIssuer(cfg.JWTSecret)
		clk      = clock.Clock{}
	)

	verifiers := make(map[user.Provider]commands.IdentityVerifier)
	if cfg.GoogleClientID != "" {
		verifiers[user.ProviderGoogle] = google.NewVerifier(cfg.GoogleClientID)
	}
	if cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "" {
		verifiers[user.ProviderGitHub] = github.NewVerifier(cfg.GitHubClientID, cfg.GitHubClientSecret)
	}

	if cfg.YandexClientID != "" && cfg.YandexClientSecret != "" {
		verifiers[user.ProviderYandex] = yandex.NewVerifier(cfg.YandexClientID, cfg.YandexClientSecret)
	}

	if len(verifiers) == 0 {
		return nil, errors.New("no login provider configured: set GOOGLE_CLIENT_ID or GITHUB_CLIENT_ID")
	}

	return rpc.NewServer(
		commands.NewLoginHandler(verifiers, users, sessions, access, clk, cfg.AccessTokenTTL, cfg.RefreshTokenTTL),
		commands.NewRefreshHandler(sessions, access, clk, cfg.AccessTokenTTL, cfg.RefreshTokenTTL),
		commands.NewLogoutHandler(sessions),
		queries.NewGetMeHandler(users),
		queries.NewListIdentitiesHandler(users),
	), nil
}
