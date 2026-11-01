package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

var ErrIdentityTaken = errors.New("identity already linked to a user")

type UserRepository interface {
	FindByIdentity(ctx context.Context, provider user.Provider, subject string) (*user.User, error)
	Create(ctx context.Context, u *user.User, identity user.ExternalIdentity) error
}

type SessionWriter interface {
	Create(ctx context.Context, s *session.Session) error
}

type LoginHandler struct {
	identities map[user.Provider]IdentityVerifier
	users      UserRepository
	sessions   SessionWriter
	access     AccessTokenIssuer
	clock      Clock
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewLoginHandler(
	identities map[user.Provider]IdentityVerifier,
	users UserRepository,
	sessions SessionWriter,
	access AccessTokenIssuer,
	clock Clock,
	accessTTL, refreshTTL time.Duration,
) *LoginHandler {
	return &LoginHandler{
		identities: identities,
		users:      users,
		sessions:   sessions,
		access:     access,
		clock:      clock,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (LoginResult, error) {
	verifier, ok := h.identities[cmd.Provider]
	if !ok {
		return LoginResult{}, fmt.Errorf("provider %s is not configured: %w", cmd.Provider, errs.ErrInvalidArgument)
	}

	identity, err := verifier.Verify(ctx, cmd.Creds)
	if err != nil {
		return LoginResult{}, fmt.Errorf("verify identity: %w", err)
	}

	now := h.clock.Now()

	u, err := h.findOrCreate(ctx, identity, now)
	if err != nil {
		return LoginResult{}, err
	}

	sess, refreshToken := session.New(u.ID(), cmd.UserAgent, cmd.IP, now, h.refreshTTL)
	if err := h.sessions.Create(ctx, sess); err != nil {
		return LoginResult{}, fmt.Errorf("create session: %w", err)
	}

	accessToken, err := h.access.Issue(u.ID(), now.Add(h.accessTTL))
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue access token: %w", err)
	}

	return LoginResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (h *LoginHandler) findOrCreate(ctx context.Context, identity user.ExternalIdentity, now time.Time) (*user.User, error) {
	u, err := h.users.FindByIdentity(ctx, identity.Provider, identity.Subject)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("find user: %w", err)
	}

	u, err = user.New(identity, now)
	if err != nil {
		return nil, err
	}
	err = h.users.Create(ctx, u, identity)
	if errors.Is(err, ErrIdentityTaken) {
		u, err = h.users.FindByIdentity(ctx, identity.Provider, identity.Subject)
		if err != nil {
			return nil, fmt.Errorf("find user after conflict: %w", err)
		}
		return u, nil
	}
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}
