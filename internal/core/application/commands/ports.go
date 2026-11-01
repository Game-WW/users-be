package commands

import (
	"context"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/user"
)

type Clock interface {
	Now() time.Time
}

type AccessTokenIssuer interface {
	Issue(userID string, expiresAt time.Time) (string, error)
}

type AuthCredentials struct {
	Token        string // Google
	Code         string // GitHub
	CodeVerifier string // GitHub PKCE
}

type IdentityVerifier interface {
	Verify(ctx context.Context, creds AuthCredentials) (user.ExternalIdentity, error)
}
